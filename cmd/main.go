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

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Тут будет тайтл аля название приложение")
	myWindow.Resize(fyne.NewSize(750, 450))

	fileList := container.NewVBox()   // Контейнер для списка файлов
	signerList := container.NewVBox() // Контейнер для списка подписантов

	openFileButton := widget.NewButton("Выберите файл", func() {
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

			// Создаем новый элемент для списка файлов
			fileLabel := widget.NewLabel(newFile.Name)

			// Переменная для хранения горизонтального контейнера с меткой и кнопкой
			fileContainer := container.NewHBox()

			// Используем NewButtonWithIcon для создания кнопки с иконкой
			removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				// Удаляем файл из списка
				for i, f := range files {
					if f.Path == newFile.Path {
						files = append(files[:i], files[i+1:]...)
						break
					}
				}
				// Удаляем элемент из контейнера
				fileList.Remove(fileContainer)
				fmt.Println("Текущие файлы:", files)
			})

			// Добавляем метку и кнопку в контейнер
			fileContainer.Add(fileLabel)
			fileContainer.Add(removeButton)

			// Добавляем контейнер с файлом в список
			fileList.Add(fileContainer)

			fmt.Println("Выбранные файлы:", files)
		}
	})

	templateOptions := []string{"Шаблон 1", "Шаблон 2"}
	templateSelect := widget.NewSelect(templateOptions, func(value string) {
		selectedTemplate = value
		fmt.Println("Выбранный шаблон:", selectedTemplate)
	})

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Введите ID сотрудников через пробел")

	confirmWorkersButton := widget.NewButton("Подтвердить подписантов", func() {
		ids := strings.Split(idEntry.Text, " ")
		fmt.Println("IDS", ids)
		users, err := internal.ReadDataFromExcel(ids)
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			log.Println("Ошибка при чтении данных из Excel:", err)
		} else {
			for _, user := range users {
				// Создаем новый элемент для списка подписантов
				userLabel := widget.NewLabel(user.WorkType + " " + user.FullName)

				// Переменная для хранения горизонтального контейнера с меткой и кнопкой
				signerContainer := container.NewHBox()

				// Используем NewButtonWithIcon для создания кнопки с иконкой
				removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
					// Удаляем подписанта из списка
					for i, u := range users {
						if u.FullName == user.FullName && u.WorkType == user.WorkType {
							users = append(users[:i], users[i+1:]...)
							break
						}
					}
					// Удаляем элемент из контейнера
					signerList.Remove(signerContainer)
					fmt.Println("Текущие подписанты:", users)
				})

				// Добавляем метку и кнопку в контейнер
				signerContainer.Add(userLabel)
				signerContainer.Add(removeButton)

				// Добавляем контейнер с подписантом в список
				signerList.Add(signerContainer)
			}
		}
	})

	content := container.NewVBox(
		openFileButton,
		fileList,
		widget.NewLabel("Выберите шаблон:"),
		templateSelect,
		widget.NewLabel("Введите IDшники:"),
		idEntry,
		confirmWorkersButton,
		signerList,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
