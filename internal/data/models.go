package data

import (
	"time"

	"fyne.io/fyne/v2"
)

type Input struct {
	Users                   []User
	Files                   []File
	Template                string
	ObjectName              string //Наименование объекта
	SerialNumber            string //Порядковый номер документа
	DocumentDefiniton       string //Обозначение ДЭ
	DocumentName            string //Наименование документа
	LastVersionUpdateNumber string //Номер последнего изменения
	InfoCertifyingSheet     string //Обозначение ИУЛ
	Page                    int    //Номер страницы ИУЛ
	Limit                   int    //Количество листов ИУЛ
	FileList                *fyne.Container
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
