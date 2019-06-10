package main

import (
	"os"
	"fmt"
	"path"
	"bufio"
	"strings"
	"io"
	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	currentDir, _ := os.Getwd()
	fmt.Println("当前文件路径：", currentDir)
	dirFile, _ := os.Open(currentDir)
	fileInfoList, _ := dirFile.Readdir(0)
	defer dirFile.Close()
	for _, fileInfo := range fileInfoList {
		fileSuffix := path.Ext(fileInfo.Name()) //获取文件后缀
		if (fileSuffix != ".txt" && fileSuffix != "") || fileInfo.Name() == "txt2excel" {
			continue
		}

		filePath := path.Join(currentDir, fileInfo.Name())
		fmt.Println("filePath: ", filePath)
		results, err := readTxt(filePath)
		if err != nil {
			fmt.Printf("read file: '%s' fialed, err: %s\n", filePath, err)
			continue
		}

		toExcel(results)
	}

	os.Exit(0)
}

func toExcel(contentList []string) {
	xlsx := excelize.NewFile()
	var fileName string
	for index, content := range contentList {
		if index == 0 {
			fileName = contentList[index]
			continue
		}

		r := strings.Split(content, ",")
		xlsx.SetCellValue("Sheet1", fmt.Sprintf("A%s", r[0]), r[1])
	}

	fileName = fmt.Sprintf("%s.xlsx", fileName)
	currendDir, _ := os.Getwd()
	filePath := path.Join(currendDir, fileName)
	err := xlsx.SaveAs(filePath)
	if err != nil {
		fmt.Printf("创建excel: %s 失败, err: %s\n", fileName, err)
		return
	}

	fmt.Printf("创建 excel: %s 成功，行数： %d\n", fileName, len(contentList))
}

func readTxt(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	var result []string
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				return result, nil
			}

			return nil, err
		}

		result = append(result, line)
	}

	return result, nil
}
