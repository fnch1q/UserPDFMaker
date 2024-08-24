package main

import (
	"UserPDFMaker/internal"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

var filePathes []string
var filePathesString string = "Выбранные файлы:"
var selectedTemplate string

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Тут будет тайтл аля название приложение")

	pathes := widget.NewLabel(filePathesString)
	openFileButton := widget.NewButton("Выберите файл", func() {
		selectedPath, err := dialog.File().Filter("Все файлы", "*").Load()
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка")
			return
		}
		if selectedPath != "" {
			filePathes = append(filePathes, selectedPath)
			filePathesString += "\n" + selectedPath
			pathes.SetText(filePathesString)
			fmt.Println("Выбранный путь:", filePathes)
		}
	})

	templateOptions := []string{"Шаблон 1", "Шаблон 2"}
	templateSelect := widget.NewSelect(templateOptions, func(value string) {
		selectedTemplate = value
		fmt.Println("Выбранный шаблон:", selectedTemplate)
	})

	names := widget.NewLabel("Выбранные подписанты:")
	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Введите ID сотрудников через пробел")

	confirmWorkersButton := widget.NewButton("Подтвердить подписантов", func() {
		var usersName string = "Выбранные подписанты:"
		ids := strings.Split(idEntry.Text, " ")
		fmt.Println("IDS", ids)
		users := internal.ReadDataFromExcel(ids)
		for _, user := range users {
			usersName += "\n" + user.WorkType + " " + user.FullName
		}
		names.SetText(usersName)
	})

	content := container.NewVBox(
		openFileButton,
		pathes,
		widget.NewLabel("Выберите шаблон:"),
		templateSelect,
		widget.NewLabel("Введите IDшники:"),
		idEntry,
		confirmWorkersButton,
		names,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(750, 450))
	myWindow.ShowAndRun()
}
