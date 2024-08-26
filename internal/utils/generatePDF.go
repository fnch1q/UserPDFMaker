package utils

import (
	"UserPDFMaker/internal/data"
	_ "image/jpeg"

	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/encoding/charmap"
)

func encodeToWindows1251(input string) string { // Не работает
	encoder := charmap.Windows1251.NewEncoder()
	output, _ := encoder.String(input)
	return output
}

func GeneratePDF(user data.User) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, encodeToWindows1251("ФИО: "+user.FullName))
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, encodeToWindows1251("Вид работы: "+user.WorkType))
	pdf.Ln(12)

	// Вставка подписи
	if user.Signature != "" {
		pdf.Image(user.Signature, 10, 50, 40, 20, false, "", 0, "")
	}

	// Сохранение PDF файла
	err := pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		return err
	}

	return nil
}

// TO DO (Кириллица)
