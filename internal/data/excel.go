package data

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

var (
	ErrExcelFileNotExists = errors.New("err_when_open_file")
	ErrWhenReadingRows    = errors.New("err_when_reading_rows")
)

func ReadDataFromExcel(ids []string) ([]User, error) {
	var users []User
	file, err := excelize.OpenFile("workers.xlsx")
	if err != nil {
		return users, ErrExcelFileNotExists
	}

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return users, ErrWhenReadingRows
	}

	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	for _, row := range rows {
		if !compareID(row[0], idMap) {
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
	return users, nil
}

func compareID(id string, ids map[string]bool) bool {
	return ids[id]
}
