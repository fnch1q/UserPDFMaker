package ui

import (
	"UserPDFMaker/internal/data"
	"fmt"
	"strconv"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateTemplateGroup(input *data.Input) *fyne.Container {
	// Убедитесь, что TemplateSelect создан
	if input.Widgets.TemplateSelect == nil {
		templateOptions := []string{template1, template2}
		input.Widgets.TemplateSelect = widget.NewSelect(templateOptions, func(value string) {
			input.Template = value
			// Очищаем файлы при смене шаблона
			input.Files = nil
			if input.Widgets.FileList != nil {
				input.Widgets.FileList.RemoveAll() // Очищаем визуальный список файлов
			}
			fmt.Println("Файлы удалены из списка и структуры после смены шаблона.")
		})
	}

	label := widget.NewLabel("Выберите шаблон:")
	label.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		label,
		input.Widgets.TemplateSelect,
	)
}

func CreateDocumentDetailsGroup(input *data.Input) *fyne.Container {
	maxLength := 256
	// Инициализация виджетов, если они не были созданы
	if input.Widgets.ObjectName == nil {
		input.Widgets.ObjectName = widget.NewEntry()
	}
	if input.Widgets.SerialNumber == nil {
		input.Widgets.SerialNumber = widget.NewEntry()
	}
	if input.Widgets.DocumentDefiniton == nil {
		input.Widgets.DocumentDefiniton = widget.NewEntry()
	}
	if input.Widgets.DocumentName == nil {
		input.Widgets.DocumentName = widget.NewEntry()
	}
	if input.Widgets.LastVersionUpdateNumber == nil {
		input.Widgets.LastVersionUpdateNumber = widget.NewEntry()
	}
	if input.Widgets.InfoCertifyingSheet == nil {
		input.Widgets.InfoCertifyingSheet = widget.NewEntry()
	}
	if input.Widgets.Page == nil {
		input.Widgets.Page = widget.NewEntry()
	}
	if input.Widgets.Limit == nil {
		input.Widgets.Limit = widget.NewEntry()
	}

	input.Widgets.ObjectName.SetPlaceHolder("Наименование объекта")
	input.Widgets.ObjectName.OnChanged = func(text string) {
		if len(text) > maxLength {
			text = text[:maxLength]
			input.Widgets.ObjectName.SetText(text) // update the widget's text
		}
		input.ObjectName = text
		fmt.Println(input)
	}
	input.Widgets.ObjectName.Validator = (func(text string) error {
		if len(text) > maxLength {
			input.Widgets.ObjectName.SetText(text[:maxLength])
		}
		return nil
	})

	input.Widgets.SerialNumber.SetPlaceHolder("Порядковый номер документа")
	input.Widgets.SerialNumber.Validator = (func(text string) error {
		if !isNumeric(text) {
			input.Widgets.SerialNumber.SetText(text[:len(text)-2])
		} else {
			input.SerialNumber = text
			fmt.Println(input)
		}
		return nil
	})

	input.Widgets.DocumentDefiniton.SetPlaceHolder("Обозначение ДЭ")
	input.Widgets.DocumentDefiniton.OnChanged = func(text string) {
		if len(text) > maxLength {
			text = text[:maxLength]
			input.Widgets.DocumentDefiniton.SetText(text) // update the widget's text
		}
		input.DocumentDefiniton = text
		fmt.Println(input)
	}
	input.Widgets.DocumentDefiniton.Validator = (func(text string) error {
		if len(text) > maxLength {
			input.Widgets.DocumentDefiniton.SetText(text[:maxLength])
		}
		return nil
	})

	input.Widgets.DocumentName.SetPlaceHolder("Наименование документа")
	input.Widgets.DocumentName.OnChanged = func(text string) {
		if len(text) > maxLength {
			text = text[:maxLength]
			input.Widgets.DocumentName.SetText(text) // update the widget's text
		}
		input.DocumentName = text
		fmt.Println(input)
	}
	input.Widgets.DocumentName.Validator = (func(text string) error {
		if len(text) > maxLength {
			input.Widgets.DocumentName.SetText(text[:maxLength])
		}
		return nil
	})

	input.Widgets.LastVersionUpdateNumber.SetPlaceHolder("Номер последнего изменения")
	input.Widgets.LastVersionUpdateNumber.Validator = (func(text string) error {
		if !isNumeric(text) {
			input.Widgets.LastVersionUpdateNumber.SetText(text[:len(text)-2])
		} else {
			input.LastVersionUpdateNumber = text
			fmt.Println(input)
		}
		return nil
	})

	input.Widgets.InfoCertifyingSheet.SetPlaceHolder("Обозначение ИУЛ")
	input.Widgets.InfoCertifyingSheet.OnChanged = func(text string) {
		if len(text) > maxLength {
			text = text[:maxLength]
			input.Widgets.InfoCertifyingSheet.SetText(text) // update the widget's text
		}
		input.InfoCertifyingSheet = text
		fmt.Println(input)
	}
	input.Widgets.InfoCertifyingSheet.Validator = (func(text string) error {
		if len(text) > maxLength {
			input.Widgets.InfoCertifyingSheet.SetText(text[:maxLength])
		}
		return nil
	})

	input.Widgets.Page.SetPlaceHolder("Номер страницы ИУЛ")
	input.Widgets.Page.Validator = (func(text string) error {
		if !isNumeric(text) {
			input.Widgets.Page.SetText(text[:len(text)-2])
		} else {
			input.Page = textToInt(text)
			fmt.Println(input)
		}
		return nil
	})

	input.Widgets.Limit.SetPlaceHolder("Количество листов ИУЛ")
	input.Widgets.Limit.Validator = (func(text string) error {
		if !isNumeric(text) {
			input.Widgets.Limit.SetText(text[:len(text)-2])
		} else {
			input.Limit = textToInt(text)
			fmt.Println(input)
		}
		return nil
	})

	// Заголовок группы
	label := widget.NewLabel("Детали документа")
	label.TextStyle = fyne.TextStyle{Bold: true}

	// Собираем все в контейнер с использованием Grid
	return container.NewVBox(
		label,
		container.NewGridWithColumns(2,
			widget.NewLabel("Наименование объекта:"), input.Widgets.ObjectName,
			widget.NewLabel("Порядковый номер документа:"), input.Widgets.SerialNumber,
			widget.NewLabel("Обозначение ДЭ:"), input.Widgets.DocumentDefiniton,
			widget.NewLabel("Наименование документа:"), input.Widgets.DocumentName,
			widget.NewLabel("Номер последнего изменения:"), input.Widgets.LastVersionUpdateNumber,
			widget.NewLabel("Обозначение ИУЛ:"), input.Widgets.InfoCertifyingSheet,
			widget.NewLabel("Номер страницы ИУЛ:"), input.Widgets.Page,
			widget.NewLabel("Количество листов ИУЛ:"), input.Widgets.Limit,
		),
	)
}

func textToInt(text string) int {
	// Преобразует строку в целое число, если возможно
	val, err := strconv.Atoi(text)
	if err != nil {
		return 0 // Возвращаем 0, если ошибка преобразования
	}
	return val
}

func isNumeric(text string) bool {
	for _, char := range text {
		if !unicode.IsNumber(char) {
			return false
		}
	}
	return true
}
