package main

import (
	"UserPDFMaker/internal"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

var selectedTemplate string
var files []internal.File
var users []internal.User // Массив подписантов

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("PDF Maker")
	myWindow.Resize(fyne.NewSize(800, 500))

	// Контейнеры для файлов и подписантов
	fileList := container.NewVBox()
	signerList := container.NewVBox()

	// Кнопка для выбора файлов
	openFileButton := widget.NewButtonWithIcon("Выберите файл", theme.FolderIcon(), func() {
		selectedPath, err := dialog.File().Filter("Все файлы", "*").Load()
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			return
		}
		if selectedPath != "" {
			newFile, err := internal.NewFile(selectedPath)
			if err != nil {
				dialog.Message("%s", err.Error()).Title("Ошибка").Error()
				return
			}

			files = append(files, *newFile)

			// Добавляем файл в интерфейс
			fileLabel := widget.NewLabel(newFile.Name)
			fileContainer := container.NewHBox()
			removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				// Удаление файла из списка
				for i, f := range files {
					if f.Path == newFile.Path {
						files = append(files[:i], files[i+1:]...)
						break
					}
				}
				fileList.Remove(fileContainer)
				fmt.Println("Текущие файлы:", files)
			})

			fileContainer.Add(fileLabel)
			fileContainer.Add(removeButton)
			fileList.Add(fileContainer)

			fmt.Println("Выбранные файлы:", files)
		}
	})

	// Выпадающий список шаблонов
	templateOptions := []string{"Шаблон 1", "Шаблон 2"}
	templateSelect := widget.NewSelect(templateOptions, func(value string) {
		selectedTemplate = value
		fmt.Println("Выбранный шаблон:", selectedTemplate)
	})

	// Поле для ввода ID сотрудников
	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Введите ID сотрудников через пробел")

	// Кнопка для подтверждения подписантов
	confirmWorkersButton := widget.NewButtonWithIcon("Добавить подписантов", theme.ConfirmIcon(), func() {
		ids := strings.Split(idEntry.Text, " ")

		// Чтение подписантов и добавление их в массив
		newUsers, err := internal.ReadDataFromExcel(ids)
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			log.Println("Ошибка при чтении данных из Excel:", err)
		} else {
			for _, newUser := range newUsers {
				found := false
				for _, existingUser := range users {
					if existingUser.FullName == newUser.FullName && existingUser.WorkType == newUser.WorkType {
						found = true
						break
					}
				}
				if !found {
					users = append(users, newUser)
				}
			}

			// Обновляем список подписантов в интерфейсе
			signerList.RemoveAll() // Очищаем контейнер перед добавлением
			for _, user := range users {
				userLabel := widget.NewLabel(user.WorkType + " " + user.FullName)
				signerContainer := container.NewHBox()

				removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
					for i, u := range users {
						if u.FullName == user.FullName && u.WorkType == user.WorkType {
							users = append(users[:i], users[i+1:]...)
							break
						}
					}
					signerList.Remove(signerContainer)
					fmt.Println("Текущие подписанты:", users)
				})

				signerContainer.Add(userLabel)
				signerContainer.Add(removeButton)
				signerList.Add(signerContainer)
			}

			fmt.Println("Текущие подписанты:", users)
		}
	})

	// Организация группировки элементов интерфейса
	fileGroup := container.NewVBox(
		widget.NewLabel("Файлы:"),
		fileList,
		openFileButton,
	)

	templateGroup := container.NewVBox(
		widget.NewLabel("Выберите шаблон:"),
		templateSelect,
	)

	signerGroup := container.NewVBox(
		widget.NewLabel("Подписанты:"),
		signerList,
		idEntry,
		confirmWorkersButton,
	)

	content := container.NewVBox(
		fileGroup,
		templateGroup,
		signerGroup,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
