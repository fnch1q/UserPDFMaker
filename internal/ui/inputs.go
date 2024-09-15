package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func CreateTemplateGroup() *fyne.Container {
	templateOptions := []string{"Шаблон 1", "Шаблон 2"}
	templateSelect := widget.NewSelect(templateOptions, func(value string) {
		input.Template = value
	})

	label := widget.NewLabel("Выберите шаблон:")
	label.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		label,
		templateSelect,
	)
}

func CreateDocumentDetailsGroup() *fyne.Container {
	// Создаем поля для ввода
	objectNameEntry := widget.NewEntry()
	objectNameEntry.SetPlaceHolder("Наименование объекта")
	objectNameEntry.OnChanged = func(text string) {
		input.ObjectName = text
	}

	serialNumberEntry := widget.NewEntry()
	serialNumberEntry.SetPlaceHolder("Порядковый номер документа")
	serialNumberEntry.OnChanged = func(text string) {
		input.SerialNumber = text
	}

	documentDefinitionEntry := widget.NewEntry()
	documentDefinitionEntry.SetPlaceHolder("Обозначение ДЭ")
	documentDefinitionEntry.OnChanged = func(text string) {
		input.DocumentDefiniton = text
	}

	documentNameEntry := widget.NewEntry()
	documentNameEntry.SetPlaceHolder("Наименование документа")
	documentNameEntry.OnChanged = func(text string) {
		input.DocumentName = text
	}

	lastVersionUpdateNumberEntry := widget.NewEntry()
	lastVersionUpdateNumberEntry.SetPlaceHolder("Номер последнего изменения")
	lastVersionUpdateNumberEntry.OnChanged = func(text string) {
		input.LastVersionUpdateNumber = text
	}

	infoCertifyingSheetEntry := widget.NewEntry()
	infoCertifyingSheetEntry.SetPlaceHolder("Обозначение ИУЛ")
	infoCertifyingSheetEntry.OnChanged = func(text string) {
		input.InfoCertifyingSheet = text
	}

	pageEntry := widget.NewEntry()
	pageEntry.SetPlaceHolder("Номер страницы ИУЛ")
	pageEntry.OnChanged = func(text string) {
		input.Page = textToInt(text)
	}

	limitEntry := widget.NewEntry()
	limitEntry.SetPlaceHolder("Количество листов ИУЛ")
	limitEntry.OnChanged = func(text string) {
		input.Limit = textToInt(text)
	}

	// Заголовок группы
	label := widget.NewLabel("Детали документа")
	label.TextStyle = fyne.TextStyle{Bold: true}

	// Собираем все в контейнер с использованием Grid
	return container.NewVBox(
		label,
		container.NewGridWithColumns(2,
			widget.NewLabel("Наименование объекта:"), objectNameEntry,
			widget.NewLabel("Порядковый номер документа:"), serialNumberEntry,
			widget.NewLabel("Обозначение ДЭ:"), documentDefinitionEntry,
			widget.NewLabel("Наименование документа:"), documentNameEntry,
			widget.NewLabel("Номер последнего изменения:"), lastVersionUpdateNumberEntry,
			widget.NewLabel("Обозначение ИУЛ:"), infoCertifyingSheetEntry,
			widget.NewLabel("Номер страницы ИУЛ:"), pageEntry,
			widget.NewLabel("Количество листов ИУЛ:"), limitEntry,
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
