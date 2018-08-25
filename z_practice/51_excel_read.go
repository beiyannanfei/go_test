package main

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"path"
	"os"
	"fmt"
)

//读excels
//	cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 51_excel_read.go
func main() {
	currendDir, _ := os.Getwd()
	filePath := path.Join(currendDir, "Book2.xlsx")
	fmt.Println("filePath", filePath)
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println("openfile err:", err)
		return
	}

	cellValue_sheet1_b2 := xlsx.GetCellValue("Sheet1", "B2") //获取单元格的值
	fmt.Println("cellValue_sheet1_b2:", cellValue_sheet1_b2)

	allSheet := xlsx.GetSheetMap()
	fmt.Printf("allSheet: %v\n", allSheet)

	rows := xlsx.GetRows(allSheet[1])
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
