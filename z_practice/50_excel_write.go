package main

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"fmt"
	"path"
	"os"
)

//写excel

func main() {
	xlsx := excelize.NewFile()

	index := xlsx.NewSheet("Sheet2") //创建一个新sheet

	xlsx.SetCellValue("Sheet2", "A2", "Hello world") //设置单元格值
	xlsx.SetCellValue("Sheet1", "B2", 100)

	xlsx.SetActiveSheet(index)

	currendDir, _ := os.Getwd()
	filePath := path.Join(currendDir, "Book1.xlsx")
	fmt.Println("filePath:", filePath)
	err := xlsx.SaveAs(filePath)

	if err != nil {
		fmt.Println("save excel err:", err)
	}

}
