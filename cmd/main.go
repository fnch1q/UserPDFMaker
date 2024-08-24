package main

import (
	"UserPDFMaker/internal"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	// "os"
	// "path/filepath"
	// d2 "fyne.io/fyne/v2/dialog"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Выбор файла и шаблона")

	var filePath []string
	var selectedTemplate string
	var ids []string

	// widget.NewButton("custom button", func() {
	// 	d2.NewFileOpen(func(r []fyne.URIReadCloser, err error) {
	// 		if err != nil {
	// 			d2.ShowError(err, myWindow)
	// 			return
	// 		}
	// 		if r != nil {
	// 			uri := r.URI()
	// 			if fileInfo, err := os.Stat(uri.Path()); err == nil {
	// 				if fileInfo.IsDir() {
	// 					// Если это директория, получаем все файлы в ней
	// 					files, err := os.ReadDir(uri.Path())
	// 					if err != nil {
	// 						d2.ShowError(err, myWindow)
	// 						return
	// 					}
	// 					for _, file := range files {
	// 						if !file.IsDir() {
	// 							filePath = append(filePath, filepath.Join(uri.Path(), file.Name()))
	// 						}
	// 					}
	// 				} else {
	// 					// Если это файл, добавляем его путь в filePath
	// 					filePath = append(filePath, uri.Path())
	// 				}
	// 				fmt.Println("Выбранные пути:", filePath)
	// 			}
	// 		}
	// 	}, myWindow)
	// })

	openFileButton := widget.NewButton("Выбрать файл", func() {
		selectedPath, err := dialog.File().Filter("Все файлы", "*").Load()
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка")
			return
		}
		if selectedPath != "" {
			// Добавляем путь в список
			filePath = append(filePath, selectedPath)
			fmt.Println("Выбранный путь:", filePath)
		}
	})

	templateOptions := []string{"Шаблон 1", "Шаблон 2"}
	templateSelect := widget.NewSelect(templateOptions, func(value string) {
		selectedTemplate = value
		fmt.Println("Выбранный шаблон:", selectedTemplate)
	})

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Введите ID (0 для завершения ввода)")
	idEntry.OnSubmitted = func(value string) {
		if value == "0" {
			fmt.Println("Введенные ID:", ids)
		} else {
			ids = append(ids, value)
			idEntry.SetText("")
		}
	}

	stopButton := widget.NewButton("Завершить ввод ID", func() {
		fmt.Println("Введенные ID:", ids)
	})

	content := container.NewVBox(
		openFileButton,
		// openfolderbtton,
		widget.NewLabel("Выберите шаблон:"),
		templateSelect,
		widget.NewLabel("Введите IDшники:"),
		idEntry,
		stopButton,
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(400, 300))
	myWindow.ShowAndRun()

	//для теста экселя
	idsForCheck := make([]string, 0)
	idsForCheck = append(idsForCheck, "1")
	idsForCheck = append(idsForCheck, "2")
	idsForCheck = append(idsForCheck, "6")
	internal.ReadDataFromExcel(idsForCheck)
}
