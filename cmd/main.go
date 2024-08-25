package main

import (
	"UserPDFMaker/internal"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

var fileNames strings.Builder
var selectedTemplate string
var files []internal.File // Массив структур File

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Тут будет тайтл аля название приложение")
	myWindow.Resize(fyne.NewSize(750, 450))

	// Инициализация строки с выбранными файлами
	fileNames.WriteString("Выбранные файлы:")
	pathes := widget.NewLabel(fileNames.String())

	openFileButton := widget.NewButton("Выберите файл", func() {
		selectedPath, err := dialog.File().Filter("Все файлы", "*").Load()
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			return
		}
		if selectedPath != "" {
			// Создаем новый объект File и добавляем его в массив
			newFile, err := internal.NewFile(selectedPath)
			if err != nil {
				dialog.Message("%s", err.Error()).Title("Ошибка").Error()
				return
			}

			files = append(files, *newFile) // Добавляем в массив файлов
			fileNames.WriteString("\n" + newFile.Name)
			pathes.SetText(fileNames.String())
			fmt.Println("Выбранные файлы:", files)
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
		var usersName strings.Builder
		usersName.WriteString("Выбранные подписанты:")
		ids := strings.Split(idEntry.Text, " ")
		fmt.Println("IDS", ids)
		users, err := internal.ReadDataFromExcel(ids)
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			log.Println("Ошибка при чтении данных из Excel:", err)
		} else {
			for _, user := range users {
				usersName.WriteString("\n" + user.WorkType + " " + user.FullName)
			}
		}
		names.SetText(usersName.String())
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
	myWindow.ShowAndRun()
}
