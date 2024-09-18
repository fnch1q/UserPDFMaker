package data

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Input struct {
	Users    []User
	Files    []File
	Template string
	DocumentDetails
	Widgets
}

type User struct {
	ID        string
	WorkType  string
	FullName  string
	Signature string // путь к файлу с изображением подписи
}

type File struct {
	Path       string
	Name       string
	Hash       string
	UpdateTime time.Time
	Size       int64
}

type DocumentDetails struct {
	ObjectName              string
	SerialNumber            string
	DocumentDefiniton       string
	DocumentName            string
	LastVersionUpdateNumber string
	InfoCertifyingSheet     string
	Page                    int
	Limit                   int
}

type Widgets struct {
	TemplateSelect *widget.Select
	FileList       *fyne.Container
	SignerList     *fyne.Container
	IDEntry        *widget.Entry
	DocumentDetailWidgets
}

type DocumentDetailWidgets struct {
	ObjectName              *widget.Entry
	SerialNumber            *widget.Entry
	DocumentDefiniton       *widget.Entry
	DocumentName            *widget.Entry
	LastVersionUpdateNumber *widget.Entry
	InfoCertifyingSheet     *widget.Entry
	Page                    *widget.Entry
	Limit                   *widget.Entry
}
