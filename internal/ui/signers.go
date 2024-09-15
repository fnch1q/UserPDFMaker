package ui

import (
	"UserPDFMaker/internal/data"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func CreateSignerGroup() *fyne.Container {
	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Введите ID сотрудников через пробел")

	signerList := container.NewVBox() // Создаем контейнер для списка подписантов

	confirmWorkersButton := widget.NewButtonWithIcon("Добавить подписантов", theme.ConfirmIcon(), func() {
		ids := strings.Split(idEntry.Text, " ")
		newUsers, err := data.ReadDataFromExcel(ids)
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			log.Println("Ошибка при чтении данных из Excel:", err)
		} else {
			addUsersToList(newUsers, signerList) // Передаем список в функцию
		}
	})

	label := widget.NewLabel("Подписанты:")
	label.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		label,
		signerList, // Включаем контейнер для списка подписантов
		idEntry,
		confirmWorkersButton,
	)
}

func addUsersToList(newUsers []data.User, signerList *fyne.Container) {
	for _, newUser := range newUsers {
		found := false
		for _, existingUser := range input.Users {
			if existingUser.FullName == newUser.FullName && existingUser.WorkType == newUser.WorkType {
				found = true
				break
			}
		}
		if !found {
			input.Users = append(input.Users, newUser)
		}
	}
	updateSignerList(signerList)
}

func updateSignerList(signerList *fyne.Container) {
	signerList.RemoveAll() // Очищаем контейнер перед обновлением

	for _, user := range input.Users {
		userLabel := widget.NewLabel(user.WorkType + " " + user.FullName)
		signerContainer := container.NewHBox()

		removeButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			for i, u := range input.Users {
				if u.FullName == user.FullName && u.WorkType == user.WorkType {
					input.Users = append(input.Users[:i], input.Users[i+1:]...)
					break
				}
			}
			updateSignerList(signerList) // Обновляем список подписантов после удаления
		})

		signerContainer.Add(userLabel)
		signerContainer.Add(removeButton)
		signerList.Add(signerContainer)
	}

	signerList.Refresh() // Обновляем отображение контейнера на экране
}
