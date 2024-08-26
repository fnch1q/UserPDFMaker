package data

import "time"

type User struct {
	ID        string
	WorkType  string
	FullName  string
	Signature string // путь к файлу с изображением подписи
}

type File struct {
	ID         string
	Path       string
	Name       string
	Hash       string
	UpdateTime time.Time
	Size       int64
}
