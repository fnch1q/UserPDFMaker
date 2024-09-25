package utils

import (
	"UserPDFMaker/internal/data"
	"embed"
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
	label9            = "Характер работы"
	label10           = "Фамилия"
	label11           = "Подпись"
	label12           = "Дата подписания"
	defaultX          = 22.0
	defaultY          = 10.0
	bottomMargin      = 20.0
	pageHeight        = 297.0
	defaultCellHeight = 5.0
)

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

	x := 90.0
	pdf.SetX(x)
	pdf.MultiCell(108, 5, input.ObjectName, "1", "C", false)
	y := pdf.GetY() - defaultY
	pdf.SetXY(defaultX, defaultY)
	firstCellX := x - defaultX
	pdf.CellFormat(firstCellX, y, label1, "1", 1, "C", false, 0, "")

	pdf.SetFont("TNR", "", 11.5)

	pdf.SetX(defaultX)
	pdf.CellFormat(firstCellX/2, 5, label2, "LTR", 1, "C", false, 0, "")
	pdf.SetX(defaultX)
	pdf.CellFormat(firstCellX/2, 5, label3, "LBR", 0, "C", false, 0, "")
	pdf.SetXY(pdf.GetX(), pdf.GetY()-5)
	pdf.CellFormat(firstCellX/2, 5, label4, "LTR", 1, "C", false, 0, "")
	pdf.SetXY(defaultX+firstCellX/2, pdf.GetY())
	pdf.CellFormat(firstCellX/2, 5, label5, "LBR", 0, "C", false, 0, "")
	pdf.SetXY(pdf.GetX(), pdf.GetY()-5)
	pdf.CellFormat(60, 10, label6, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(48, 5, label7, "LTR", 2, "CM", false, 0, "")
	pdf.CellFormat(48, 5, label8, "LBR", 1, "CM", false, 0, "")

	var maxY float64
	currentY := pdf.GetY()
	pdf.SetX(defaultX)
	pdf.MultiCell(34, 5, input.SerialNumber, "", "C", false)
	maxY = pdf.GetY()
	pdf.SetXY(defaultX+34, currentY)
	pdf.MultiCell(34, 5, input.DocumentDefiniton, "", "C", false)
	if maxY < pdf.GetY() {
		maxY = pdf.GetY()
	}
	pdf.SetXY(defaultX+68, currentY)
	pdf.MultiCell(60, 5, input.DocumentName, "", "C", false)
	if maxY < pdf.GetY() {
		maxY = pdf.GetY()
	}
	pdf.SetXY(defaultX+128, currentY)
	pdf.MultiCell(48, 5, input.LastVersionUpdateNumber, "", "C", false)
	if maxY < pdf.GetY() {
		maxY = pdf.GetY()
	}

	pdf.SetXY(defaultX, currentY)
	pdf.CellFormat(34, maxY-currentY, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(34, maxY-currentY, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(60, maxY-currentY, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(48, maxY-currentY, "", "1", 1, "", false, 0, "")

	pdf.SetXY(defaultX, maxY)
	pdf.CellFormat(68, 5, "CRC32", "1", 0, "C", false, 0, "")
	pdf.CellFormat(108, 5, input.Files[0].Hash, "1", 1, "C", false, 0, "")

	y = pdf.GetY()
	pdf.SetX(defaultX)
	pdf.CellFormat(68, 10, "Наименование файла", "1", 0, "CM", false, 0, "")
	x = pdf.GetX()
	pdf.MultiCell(52, 5, "Дата и время последнего изменения файла", "1", "CM", false)
	pdf.SetXY(x+52, y)
	pdf.CellFormat(56, 10, "Размер файла, байт", "1", 1, "CM", false, 0, "")

	pdf.SetX(defaultX)
	y = pdf.GetY()
	pdf.MultiCell(68, 5, input.Files[0].Name, "1", "CM", false)
	calcY := pdf.GetY() - y
	pdf.SetXY(defaultX+68, y)
	updateTime := input.Files[0].UpdateTime
	pdf.CellFormat(52, calcY, updateTime.Format("02.01.2006 15:04"), "1", 0, "CM", false, 0, "")
	pdf.CellFormat(56, calcY, strconv.Itoa(int(input.Files[0].Size)), "1", 1, "CM", false, 0, "")

	pdf.SetX(defaultX)
	pdf.CellFormat(46, 5, label9, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(46, 5, label10, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(38, 5, label11, "1", 0, "CM", false, 0, "")
	pdf.CellFormat(46, 5, label12, "1", 1, "CM", false, 0, "")

	for _, user := range input.Users {
		pdf.SetX(defaultX)
		printUserInfo(user, pdf.GetY(), pdf)
	}

	printInfoCertifyingSheet(input, pdf.GetY(), pdf)

	err = pdf.OutputFileAndClose(savePath)
	if err != nil {
		return err
	}

	return nil
}

// вывод строки с информацией об ИУЛ и количестве листов (последняя строка)
func printInfoCertifyingSheet(input data.Input, currentY float64, pdf *gofpdf.Fpdf) {
	maxY := 10.0
	icsArr := pdf.SplitText(input.InfoCertifyingSheet, 75)
	lenIcs := len(icsArr)

	if lenIcs*defaultCellHeight > 10 {
		maxY = float64(lenIcs) * defaultCellHeight
	}

	remainingSpace := pageHeight - (currentY + maxY)

	// Если до конца страницы осталось <= 20, создаём новую страницу
	if remainingSpace < bottomMargin {
		pdf.AddPage()
		currentY = defaultY
	}

	pdf.SetXY(defaultX+55, currentY)
	if lenIcs <= 1 {
		pdf.SetXY(defaultX+55, currentY+2.5)
	}
	pdf.MultiCell(75, 5, input.InfoCertifyingSheet, "", "CM", false)

	topMargin := 0.0
	y := pdf.GetY() - currentY
	if y > 10 {
		topMargin = (y - 10) / 2
	}

	pdf.SetXY(defaultX, currentY+topMargin)

	pdf.CellFormat(55, 5, "Информационно-", "", 2, "CM", false, 0, "")
	pdf.CellFormat(55, 5, "удостоверяющий лист", "", 0, "CM", false, 0, "")

	pdf.SetXY(defaultX+130, currentY)
	pdf.CellFormat(21, 5, "Лист", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 5, "Листов", "1", 1, "CM", false, 0, "")

	pdf.SetX(defaultX + 130)
	pdf.CellFormat(21, 5, fmt.Sprintf("%d", input.Page), "", 0, "CM", false, 0, "")
	pdf.CellFormat(25, 5, fmt.Sprintf("%d", input.Limit), "", 0, "CM", false, 0, "")

	pdf.SetXY(defaultX, currentY)
	pdf.CellFormat(55, maxY, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(75, maxY, "", "1", 0, "", false, 0, "")

	pdf.SetXY(defaultX+130, currentY+5)
	pdf.CellFormat(21, maxY-5, "", "1", 0, "", false, 0, "")
	pdf.CellFormat(25, maxY-5, "", "1", 0, "", false, 0, "")
}

// вывод пользователей
func printUserInfo(user data.User, currentY float64, pdf *gofpdf.Fpdf) {
	workTypeArr := pdf.SplitText(user.WorkType, 46)
	fullNameArr := pdf.SplitText(user.FullName, 46)
	maxKoef := math.Max(float64(len(workTypeArr)), float64(len(fullNameArr)))

	maxY := 10.0
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

	var y float64
	isImage, _ := checkImageFormat(user.Signature)
	lenFullName, lenWorkType := len(fullNameArr), len(workTypeArr)

	if lenFullName == 0 {
		lenFullName = 1
	}
	if lenWorkType == 0 {
		lenWorkType = 1
	}

	if lenFullName == 1 && lenWorkType == 1 {
		y = 2.5
	}

	if lenFullName > lenWorkType {
		pdf.SetXY(defaultX+46, currentY)
		pdf.MultiCell(46, 5, user.FullName, "", "L", false)
		if y != 2.5 {
			y = ((pdf.GetY() - currentY) - (float64(lenWorkType) * 5)) / 2
		}
		pdf.SetXY(defaultX, currentY+y)
		pdf.MultiCell(46, 5, user.WorkType, "", "CM", false)
	} else {
		pdf.SetXY(defaultX, currentY+y)
		pdf.MultiCell(46, 5, user.WorkType, "", "CM", false)
		if y != 2.5 {
			y = ((pdf.GetY() - currentY) - (float64(lenFullName) * 5)) / 2
		}
		pdf.SetXY(defaultX+46, currentY+y)
		pdf.MultiCell(46, 5, user.FullName, "", "L", false)
	}

	pdf.SetXY(defaultX+92, currentY)

	if isImage {
		// Вставляем изображение
		imageHeight := 10.0
		if maxY == 5 {
			maxY = 10
		}
		imageY := currentY
		if maxY > 10 {
			imageY += (maxY - imageHeight) / 2
		}
		pdf.ImageOptions(user.Signature, defaultX+92, imageY, 38, imageHeight, false, gofpdf.ImageOptions{}, 0, "")
	} else {
		// Если изображение не найдено
		pdf.CellFormat(38, maxY, "", "1", 0, "CM", false, 0, "")
	}

	pdf.SetXY(defaultX+130, currentY)
	t := time.Now()
	pdf.CellFormat(46, maxY, t.Format("02.01.2006"), "1", 1, "CM", false, 0, "")

	pdf.SetXY(defaultX, currentY)
	pdf.CellFormat(46, maxY, "", "1", 0, "", false, 0, "")
	pdf.SetXY(defaultX+46, currentY)
	pdf.CellFormat(46, maxY, "", "1", 0, "", false, 0, "")
	pdf.SetXY(defaultX+92, currentY)
	pdf.CellFormat(38, maxY, "", "1", 0, "", false, 0, "")
	pdf.SetXY(defaultX+130, currentY)
	pdf.CellFormat(46, maxY, "", "1", 1, "", false, 0, "")
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
