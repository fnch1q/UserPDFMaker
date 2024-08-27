package utils

import (
	"UserPDFMaker/internal/data"
	_ "image/jpeg"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePDF(user data.User) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Добавляем шрифт для поддержки кириллицы
	pdf.AddUTF8Font("TimesNewRoman", "", "fonts/times_new_roman.ttf")
	pdf.SetFont("TimesNewRoman", "", 12)

	// Устанавливаем ширину для каждой ячейки
	colWidth1 := 95.0
	colWidth2 := 95.0

	// Первый столбец
	pdf.SetXY(10, 10)
	pdf.MultiCell(colWidth1, 10, "Длинный текст для первого столбца, который может занимать несколько строк, и его высота будет автоматически подстроена.", "1", "L", false)
	y1 := pdf.GetY() // Получаем высоту после первой ячейки

	// Второй столбец
	pdf.SetXY(105, 10) // Позиционируем по горизонтали, оставляя место для первого столбца
	pdf.MultiCell(colWidth2, 10, "Короткий текст для второго столбца.", "1", "L", false)
	y2 := pdf.GetY() // Получаем высоту после второй ячейки

	// Определяем максимальную высоту для выравнивания
	maxHeight := y1
	if y2 > maxHeight {
		maxHeight = y2
	}

	// Устанавливаем высоту для следующего контента
	pdf.SetY(maxHeight + 10)

	// Сохраняем PDF файл
	err := pdf.OutputFileAndClose("table_output.pdf")
	if err != nil {
		return err
	}

	return nil
}
