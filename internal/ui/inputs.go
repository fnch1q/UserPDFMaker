package ui

import (
	"UserPDFMaker/internal/data"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateTemplateGroup(input *data.Input) *fyne.Container {
	// Убедитесь, что TemplateSelect создан
	if input.Widgets.TemplateSelect == nil {
		templateOptions := []string{"Шаблон 1", "Шаблон 2"}
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
		input.ObjectName = text
		fmt.Println(input)
	}

	input.Widgets.SerialNumber.SetPlaceHolder("Порядковый номер документа")
	input.Widgets.SerialNumber.OnChanged = func(text string) {
		input.SerialNumber = text
		fmt.Println(input)
	}

	input.Widgets.DocumentDefiniton.SetPlaceHolder("Обозначение ДЭ")
	input.Widgets.DocumentDefiniton.OnChanged = func(text string) {
		input.DocumentDefiniton = text
		fmt.Println(input)
	}

	input.Widgets.DocumentName.SetPlaceHolder("Наименование документа")
	input.Widgets.DocumentName.OnChanged = func(text string) {
		input.DocumentName = text
		fmt.Println(input)
	}

	input.Widgets.LastVersionUpdateNumber.SetPlaceHolder("Номер последнего изменения")
	input.Widgets.LastVersionUpdateNumber.OnChanged = func(text string) {
		input.LastVersionUpdateNumber = text
		fmt.Println(input)
	}

	input.Widgets.InfoCertifyingSheet.SetPlaceHolder("Обозначение ИУЛ")
	input.Widgets.InfoCertifyingSheet.OnChanged = func(text string) {
		input.InfoCertifyingSheet = text
		fmt.Println(input)
	}

	input.Widgets.Page.SetPlaceHolder("Номер страницы ИУЛ")
	input.Widgets.Page.OnChanged = func(text string) {
		input.Page = textToInt(text)
		fmt.Println(input)
	}

	input.Widgets.Limit.SetPlaceHolder("Количество листов ИУЛ")
	input.Widgets.Limit.OnChanged = func(text string) {
		input.Limit = textToInt(text)
		fmt.Println(input)
	}

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
