package utils

import (
	"UserPDFMaker/internal/data"
	"image"
	_ "image/jpeg"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const (
	label1   = "Наименование объекта"
	label2   = "Номер"
	label3   = "п/п"
	label4   = "Обозначение"
	label5   = "документа"
	label6   = "Наименование документа"
	label7   = "Номер последнего"
	label8   = "изменения (версии)"
	label9   = "Характер работы"
	label10  = "Фамилия"
	label11  = "Подпись"
	label12  = "Дата подписания"
	defaultX = 22.0
)

func GeneratePDF(input data.Input) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("TNR", "", "fonts/times_new_roman.ttf")
	pdf.AddUTF8Font("TNRBold", "", "fonts/times_new_roman_bold.ttf")
	pdf.SetFont("TNRBold", "", 11.5)

	defaultY := 10.0
	x := 90.0
	y := pdf.GetY()
	pdf.SetX(x)
	pdf.MultiCell(108, 5, input.ObjectName, "1", "C", false)
	y = pdf.GetY() - y

	pdf.SetXY(defaultX, defaultY)
	firstCellX := x - defaultX
	pdf.CellFormat(firstCellX, y, label1, "1", 1, "C", false, 0, "")
	pdf.SetX(defaultX)

	pdf.SetFont("TNR", "", 11.5)
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
		outputUserInfo(user, pdf.GetY(), pdf)
	}

	err := pdf.OutputFileAndClose("output_" + time.Now().Format("02.01.2006 15:04") + ".pdf")
	if err != nil {
		return err
	}

	return nil
}

func outputUserInfo(user data.User, currentY float64, pdf *gofpdf.Fpdf) {
	workTypeArr := pdf.SplitText(user.WorkType, 46)
	fullNameArr := pdf.SplitText(user.FullName, 46)
	maxKoef := math.Max(float64(len(workTypeArr)), float64(len(fullNameArr)))

	maxY := 5.0
	if maxKoef != 0 {
		maxY *= maxKoef
	}

	var y float64
	isImage, _ := checkImageFormat(user.Signature)
	if len(fullNameArr) == 1 && len(workTypeArr) == 1 && isImage {
		y = 2.5
	}

	if len(fullNameArr) > len(workTypeArr) {
		pdf.SetXY(defaultX+46, currentY)
		pdf.MultiCell(46, 5, user.FullName, "", "L", false)
		if y != 2.5 {
			y = ((pdf.GetY() - currentY) - (float64(len(workTypeArr)) * 5)) / 2
		}
		pdf.SetXY(defaultX, currentY+y)
		pdf.MultiCell(46, 5, user.WorkType, "", "CM", false)
	} else {
		pdf.SetXY(defaultX, currentY+y)
		pdf.MultiCell(46, 5, user.WorkType, "", "CM", false)
		if y != 2.5 {
			y = ((pdf.GetY() - currentY) - (float64(len(fullNameArr)) * 5)) / 2
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

	if format == "jpeg" || format == "png" {
		return true, nil
	} else {
		return false, nil
	}
}
