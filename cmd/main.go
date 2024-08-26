package main

import (
	"UserPDFMaker/internal/ui"

	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	myWindow := ui.CreateMainWindow(myApp)
	myWindow.ShowAndRun()
}
