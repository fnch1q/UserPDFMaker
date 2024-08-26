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

	return container.NewVBox(
		widget.NewLabel("Подписанты:"),
		signerList, // Включаем контейнер для списка подписантов
		idEntry,
		confirmWorkersButton,
	)
}

func addUsersToList(newUsers []data.User, signerList *fyne.Container) {
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
	updateSignerList(signerList)
}

func updateSignerList(signerList *fyne.Container) {
	signerList.RemoveAll() // Очищаем контейнер перед обновлением

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
			updateSignerList(signerList) // Обновляем список подписантов после удаления
		})

		signerContainer.Add(userLabel)
		signerContainer.Add(removeButton)
		signerList.Add(signerContainer)
	}

	signerList.Refresh() // Обновляем отображение контейнера на экране
}
