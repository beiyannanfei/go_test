package main

import (
	"html/template"
	"os"
	"strings"
	"fmt"
)

//如果我们想要的email函数的模板函数名是emailDeal，它关联的Go函数名称是EmailDealWith,n那么我们 可以通过下面的方式来注册这个函数
//t = t.Funcs(template.FuncMap{"emailDeal": EmailDealWith}) EmailDealWith这个函数的参数和返回值定义如下:
//func EmailDealWith(args ...interface{}) string

type Friend72 struct {
	Fname string
}

type Person72 struct {
	UserName string
	Emails   []string
	Friends  []*Friend72
}

func EmailDealWith(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		return "aaa"
	}

	substrs := strings.Split(s, "@")
	if len(substrs) != 2 {
		return s
	}

	return fmt.Sprintf("%s at %s", substrs[0], substrs[1])
}

func main() {
	f1 := Friend72{Fname: "zhangsan"}
	f2 := Friend72{Fname: "lisi"}
	t := template.New("feg")
	t = t.Funcs(template.FuncMap{"emailDeal": EmailDealWith})
	t, _ = t.Parse(
		`hello {{.UserName}}
{{range .Emails}}
	an emails {{.|emailDeal}}{{end}}
{{with .Friends}}
{{range .}}
	my friend name is {{.Fname}}{{end}}
{{end}}
`,
	)
	p := Person72{
		UserName: "beiyannanfei",
		Emails:   []string{"aa@bb.com", "cc@dd.com"},
		Friends:  []*Friend72{&f1, &f2},
	}
	t.Execute(os.Stdout, p)
}
