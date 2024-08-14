package main

import (
	"UserPDFMaker/internal"
	"bufio"
	"fmt"
	"os"
)

func main() {
	// db, err := initDB()
	// if err != nil {
	// 	fmt.Println("Ошибка подключения к базе данных:", err)
	// 	return
	// }

	// Создание нового reader для чтения ввода с клавиатуры
	reader := bufio.NewReader(os.Stdin)

	// Запрос ФИО
	fmt.Print("Введите ФИО: ")
	fullName, _ := reader.ReadString('\n')
	fullName = fullName[:len(fullName)-1] // Удаление символа новой строки

	// Запрос Типа работы
	fmt.Print("Введите тип работы: ")
	workType, _ := reader.ReadString('\n')
	workType = workType[:len(workType)-1] // Удаление символа новой строки

	// Запрос пути к файлу подписи (необязательный)
	/*fmt.Print("Введите путь к файлу подписи (или нажмите Enter, чтобы пропустить): ")
	signature, _ := reader.ReadString('\n')
	signature = signature[:len(signature)-1] // Удаление символа новой строки*/
	var user = internal.User{
		FullName:  fullName,
		WorkType:  workType,
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
