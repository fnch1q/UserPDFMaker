package utils

import (
	"UserPDFMaker/internal/data"
	"embed"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/sqweek/dialog"
)

//go:embed fonts/*
var fontFS embed.FS

const (
	label1            = "Наименование объекта"
	label2            = "Номер"
	label3            = "п/п"
	label4            = "Обозначение"
	label5            = "документа"
	label6            = "Наименование документа"
	label7            = "Номер последнего"
	label8            = "изменения (версии)"
	label9            = "Наименование файла"
	label10           = "Дата и время последнего изменения файла"
	label11           = "Размер файла, байт"
	label12           = "Значение контрольной суммы"
	label13           = "Характер работы"
	label14           = "Фамилия"
	label15           = "Подпись"
	label16           = "Дата подписания"
	defaultX          = 22.0
	defaultY          = 10.0
	bottomMargin      = 20.0
	rightMargin       = 12.0
	pageHeight        = 297.0
	pageWidth         = 210.0
	defaultCellHeight = 5.0
	template1         = "Шаблон 1"
	template2         = "Шаблон 2"
)

type cellSize struct {
	objectNameWidth          float64
	serialNumberWidth        float64
	documentDefinitionWidth  float64
	documentNameWidth        float64
	lastVersionUpdateWidth   float64
	hashTypeWidth            float64
	hashValueWidth           float64
	fileNameWidth            float64
	fileUpdatedDateWidth     float64
	fileSizeWidth            float64
	multiFileHashValueWidth  float64
	workTypeWidth            float64
	fullNameWidth            float64
	signatureWidth           float64
	signDateWidth            float64
	labelInfoCertSheetWidth  float64
	infoCertifyingSheetWidth float64
	pageNumberWidth          float64
	limitWidth               float64
}

func GeneratePDF(input data.Input) error {
	savePath, err := dialog.File().Title("Выберите место для сохранения PDF").Filter("PDF file", "pdf").Save()
	if err != nil {
		return err // Пользователь отменил выбор или произошла ошибка
	}

	if filepath.Ext(savePath) != ".pdf" {
		savePath += ".pdf"
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	// Чтение шрифтов из встраиваемых файлов
	fontBytes, err := fontFS.ReadFile("fonts/times_new_roman.ttf")
	if err != nil {
		return err
	}
	pdf.AddUTF8FontFromBytes("TNR", "", fontBytes)

	fontBoldBytes, err := fontFS.ReadFile("fonts/times_new_roman_bold.ttf")
	if err != nil {
		return err
	}
	pdf.AddUTF8FontFromBytes("TNRBold", "", fontBoldBytes)

	pdf.SetFont("TNRBold", "", 11.5)

	size, err := getCellSize(input.Template)
	if err != nil {
		panic("wrong_template")
	}

	x := pageWidth - size.objectNameWidth - rightMargin
	//1-я строка
	pdf.SetX(x)
	pdf.MultiCell(size.objectNameWidth, defaultCellHeight, input.ObjectName, "1", "C", false)
	y := pdf.GetY() - defaultY
	pdf.SetXY(defaultX, defaultY)
	firstCellX := x - defaultX
	pdf.CellFormat(firstCellX, y, label1, "1", 1, "C", false, 0, "")

	pdf.SetFont("TNR", "", 11.5)

	printDocumentLabels(input, size, firstCellX, pdf)

	//3-я строка
	var maxY float64
	currentY := pdf.GetY()
	pdf.SetX(defaultX)
	pdf.MultiCell(size.serialNumberWidth, defaultCellHeight, input.SerialNumber, "", "C", false)
	maxY = pdf.GetY()
	x = defaultX + size.serialNumberWidth
	pdf.SetXY(x, currentY)
	pdf.MultiCell(size.documentDefinitionWidth, defaultCellHeight, input.DocumentDefiniton, "", "C", false)
	if maxY < pdf.GetY() {
		maxY = pdf.GetY()
	}
	x += size.documentDefinitionWidth
	pdf.SetXY(x, currentY)
	pdf.MultiCell(size.documentNameWidth, defaultCellHeight, input.DocumentName, "", "C", false)
	if maxY < pdf.GetY() {
		maxY = pdf.GetY()
	}
	x += size.documentNameWidth
	pdf.SetXY(x, currentY)
	pdf.MultiCell(size.lastVersionUpdateWidth, defaultCellHeight, input.LastVersionUpdateNumber, "", "C", false)
	if maxY < pdf.GetY() {
		maxY = pdf.GetY()
	}

	//Обводка 3-й строки
	pdf.SetXY(defaultX, currentY)
	pdf.CellFormat(size.serialNumberWidth, maxY-currentY, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(size.documentDefinitionWidth, maxY-currentY, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(size.documentNameWidth, maxY-currentY, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(size.lastVersionUpdateWidth, maxY-currentY, "", "1", 1, "", false, 0, "")

	printFilesInfo(input, maxY, size, pdf)

	//7-я строка
	pdf.SetX(defaultX)
	if input.Template == template1 {
		pdf.CellFormat(size.workTypeWidth, defaultCellHeight, label13, "1", 0, "CM", false, 0, "")
		pdf.CellFormat(size.fullNameWidth, defaultCellHeight, label14, "1", 0, "CM", false, 0, "")
		pdf.CellFormat(size.signatureWidth, defaultCellHeight, label15, "1", 0, "CM", false, 0, "")
		pdf.CellFormat(size.signDateWidth, defaultCellHeight, label16, "1", 1, "CM", false, 0, "")
	} else {
		pdf.CellFormat(size.workTypeWidth, defaultCellHeight*2, label13, "1", 0, "CM", false, 0, "")
		pdf.CellFormat(size.fullNameWidth, defaultCellHeight*2, label14, "1", 0, "CM", false, 0, "")
		pdf.CellFormat(size.signatureWidth, defaultCellHeight*2, label15, "1", 0, "CM", false, 0, "")
		pdf.CellFormat(size.signDateWidth, defaultCellHeight, "Дата", "LTR", 2, "CM", false, 0, "")
		pdf.CellFormat(size.signDateWidth, defaultCellHeight, "Подписания", "LBR", 1, "CM", false, 0, "")
	}

	for _, user := range input.Users {
		pdf.SetX(defaultX)
		printUserInfo(user, pdf.GetY(), size, pdf)
	}

	printInfoCertifyingSheet(input, pdf.GetY(), size, pdf)

	err = pdf.OutputFileAndClose(savePath)
	if err != nil {
		return err
	}

	return nil
}

func printDocumentLabels(input data.Input, size cellSize, firstCellX float64, pdf *gofpdf.Fpdf) {
	//2-я строка
	if input.Template == template1 {
		pdf.SetX(defaultX)
		pdf.CellFormat(size.serialNumberWidth, defaultCellHeight, label2, "LTR", 1, "C", false, 0, "")
		pdf.SetX(defaultX)
		pdf.CellFormat(size.serialNumberWidth, defaultCellHeight, label3, "LBR", 0, "C", false, 0, "")
		pdf.SetXY(pdf.GetX(), pdf.GetY()-5)
		pdf.CellFormat(size.documentDefinitionWidth, defaultCellHeight, label4, "LTR", 1, "C", false, 0, "")
		pdf.SetXY(defaultX+firstCellX/2, pdf.GetY())
		pdf.CellFormat(size.documentDefinitionWidth, defaultCellHeight, label5, "LBR", 0, "C", false, 0, "")
		pdf.SetXY(pdf.GetX(), pdf.GetY()-5)
		pdf.CellFormat(size.documentNameWidth, defaultCellHeight*2, label6, "1", 0, "CM", false, 0, "")
		pdf.CellFormat(size.lastVersionUpdateWidth, defaultCellHeight, label7, "LTR", 2, "CM", false, 0, "")
		pdf.CellFormat(size.lastVersionUpdateWidth, defaultCellHeight, label8, "LBR", 1, "CM", false, 0, "")
	} else {
		pdf.SetX(defaultX)
		pdf.CellFormat(size.serialNumberWidth, defaultCellHeight, "", "LTR", 2, "C", false, 0, "")
		pdf.CellFormat(size.serialNumberWidth, defaultCellHeight, label2, "LR", 2, "C", false, 0, "")
		pdf.CellFormat(size.serialNumberWidth, defaultCellHeight, label3, "LR", 2, "C", false, 0, "")
		pdf.CellFormat(size.serialNumberWidth, defaultCellHeight, "", "LBR", 0, "C", false, 0, "")
		pdf.SetXY(pdf.GetX(), pdf.GetY()-15)
		pdf.CellFormat(size.documentDefinitionWidth, defaultCellHeight, "", "LTR", 2, "C", false, 0, "")
		pdf.CellFormat(size.documentDefinitionWidth, defaultCellHeight, label4, "LR", 2, "C", false, 0, "")
		pdf.CellFormat(size.documentDefinitionWidth, defaultCellHeight, label5, "LR", 2, "C", false, 0, "")
		pdf.CellFormat(size.documentDefinitionWidth, defaultCellHeight, "", "LBR", 0, "C", false, 0, "")
		pdf.SetXY(pdf.GetX(), pdf.GetY()-15)
		pdf.CellFormat(size.documentNameWidth, defaultCellHeight*4, label6, "1", 0, "CM", false, 0, "")
		pdf.CellFormat(size.lastVersionUpdateWidth, defaultCellHeight, "Номер", "LTR", 2, "CM", false, 0, "")
		pdf.CellFormat(size.lastVersionUpdateWidth, defaultCellHeight, "последнего", "LR", 2, "CM", false, 0, "")
		pdf.CellFormat(size.lastVersionUpdateWidth, defaultCellHeight, "изменения", "LR", 2, "CM", false, 0, "")
		pdf.CellFormat(size.lastVersionUpdateWidth, defaultCellHeight, "версии", "LBR", 1, "CM", false, 0, "")
	}
}

func printFilesInfo(input data.Input, currentY float64, size cellSize, pdf *gofpdf.Fpdf) {
	//4-я строка
	pdf.SetXY(defaultX, currentY)
	pdf.CellFormat(size.hashTypeWidth, defaultCellHeight, "CRC32", "1", 0, "C", false, 0, "")

	var hashValue string
	if input.Template == template1 {
		hashValue = input.Files[0].Hash
	}
	pdf.CellFormat(size.hashValueWidth, defaultCellHeight, hashValue, "1", 1, "C", false, 0, "")

	y := pdf.GetY()
	//5-я строка
	if input.Template == template1 {
		pdf.SetX(defaultX)
		pdf.CellFormat(size.fileNameWidth, defaultCellHeight*2, "Наименование файла", "1", 0, "CM", false, 0, "")
		x := pdf.GetX()
		pdf.MultiCell(size.fileUpdatedDateWidth, defaultCellHeight, "Дата и время последнего изменения файла", "1", "CM", false)
		pdf.SetXY(x+size.fileUpdatedDateWidth, y)
		pdf.CellFormat(size.fileSizeWidth, defaultCellHeight*2, "Размер файла, байт", "1", 1, "CM", false, 0, "")
	} else {
		pdf.SetX(defaultX)
		pdf.CellFormat(size.fileNameWidth, defaultCellHeight*3, "Наименование файла", "1", 0, "CM", false, 0, "")
		x := pdf.GetX()
		pdf.MultiCell(size.fileUpdatedDateWidth, defaultCellHeight, "Дата и время последнего изменения файла", "1", "CM", false)
		pdf.SetXY(x+size.fileUpdatedDateWidth, y)
		pdf.CellFormat(size.fileSizeWidth, defaultCellHeight*3, "Размер файла, байт", "1", 0, "CM", false, 0, "")
		pdf.MultiCell(size.multiFileHashValueWidth, defaultCellHeight, "Значение контрольной суммы", "1", "CM", false)
	}

	if input.Template == template1 {
		//6-я строка
		pdf.SetX(defaultX)
		y = pdf.GetY()
		pdf.MultiCell(size.fileNameWidth, defaultCellHeight, input.Files[0].Name, "1", "CM", false)
		calcY := pdf.GetY() - y
		pdf.SetXY(defaultX+size.fileNameWidth, y)
		updateTime := input.Files[0].UpdateTime
		pdf.CellFormat(size.fileUpdatedDateWidth, calcY, updateTime.Format("02.01.2006 15:04"), "1", 0, "CM", false, 0, "")
		pdf.CellFormat(size.fileSizeWidth, calcY, strconv.Itoa(int(input.Files[0].Size)), "1", 1, "CM", false, 0, "")
	} else {
		for _, file := range input.Files {
			fileNameArr := pdf.SplitText(file.Name, size.fileNameWidth)
			fileNameTopMargin := 0.0
			if len(fileNameArr) <= 1 {
				fileNameTopMargin = 2.5
			}

			calcHeight := 10.0
			if len(fileNameArr) > 2 {
				calcHeight = float64(len(fileNameArr)) * defaultCellHeight
			}

			remainingSpace := pageHeight - (currentY + calcHeight)

			// Если до конца страницы осталось <= 20, создаём новую страницу
			if remainingSpace < bottomMargin {
				pdf.AddPage()
				pdf.SetXY(defaultX, defaultY)
			}

			y = pdf.GetY()
			pdf.SetXY(defaultX, pdf.GetY()+fileNameTopMargin)
			pdf.MultiCell(size.fileNameWidth, defaultCellHeight, file.Name, "", "LM", false)
			pdf.SetXY(defaultX, y)
			pdf.CellFormat(size.fileNameWidth, calcHeight, "", "1", 0, "", false, 0, "")

			pdf.SetXY(defaultX+size.fileNameWidth, y)
			updateTime := file.UpdateTime
			pdf.CellFormat(size.fileUpdatedDateWidth, calcHeight, updateTime.Format("02.01.2006 15:04"), "1", 0, "CM", false, 0, "")
			pdf.CellFormat(size.fileSizeWidth, calcHeight, strconv.Itoa(int(file.Size)), "1", 0, "CM", false, 0, "")
			pdf.CellFormat(size.multiFileHashValueWidth, calcHeight, file.Hash, "1", 1, "CM", false, 0, "")
		}
	}

}

// вывод строки с информацией об ИУЛ и количестве листов (последняя строка)
func printInfoCertifyingSheet(input data.Input, currentY float64, size cellSize, pdf *gofpdf.Fpdf) {
	maxY := defaultCellHeight * 2
	icsArr := pdf.SplitText(input.InfoCertifyingSheet, size.infoCertifyingSheetWidth)
	lenIcs := len(icsArr)

	if lenIcs*defaultCellHeight > int(maxY) {
		maxY = float64(lenIcs) * defaultCellHeight
	}

	remainingSpace := pageHeight - (currentY + maxY)

	// Если до конца страницы осталось <= 20, создаём новую страницу
	if remainingSpace < bottomMargin {
		pdf.AddPage()
		currentY = defaultY
	}

	pdf.SetXY(defaultX+size.labelInfoCertSheetWidth, currentY)
	if lenIcs <= 1 {
		pdf.SetXY(defaultX+size.labelInfoCertSheetWidth, currentY+2.5)
	}
	pdf.MultiCell(size.infoCertifyingSheetWidth, defaultCellHeight, input.InfoCertifyingSheet, "", "CM", false)

	topMargin := 0.0
	y := pdf.GetY() - currentY
	if y > defaultCellHeight*2 {
		topMargin = (y - 10) / 2
	}

	pdf.SetXY(defaultX, currentY+topMargin)

	pdf.CellFormat(size.labelInfoCertSheetWidth, defaultCellHeight, "Информационно-", "", 2, "CM", false, 0, "")
	pdf.CellFormat(size.labelInfoCertSheetWidth, defaultCellHeight, "удостоверяющий лист", "", 0, "CM", false, 0, "")

	x := defaultX + size.labelInfoCertSheetWidth + size.infoCertifyingSheetWidth
	pdf.SetXY(x, currentY)
	pdf.CellFormat(size.pageNumberWidth, defaultCellHeight, "Лист", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(size.limitWidth, defaultCellHeight, "Листов", "1", 1, "CM", false, 0, "")

	pdf.SetX(x)
	pdf.CellFormat(size.pageNumberWidth, defaultCellHeight, fmt.Sprintf("%d", input.Page), "", 0, "CM", false, 0, "")
	pdf.CellFormat(size.limitWidth, defaultCellHeight, fmt.Sprintf("%d", input.Limit), "", 0, "CM", false, 0, "")

	pdf.SetXY(defaultX, currentY)
	pdf.CellFormat(size.labelInfoCertSheetWidth, maxY, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(size.infoCertifyingSheetWidth, maxY, "", "1", 0, "", false, 0, "")

	pdf.SetXY(x, currentY+5)
	pdf.CellFormat(size.pageNumberWidth, maxY-5, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(size.limitWidth, maxY-5, "", "1", 0, "", false, 0, "")
}

// вывод пользователей
func printUserInfo(user data.User, currentY float64, size cellSize, pdf *gofpdf.Fpdf) {
	workTypeArr := pdf.SplitText(user.WorkType, size.workTypeWidth)
	fullNameArr := pdf.SplitText(user.FullName, size.fullNameWidth)
	maxKoef := math.Max(float64(len(workTypeArr)), float64(len(fullNameArr)))

	maxY := defaultCellHeight * 2
	cellHeight := defaultCellHeight
	if maxKoef != 0 {
		cellHeight *= maxKoef
	}
	if cellHeight > maxY {
		maxY = cellHeight
	}

	remainingSpace := pageHeight - (currentY + maxY)

	// Если до конца страницы осталось <= 20, создаём новую страницу
	if remainingSpace < bottomMargin {
		pdf.AddPage()
		currentY = defaultY
	}

	var topMargin float64
	isImage, _ := checkImageFormat(user.Signature)
	lenFullName, lenWorkType := len(fullNameArr), len(workTypeArr)

	if lenFullName == 0 {
		lenFullName = 1
	}
	if lenWorkType == 0 {
		lenWorkType = 1
	}

	if lenFullName == 1 && lenWorkType == 1 {
		topMargin = 2.5
	}

	if lenFullName > lenWorkType {
		pdf.SetXY(defaultX+size.workTypeWidth, currentY)
		pdf.MultiCell(size.fullNameWidth, defaultCellHeight, user.FullName, "", "L", false)
		if topMargin != 2.5 {
			topMargin = ((pdf.GetY() - currentY) - (float64(lenWorkType) * 5)) / 2
		}
		pdf.SetXY(defaultX, currentY+topMargin)
		pdf.MultiCell(size.workTypeWidth, defaultCellHeight, user.WorkType, "", "CM", false)
	} else {
		pdf.SetXY(defaultX, currentY+topMargin)
		pdf.MultiCell(size.workTypeWidth, defaultCellHeight, user.WorkType, "", "CM", false)
		if topMargin != 2.5 {
			topMargin = ((pdf.GetY() - currentY) - (float64(lenFullName) * 5)) / 2
		}
		pdf.SetXY(defaultX+size.workTypeWidth, currentY+topMargin)
		pdf.MultiCell(size.fullNameWidth, defaultCellHeight, user.FullName, "", "L", false)
	}

	pdf.SetXY(defaultX+size.workTypeWidth+size.fullNameWidth, currentY)

	if isImage {
		// Вставляем изображение
		imageHeight := defaultCellHeight * 2
		if maxY == defaultCellHeight {
			maxY = defaultCellHeight * 2
		}
		imageY := currentY
		if maxY > defaultCellHeight*2 {
			imageY += (maxY - imageHeight) / 2
		}
		x := defaultX + size.fullNameWidth + size.workTypeWidth
		pdf.ImageOptions(user.Signature, x, imageY, size.signatureWidth, imageHeight, false, gofpdf.ImageOptions{}, 0, "")
	} else {
		// Если изображение не найдено
		pdf.CellFormat(size.signatureWidth, maxY, "", "1", 0, "CM", false, 0, "")
	}

	x := defaultX + size.fullNameWidth + size.workTypeWidth + size.signatureWidth
	pdf.SetXY(x, currentY)
	t := time.Now()
	pdf.CellFormat(size.signDateWidth, maxY, t.Format("02.01.2006"), "1", 1, "CM", false, 0, "")

	pdf.SetXY(defaultX, currentY)
	pdf.CellFormat(size.workTypeWidth, maxY, "", "1", 0, "", false, 0, "")
	x = defaultX + size.workTypeWidth
	pdf.SetXY(x, currentY)
	pdf.CellFormat(size.fullNameWidth, maxY, "", "1", 0, "", false, 0, "")
	x += size.fullNameWidth
	pdf.SetXY(x, currentY)
	pdf.CellFormat(size.signatureWidth, maxY, "", "1", 0, "", false, 0, "")
	x += size.signatureWidth
	pdf.SetXY(x, currentY)
	pdf.CellFormat(size.signDateWidth, maxY, "", "1", 1, "", false, 0, "")
}

func checkImageFormat(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Определение формата изображения
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return false, err
	}

	if format == "jpeg" || format == "png" || format == "jpg" {
		return true, nil
	} else {
		return false, nil
	}
}

func getCellSize(template string) (cellSize, error) {
	var size cellSize
	if template == template1 {
		size.objectNameWidth = 108.0
		size.serialNumberWidth = 34.0
		size.documentDefinitionWidth = 34.0
		size.documentNameWidth = 60.0
		size.lastVersionUpdateWidth = 48.0
		size.hashTypeWidth = 68.0
		size.hashValueWidth = 108.0
		size.fileNameWidth = 68.0
		size.fileUpdatedDateWidth = 52.0
		size.fileSizeWidth = 56.0
		size.workTypeWidth = 46.0
		size.fullNameWidth = 46.0
		size.signatureWidth = 38.0
		size.signDateWidth = 46.0
		size.labelInfoCertSheetWidth = 55.0
		size.infoCertifyingSheetWidth = 75.0
		size.pageNumberWidth = 21.0
		size.limitWidth = 25.0
	} else if template == template2 {
		size.objectNameWidth = 120.0
		size.serialNumberWidth = 12.0
		size.documentDefinitionWidth = 44.0
		size.documentNameWidth = 90.0
		size.lastVersionUpdateWidth = 30.0
		size.hashTypeWidth = 56.0
		size.hashValueWidth = 120.0
		size.fileNameWidth = 56.0
		size.fileUpdatedDateWidth = 45.0
		size.fileSizeWidth = 45.0
		size.multiFileHashValueWidth = 30.0
		size.workTypeWidth = 56.0
		size.fullNameWidth = 45.0
		size.signatureWidth = 45.0
		size.signDateWidth = 30.0
		size.labelInfoCertSheetWidth = 101.0
		size.infoCertifyingSheetWidth = 45.0
		size.pageNumberWidth = 12.0
		size.limitWidth = 18.0
	} else {
		return cellSize{}, errors.New("wrong_template")
	}

	return size, nil
}
