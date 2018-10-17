package main

import (
	"html/template"
	"os"
)

type Person70 struct {
	UserName string
}

func main() {
	t := template.New("example")
	t, _ = t.Parse("hellp {{.UserName}}")
	p := Person70{UserName: "bynf"}
	t.Execute(os.Stdout, p)
}
