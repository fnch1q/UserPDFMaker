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
	signerGroup := CreateSignerGroup(&input)
	documentDetailsGroup := CreateDocumentDetailsGroup(&input)

	// Генерация пдф (эндпоинт)
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
			resetInput(&input)
		}
	})

	// Создаем вкладки
	settingsTab := container.NewVBox(
		templateGroup,
		documentDetailsGroup,
	)
	filesTab := container.NewVBox(
		fileGroup,
		signerGroup,
		generatePdfButton,
	)
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

	window.SetContent(tabs)
	return window
}

func resetInput(input *data.Input) {
	input.Users = nil
	input.Files = nil
	input.ObjectName = ""
	input.SerialNumber = ""
	input.DocumentDefiniton = ""
	input.DocumentName = ""
	input.LastVersionUpdateNumber = ""
	input.InfoCertifyingSheet = ""
	input.Page = 0
	input.Limit = 0

	// Сбрасываем виджеты
	input.Widgets.IDEntry.SetText("")
	input.Widgets.ObjectName.SetText("")
	input.Widgets.SerialNumber.SetText("")
	input.Widgets.DocumentDefiniton.SetText("")
	input.Widgets.DocumentName.SetText("")
	input.Widgets.LastVersionUpdateNumber.SetText("")
	input.Widgets.InfoCertifyingSheet.SetText("")
	input.Widgets.Page.SetText("")
	input.Widgets.Limit.SetText("")

	// Если контейнер для файлов существует, очищаем его
	if input.FileList != nil {
		input.FileList.RemoveAll()
	}

	if input.SignerList != nil {
		input.SignerList.RemoveAll()
	}
}
