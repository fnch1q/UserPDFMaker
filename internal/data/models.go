package data

import "time"

type Input struct {
	Users    []User
	Files    []File
	Template string
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
