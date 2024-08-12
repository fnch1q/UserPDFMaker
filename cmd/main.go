package main

import (
	"UserPDFMaker/internal"
	"fmt"
)

func main() {
	// db, err := initDB()
	// if err != nil {
	// 	fmt.Println("Ошибка подключения к базе данных:", err)
	// 	return
	// }

	// Пример получения пользователя с ID=1
	var user = internal.User{
		FullName:  "Reshetnikov Michil",
		WorkType:  "WorkType",
		Signature: "C:\\Labis\\Projects\\UserPDFMaker\\images\\image.png",
	}

	// db.First(&user, 1)

	// Генерация PDF
	err := internal.GeneratePDF(user)
	if err != nil {
		fmt.Println("Ошибка генерации PDF:", err)
		return
	}

	fmt.Println("PDF файл успешно создан!")
}
