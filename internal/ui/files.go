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
	// Инициализация контейнера для файлов
	input.Widgets.FileList = container.NewVBox()

	// Кнопка для добавления файла
	openFileButton := widget.NewButtonWithIcon("Добавить файл", theme.FolderIcon(), func() {
		// Проверка на выбранный шаблон
		if input.Template == "Шаблон 2" && len(input.Files) >= 1 {
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

	// Заголовок группы
	label := widget.NewLabel("Файлы:")
	label.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		label,
		input.Widgets.FileList, // Используем input.Widgets.FileList
		openFileButton,
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
