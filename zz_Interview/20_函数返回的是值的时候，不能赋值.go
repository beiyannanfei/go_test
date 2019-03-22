package main

import (
	"time"
)

type Employee20 struct {
	ID        int
	Name      string
	Address   string
	DoB       time.Time
	Position  string
	Salary    int
	ManagerID int
}

func EmployeeByID20(id int) Employee20 {
	return Employee20{ID: id}
}

func main() {
	EmployeeByID20(1).Salary = 20 //cannot assign to EmployeeByID20(1).Salary
	//todo 在本例子中，函数EmployeeById(id int)返回的是值类型的，它的取值EmployeeByID(1).Salary也是一个值类型；值类型是什么概念？值类型就是和赋值语句var a = 1或var a = hello world等号=右边的1、Hello world是一个概念，他是不能够被赋值的，只有变量能够被赋值

	/* 正确写法
		var a = EmployeeByID20(1)
		a.Salary = 20
		fmt.Println(a)
	*/
}
