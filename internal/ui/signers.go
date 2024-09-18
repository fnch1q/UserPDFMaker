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

func CreateSignerGroup(input *data.Input) *fyne.Container {
	// Создаем контейнер для списка подписантов только если он не создан ранее
	if input.Widgets.SignerList == nil {
		input.Widgets.SignerList = container.NewVBox() // Создаем контейнер для списка подписантов
	}
	if input.Widgets.IDEntry == nil {
		input.Widgets.IDEntry = widget.NewEntry()
	}
	input.Widgets.IDEntry.SetPlaceHolder("Введите ID сотрудников через пробел")

	confirmWorkersButton := widget.NewButtonWithIcon("Добавить подписантов", theme.ConfirmIcon(), func() {
		ids := strings.Split(input.Widgets.IDEntry.Text, " ")
		newUsers, err := data.ReadDataFromExcel(ids)
		if err != nil {
			dialog.Message("%s", err.Error()).Title("Ошибка").Error()
			log.Println("Ошибка при чтении данных из Excel:", err)
		} else {
			addUsersToList(input, newUsers) // Передаем список в функцию
		}
	})

	label := widget.NewLabel("Подписанты:")
	label.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		label,
		input.Widgets.SignerList, // Включаем контейнер для списка подписантов
		input.Widgets.IDEntry,
		confirmWorkersButton,
	)
}

func addUsersToList(input *data.Input, newUsers []data.User) {
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
	updateSignerList(input) // Обновляем список подписантов после добавления
}

func updateSignerList(input *data.Input) {
	input.Widgets.SignerList.RemoveAll() // Очищаем контейнер перед обновлением

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
			updateSignerList(input) // Обновляем список подписантов после удаления
		})

		signerContainer.Add(userLabel)
		signerContainer.Add(removeButton)
		input.Widgets.SignerList.Add(signerContainer)
	}

	input.Widgets.SignerList.Refresh() // Обновляем отображение контейнера на экране
}
