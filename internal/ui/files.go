package ui

import (
	"UserPDFMaker/internal/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func CreateFileGroup() *fyne.Container {
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

			files = append(files, *newFile)
			addFileToList(newFile, fileList)
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

func addFileToList(newFile *data.File, fileList *fyne.Container) {
	fileLabel := widget.NewLabel(newFile.Name)
	fileContainer := container.NewHBox()
	removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		for i, f := range files {
			if f.Path == newFile.Path {
				files = append(files[:i], files[i+1:]...)
				break
			}
		}
		fileList.Remove(fileContainer)
	})

	fileContainer.Add(fileLabel)
	fileContainer.Add(removeButton)
	fileList.Add(fileContainer)
}
