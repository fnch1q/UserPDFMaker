package utils

import (
	"UserPDFMaker/internal/data"
	"fmt"
	_ "image/jpeg"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePDF(user data.User) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("TNR", "", "fonts/times_new_roman.ttf")
	pdf.AddUTF8Font("TNRBold", "", "fonts/times_new_roman_bold.ttf")
	pdf.SetFont("TNRBold", "", 11.5)

	text1 := "Наименование объекта"
	text2 := "Реконструкцивйуцвцйацкуацпкацпкпцукпкцпцкпцкпцкпцкпцпкцпцкпцкпцкпцкпя автомобильной дороги А1"

	defaultX := 22.0
	defaultY := 10.0
	x := 90.0
	y := pdf.GetY()
	fmt.Printf("%f\n", y)
	pdf.SetX(x)
	pdf.SetCellMargin(5)
	pdf.MultiCell(108, 5, text2, "1", "C", false)
	y = pdf.GetY() - y
	fmt.Printf("%f\n", y)

	pdf.SetXY(defaultX, defaultY)
	firstCellX := x - defaultX
	pdf.CellFormat(firstCellX, y, text1, "1", 1, "C", false, 0, "")
	pdf.SetX(defaultX)

	pdf.SetFont("TNR", "", 11.5)
	pdf.CellFormat(firstCellX/2, 5, "Номер", "LTR", 1, "C", false, 0, "")
	pdf.SetX(defaultX)
	pdf.CellFormat(firstCellX/2, 5, "п/п", "LBR", 0, "C", false, 0, "")
	pdf.SetXY(pdf.GetX(), pdf.GetY()-5)
	pdf.CellFormat(firstCellX/2, 5, "Обозначение", "LTR", 1, "C", false, 0, "")
	pdf.SetXY(defaultX+firstCellX/2, pdf.GetY())
	pdf.CellFormat(firstCellX/2, 5, "документа", "LBR", 0, "C", false, 0, "")
	pdf.SetXY(pdf.GetX(), 25)
	pdf.CellFormat(60, 10, "Наименование документа", "1", 0, "CM", false, 0, "")
	pdf.CellFormat(48, 5, "Номер последнего", "LTR", 2, "CM", false, 0, "")
	pdf.CellFormat(48, 5, "изменения (версии)", "LBR", 1, "CM", false, 0, "")

	txt1 := "1dddddddddddddddddddddddddddddddddddddddddddddddddddd"
	txt2 := "Раздел ПД №1_ПЗ"
	txt3 := "Раздел 1. «Пояснительная записка»"
	txt4 := "2"
	pdf.SetXY(defaultX, 35)
	pdf.MultiCell(34, 5, txt1, "LR", "C", false)
	pdf.SetXY(defaultX+34, 35)
	pdf.MultiCell(34, 5, txt2, "LR", "C", false)
	pdf.SetXY(defaultX+68, 35)
	pdf.MultiCell(60, 5, txt3, "LR", "C", false)
	pdf.SetXY(defaultX+128, 35)
	pdf.MultiCell(48, 5, txt4, "LR", "C", false)

	err := pdf.OutputFileAndClose("table_output.pdf")
	if err != nil {
		return err
	}

	return nil
}
