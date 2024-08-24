package internal

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func ReadDataFromExcel(ids []string) {
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
		for _, col := range row {
			fmt.Print(col, "\t")
		}
		fmt.Println()
	}
}

func compareID(id string, ids []string) bool {
	for _, val := range ids {
		if id == val {
			return true
		}
	}
	return false
}
