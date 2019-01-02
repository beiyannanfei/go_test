package main

import (
	"log"
	"encoding/json"
)

//https://studygolang.com/articles/17370
//Go可以声明自定义的数据类型，组合一个或多个类型，可以包含内置类型和用户自定义的类型，可以像内置类型一样使用struct类型

type rectangle01 struct {
	length  int
	breadth int
	color   string

	geometry struct {
		area      int
		perimeter int
	}
}

type rectangle02 struct {
	length  int
	breadth int
	color   string
}

type rectangle03 struct {
	length  int
	breadth int
	color   string
}

type rectangle04 struct {
	length  int
	breadth int
	color   string
}

func main() {
	//点运算符	可以使用点运算符访问结构体中的数据值
	var rec rectangle01
	rec.breadth = 19
	rec.length = 23
	rec.color = "Green"

	rec.geometry.area = rec.length * rec.breadth
	rec.geometry.perimeter = 2 * (rec.length + rec.breadth)

	log.Printf("rec: %v\n", rec)                       //2019/01/02 18:10:50 rec: {23 19 Green {437 84}}
	log.Println("area: ", rec.geometry.area)           //2019/01/02 18:10:50 area:  437
	log.Println("perimeter: ", rec.geometry.perimeter) //2019/01/02 18:10:50 perimeter:  84
	log.Println("---------------------------------------")

	//使用 var关键词和 :=运算符	如果初始化时，指定了特定的名称，那么有些字段是可以省略的
	var rect1 = rectangle02{10, 20, "red"}
	log.Println(rect1) //2019/01/02 18:10:50 {10 20 red}
	var rect2 = rectangle02{length: 10, color: "red2"}
	log.Println(rect2) //2019/01/02 18:10:50 {10 0 red2}
	rect3 := rectangle02{10, 20, "green"}
	log.Println(rect3) //2019/01/02 18:10:50 {10 20 green}
	rect4 := rectangle02{length: 100, breadth: 200, color: "green"}
	log.Println(rect4) //2019/01/02 18:10:50 {100 200 green}
	rect5 := rectangle02{breadth: 22, color: "green"}
	log.Println(rect5) //2019/01/02 18:10:50 {0 22 green}
	log.Println("---------------------------------------")

	//使用 new 关键字
	re1 := new(rectangle03)
	re1.length = 10
	re1.breadth = 20
	re1.color = "green"
	log.Println(re1) //2019/01/02 18:10:50 &{10 20 green}

	re2 := new(rectangle03)
	re2.breadth = 200
	re2.color = "red"
	log.Println(re2) //2019/01/02 18:10:50 &{0 200 red}
	log.Println("---------------------------------------")

	//使用 & 运算符
	var r1 = &rectangle04{10, 20, "red"}
	log.Println(r1) //2019/01/02 18:10:50 &{10 20 red}

	var r2 = &rectangle04{}
	r2.length = 100
	r2.color = "red"
	log.Println(r2) //2019/01/02 18:10:50 &{100 0 red}

	var r3 = &rectangle04{}
	(*r3).breadth = 20
	(*r3).color = "blue"
	log.Println(r3) //2019/01/02 18:10:50 &{0 20 blue}
	log.Println("---------------------------------------")

	//struct中的tag标签
	type Employee struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		City      string `json:"city"`
	}

	json_str := `
    {
        "firstname":"Kevin",
        "lastname":"Woo",
        "city":"Beijing"
    }`
	emp1 := new(Employee)
	err := json.Unmarshal([]byte(json_str), emp1)
	if err != nil {
		log.Fatalln("Unmarshal err: ", err)
	}
	log.Println(emp1) //2019/01/02 18:10:50 &{Kevin Woo Beijing}

	emp2 := new(Employee)
	emp2.FirstName = "John"
	emp2.LastName = "Lee"
	emp2.City = "shanghai"
	jsonStr, _ := json.Marshal(emp2)
	log.Println(string(jsonStr)) //2019/01/02 18:10:50 {"firstname":"John","lastname":"Lee","city":"shanghai"}
	log.Println("---------------------------------------")

	//内嵌的struct 类型
	type Salary struct {
		Basic float64 `json:"basic"`
		HRA   float64 `json:"hra"`
		TA    float64 `json:"ta"`
	}
	type Employee1 struct {
		FirstName     string   `json:"first_name"`
		LastName      string   `json:"last_name"`
		Email         string   `json:"email"`
		Age           int      `json:"age"`
		MonthlySalary []Salary `json:"monthly_salary"`
	}

	e := Employee1{
		FirstName: "kevin",
		LastName:  "Woo",
		Email:     "test@mail.com",
		Age:       12,
		MonthlySalary: []Salary{
			Salary{
				Basic: 15000.00,
				HRA:   5000.0,
				TA:    2000.0,
			},
			Salary{
				Basic: 16000.0,
				HRA:   6000.0,
				TA:    2100.0,
			},
		},
	}
	log.Println(e.FirstName, e.LastName) //2019/01/02 19:37:34 kevin Woo
	log.Println(e.Age)                   //2019/01/02 19:38:24 12
	log.Println(e.Email)                 //2019/01/02 19:38:24 test@mail.com
	log.Println(e.MonthlySalary[0])      //2019/01/02 19:38:24 {15000 5000 2000}
	log.Println(e.MonthlySalary[1])      //2019/01/02 19:38:24 {16000 6000 2100}
	jStr, _ := json.Marshal(e)
	log.Println(string(jStr)) //2019/01/02 19:40:15 {"first_name":"kevin","last_name":"Woo","email":"test@mail.com","age":12,"monthly_salary":[{"basic":15000,"hra":5000,"ta":2000},{"basic":16000,"hra":6000,"ta":2100}]}
	log.Println("---------------------------------------")

	e1 := EmloyeeA{
		FirstName: "KevinA",
		LastName:  "WooA",
		Email:     "test@mail.comA",
		Age:       26,
		MonthlySalary: []SalaryA{
			SalaryA{
				Basic: 15001.00,
				HRA:   5001.0,
				TA:    2001.0,
			},
			SalaryA{
				Basic: 16001.0,
				HRA:   6001.0,
				TA:    2101.0,
			},
		},
	}

	log.Println(e1.EmpInfo()) //2019/01/02 19:47:33 {"first_name":"KevinA","last_name":"WooA","email":"test@mail.comA","age":26,"monthly_salary":[{"basic":15001,"hra":5001,"ta":2001},{"basic":16001,"hra":6001,"ta":2101}]}
}

type SalaryA struct {
	Basic float64 `json:"basic"`
	HRA   float64 `json:"hra"`
	TA    float64 `json:"ta"`
}

type EmloyeeA struct {
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Age           int       `json:"age"`
	MonthlySalary []SalaryA `json:"monthly_salary"`
}

func (e EmloyeeA) EmpInfo() string {
	jStr, _ := json.Marshal(e)
	return string(jStr)
}
