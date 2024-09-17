package ui

import (
	"UserPDFMaker/internal/data"
	"UserPDFMaker/internal/utils"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func CreateMainWindow(app fyne.App) fyne.Window {
	window := app.NewWindow("PDF Maker")
	window.Resize(fyne.NewSize(800, 600))

	icon, _ := fyne.LoadResourceFromPath("../images/app_icon.png")
	window.SetIcon(icon)

	var input data.Input

	// Создание и добавление элементов интерфейса
	fileGroup := CreateFileGroup(&input)
	templateGroup := CreateTemplateGroup(&input)
	documentDetailsGroup := CreateDocumentDetailsGroup(&input)
	signerGroup := CreateSignerGroup(&input)
	generatePdfButton := widget.NewButtonWithIcon("Сгенерировать PDF", theme.FileTextIcon(), func() {
		if len(input.Users) == 0 {
			fmt.Println("Нет пользователей для генерации PDF")
			return
		}

		// Выбираем первого пользователя для примера
		user := input.Users[0]

		// Вызываем функцию для генерации PDF
		err := utils.GeneratePDF(user)
		if err != nil {
			fmt.Println("Ошибка при генерации PDF:", err)
		} else {
			fmt.Println("PDF успешно сгенерирован")
		}
	})

	settingsTab := container.NewVBox(
		templateGroup,
		documentDetailsGroup,
	)

	// Вкладка для подписантов и файлов
	filesTab := container.NewVBox(
		fileGroup,
		signerGroup,
		generatePdfButton,
	)

	// Создаем табы
	tabs := container.NewAppTabs(
		container.NewTabItem("Настройки", settingsTab),
		container.NewTabItem("Файлы и подписанты", filesTab),
	)

	// Добавляем обработчик выбора вкладки
	tabs.OnSelected = func(tab *container.TabItem) {
		if tab.Text == "Файлы и подписанты" && input.Template == "" {
			dialog.Message("Пожалуйста, выберите шаблон перед переходом на эту вкладку").Title("Выберите шаблон").Error()
			// Возвращаемся на вкладку "Настройки"
			tabs.SelectTabIndex(0)
		}
	}

	// Устанавливаем содержимое окна
	window.SetContent(tabs)

	return window
}
