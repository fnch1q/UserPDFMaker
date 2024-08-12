package internal

type User struct {
	ID        uint `gorm:"primaryKey"`
	FullName  string
	WorkType  string
	Signature string // путь к файлу с изображением подписи
}
