package internal

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func ReadDataFromExcel(ids []string) []User {
	var users []User
	file, err := excelize.OpenFile("workers.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		if !compareID(row[0], ids) {
			continue
		}
		var user User
		for i, col := range row {
			switch i {
			case 0:
				user.ID = col
			case 1:
				user.WorkType = col
			case 2:
				user.FullName = col
			case 3:
				user.Signature = col
			}
		}
		fmt.Println(user)
		users = append(users, user)
	}
	return users
}

func compareID(id string, ids []string) bool {
	for _, val := range ids {
		if id == val {
			return true
		}
	}
	return false
}
