package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Создаем приложение
	myApp := app.New()
	myWindow := myApp.NewWindow("Выбор файла и шаблона")

	// Переменные для хранения пути и ID
	var filePath string
	var selectedTemplate string
	var ids []string

	// Создаем кнопку для открытия диалогового окна
	openFileButton := widget.NewButton("Выбрать файл или директорию", func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if uc != nil {
				filePath = uc.URI().Path()
				uc.Close()
			}
		}, myWindow)
	})

	// Выпадающий список для выбора шаблона
	templateOptions := []string{"Шаблон 1", "Шаблон 2"}
	templateSelect := widget.NewSelect(templateOptions, func(value string) {
		selectedTemplate = value
	})

	_ = filePath
	_ = selectedTemplate

	// Текстовое поле для ввода IDшников
	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Введите ID (0 для завершения ввода)")
	idEntry.OnSubmitted = func(value string) {
		if value == "0" {
			// Завершаем ввод ID
			dialog.ShowInformation("Информация", "Ввод ID завершен", myWindow)
		} else {
			ids = append(ids, value)
			idEntry.SetText("") // очищаем поле для следующего ввода
		}
	}

	// Компонуем элементы на экране
	content := container.NewVBox(
		openFileButton,
		widget.NewLabel("Выберите шаблон:"),
		templateSelect,
		widget.NewLabel("Введите IDшники:"),
		idEntry,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.ShowAndRun()
}
