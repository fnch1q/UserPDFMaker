package ui

import (
	"UserPDFMaker/internal/data"
	"io/ioutil"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func CreateFileGroup(input *data.Input) *fyne.Container {
	// Инициализация контейнера для файлов
	input.Widgets.FileList = container.NewVBox()

	// Кнопка для добавления файла
	openFileButton := widget.NewButtonWithIcon("Добавить файл", theme.FileIcon(), func() {
		// Проверка на выбранный шаблон
		if input.Template == template1 && len(input.Files) >= 1 {
			dialog.Message("Вы не можете выбрать более одного файла для этого шаблона").Title("Ограничение").Error()
			return
		}

		// Открытие диалогового окна для выбора файла
		selectedPath, err := dialog.File().Filter("Все файлы", "*").Load()
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			return
		}

		if selectedPath != "" {
			// Создание нового файла
			newFile, err := data.NewFile(selectedPath)
			if err != nil {
				dialog.Message("%s", err.Error()).Title("Ошибка").Error()
				return
			}

			// Добавление файла в список и отображение
			input.Files = append(input.Files, *newFile)
			addFileToList(newFile, input)
		}
	})

	openFolderButton := widget.NewButtonWithIcon("Добавить папку", theme.FolderIcon(), func() {
		// Открытие диалогового окна для выбора папки
		if input.Template == template1 {
			dialog.Message("Вы не можете выбрать более одного файла для этого шаблона").Title("Ограничение").Error()
			return
		}
		selectedPath, err := dialog.Directory().Browse()
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			return
		}

		if selectedPath != "" {
			// Добавление всех файлов в папке в список
			files, err := ioutil.ReadDir(selectedPath)
			if err != nil {
				dialog.Message("%s", err.Error()).Title("Ошибка").Error()
				return
			}

			for _, file := range files {
				if !file.IsDir() {
					newFile, err := data.NewFile(filepath.Join(selectedPath, file.Name()))
					if err != nil {
						dialog.Message("%s", err.Error()).Title("Ошибка").Error()
						return
					}

					input.Files = append(input.Files, *newFile)
					addFileToList(newFile, input)
				}
			}
		}
	})

	label := widget.NewLabel("Файлы:")
	label.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		label,
		input.Widgets.FileList,
		openFileButton,
		openFolderButton,
	)
}

func addFileToList(newFile *data.File, input *data.Input) {
	// Создание элементов для отображения файла
	fileLabel := widget.NewLabel(newFile.Name)
	fileContainer := container.NewHBox()
	removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		// Удаление файла из списка
		for i, f := range input.Files {
			if f.Path == newFile.Path {
				input.Files = append(input.Files[:i], input.Files[i+1:]...)
				break
			}
		}
		input.Widgets.FileList.Remove(fileContainer)
	})

	// Добавление элементов в контейнер
	fileContainer.Add(fileLabel)
	fileContainer.Add(removeButton)
	input.Widgets.FileList.Add(fileContainer)
}
