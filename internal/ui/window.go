package ui

import (
	"UserPDFMaker/internal/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

var files []data.File
var users []data.User
var selectedTemplate string

func CreateMainWindow(app fyne.App) fyne.Window {
	window := app.NewWindow("PDF Maker")
	window.Resize(fyne.NewSize(800, 500))

	// Создание и добавление элементов интерфейса
	fileGroup := CreateFileGroup()
	templateGroup := CreateTemplateGroup()
	signerGroup := CreateSignerGroup()

	content := container.NewVBox(
		fileGroup,
		templateGroup,
		signerGroup,
	)

	window.SetContent(content)
	_ = selectedTemplate
	return window
}
