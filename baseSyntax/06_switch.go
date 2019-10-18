/**
 * Created by wyq on 18/2/4.
switch 语句用于基于不同条件执行不同动作，每一个 case 分支都是唯一的，从上直下逐一测试，直到匹配为止。。
switch 语句执行的过程从上至下，直到找到匹配项，匹配项后面也不需要再加break
 */

package main

import "fmt"

func switch1() {
	var marks = 60
	switch marks {
	case 90:
		fmt.Println("marks is 90")
	case 80:
		fmt.Println("marks is 80")
	case 70, 60, 50:
		fmt.Println("marks is 70 || 60 || 50")
	default:
		fmt.Println("marks is less than 50")
	}
}

func switch2() {
	var grade = "B"
	switch {
	case grade == "A":
		fmt.Println("优秀")
	case grade == "B":
		fmt.Println("良好")
	default:
		fmt.Println("及格")

	}
}

func switch_type() {
	var x interface{}
	x = 12.34

	switch i := x.(type) { //技巧 := 可以不使用var声明，所有i会根据赋值类型自动判断，所以这个很有技巧性
	case nil:
		fmt.Printf(" x 的类型 :%T", i)
	case int:
		fmt.Printf("x 是 int 型")
	case float64:
		fmt.Printf("x 是 float64 型")
	case func(int) float64:
		fmt.Printf("x 是 func(int) 型")
	case bool, string:
		fmt.Printf("x 是 bool 或 string 型")
	default:
		fmt.Printf("未知型")
	}
}

func main() {
	switch1()
	switch2()
	switch_type()
}
