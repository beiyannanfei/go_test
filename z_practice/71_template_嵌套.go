package main

import (
	"html/template"
	"os"
)

type Friend71 struct {
	Fname string
}
type Person71 struct {
	UserName string
	Emails   []string
	Friends  []*Friend71
}

func main() {
	f1 := Friend71{Fname: "zhangsan"}
	f2 := Friend71{Fname: "lisi"}
	t := template.New("eg")
	t, _ = t.Parse(
		`hello {{.UserName}}! 
{{range .Emails}}
	an email {{.}}{{end}}
{{with .Friends}}
{{range .}}
	my friend name is {{.Fname}}{{end}}
{{end}}
		`)
	p := Person71{
		UserName: "bynf",
		Emails:   []string{"a@b.com", "c@d.com"},
		Friends:  []*Friend71{&f1, &f2},
	}
	t.Execute(os.Stdout, p)
}
