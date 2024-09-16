package ui

import (
	"UserPDFMaker/internal/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func CreateFileGroup(input *data.Input) *fyne.Container {
	fileList := container.NewVBox()
	openFileButton := widget.NewButtonWithIcon("Добавить файл", theme.FolderIcon(), func() {
		selectedPath, err := dialog.File().Filter("Все файлы", "*").Load()
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			return
		}
		if selectedPath != "" {
			newFile, err := data.NewFile(selectedPath)
			if err != nil {
				dialog.Message("%s", err.Error()).Title("Ошибка").Error()
				return
			}

			input.Files = append(input.Files, *newFile)
			addFileToList(newFile, input, fileList)
		}
	})

	label := widget.NewLabel("Файлы:")
	label.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		label,
		fileList,
		openFileButton,
	)
}

func addFileToList(newFile *data.File, input *data.Input, fileList *fyne.Container) {
	fileLabel := widget.NewLabel(newFile.Name)
	fileContainer := container.NewHBox()
	removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		for i, f := range input.Files {
			if f.Path == newFile.Path {
				input.Files = append(input.Files[:i], input.Files[i+1:]...)
				break
			}
		}
		fileList.Remove(fileContainer)
	})

	fileContainer.Add(fileLabel)
	fileContainer.Add(removeButton)
	fileList.Add(fileContainer)
}
