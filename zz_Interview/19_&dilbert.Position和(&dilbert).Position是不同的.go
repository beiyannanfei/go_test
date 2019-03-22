package main

import (
	"time"
	"fmt"
)

func main() {
	type Employee struct {
		ID        int
		Name      string
		Address   string
		DoB       time.Time
		Position  string
		Salary    int
		ManagerID int
	}

	var dilbert Employee

	dilbert.Position = "123"

	position := &dilbert.Position //todo &dilbert.Position相当于&(dilbert.Position) 输出的是内存地址
	fmt.Println(position)

	position1 := (&dilbert).Position
	fmt.Println(position1)	//123
}
