package ui

import (
	"UserPDFMaker/internal/data"
	"UserPDFMaker/internal/utils"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var files []data.File
var users []data.User
var selectedTemplate string

func CreateMainWindow(app fyne.App) fyne.Window {
	window := app.NewWindow("PDF Maker")
	window.Resize(fyne.NewSize(800, 600))

	icon, _ := fyne.LoadResourceFromPath("../images/app_icon.png")
	window.SetIcon(icon)

	// Создание и добавление элементов интерфейса
	fileGroup := CreateFileGroup()
	templateGroup := CreateTemplateGroup()
	signerGroup := CreateSignerGroup()
	generatePdfButton := widget.NewButtonWithIcon("Сгенерировать PDF", theme.FileTextIcon(), func() {
		if len(users) == 0 {
			fmt.Println("Нет пользователей для генерации PDF")
			return
		}

		// Выбираем первого пользователя для примера
		user := users[0]

		// Вызываем функцию для генерации PDF
		err := utils.GeneratePDF(user)
		if err != nil {
			fmt.Println("Ошибка при генерации PDF:", err)
		} else {
			fmt.Println("PDF успешно сгенерирован")
		}
	})

	content := container.NewVBox(
		templateGroup,
		fileGroup,
		signerGroup,
		generatePdfButton,
	)

	window.SetContent(content)
	_ = selectedTemplate
	return window
}
