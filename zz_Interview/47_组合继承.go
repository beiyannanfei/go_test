package main

import "fmt"

type People47 struct{}

func (p People47) ShowA() {
	fmt.Println("ShowA")
	p.ShowB()
}

func (p People47) ShowB() {
	fmt.Println("ShowB")
}

type Teacher47 struct {
	People47
}

func (t *Teacher47) ShowB() {
	fmt.Println("teacher showB")
}

func main() {
	t := Teacher47{}
	t.ShowA()          //ShowA 	ShowB
	t.ShowB()          //teacher showB
	t.People47.ShowB() //ShowB
}
