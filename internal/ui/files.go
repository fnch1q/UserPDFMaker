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
	input.FileList = fileList
	openFileButton := widget.NewButtonWithIcon("Добавить файл", theme.FolderIcon(), func() {
		// Проверка на выбранный шаблон
		if input.Template == "Шаблон 2" && len(input.Files) >= 1 {
			dialog.Message("Вы не можете выбрать более одного файла для этого шаблона").Title("Ограничение").Error()
			return
		}

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
			addFileToList(newFile, input)
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

func addFileToList(newFile *data.File, input *data.Input) {
	fileLabel := widget.NewLabel(newFile.Name)
	fileContainer := container.NewHBox()
	removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		for i, f := range input.Files {
			if f.Path == newFile.Path {
				input.Files = append(input.Files[:i], input.Files[i+1:]...)
				break
			}
		}
		input.FileList.Remove(fileContainer)
	})

	fileContainer.Add(fileLabel)
	fileContainer.Add(removeButton)
	input.FileList.Add(fileContainer)
}
