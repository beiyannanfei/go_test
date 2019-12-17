package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

func main() {
	{
		type Person1 struct {
			UserName string
		}

		t1 := template.New("fieldname example")
		t1, _ = t1.Parse("hello {{.UserName}}!")
		p1 := Person1{UserName: "bynf"}
		_ = t1.Execute(os.Stdout, p1)
		fmt.Printf("\n==================================\n")
	}

	{
		type Friend struct {
			Fname string
		}

		type Person struct {
			UserName string
			Emails   []string
			Friends  []*Friend
		}

		parseStr := `hello {{.UserName}}!
{{range .Emails}}
	an email {{.}}
{{end}}
{{with .Friends}} 
{{range .}}
	my friends name is {{.Fname}} 
{{end}} 
{{end}}`

		f1 := Friend{Fname: "minux.ma"}
		f2 := Friend{Fname: "xushiwei"}
		t := template.New("fieldname example")
		t, _ = t.Parse(parseStr)
		p := Person{UserName: "bynf",
			Emails:  []string{"astaxie@beego.me", "astaxie@gmail.com"},
			Friends: []*Friend{&f1, &f2},
		}
		_ = t.Execute(os.Stdout, p)
		fmt.Printf("\n==================================\n")
	}

	{
		type Friend struct {
			Fname string
		}

		type Person struct {
			UserName string
			Emails   []string
			Friends  []*Friend
		}

		f1 := Friend{Fname: "minux.ma"}
		f2 := Friend{Fname: "xushiwei"}
		t := template.New("fieldname example")
		t = t.Funcs(template.FuncMap{"emailDeal": EmailDealWith})
		t, _ = t.Parse(`hello {{.UserName}}!
{{range .Emails}}
	an emails {{.|emailDeal}}
{{end}}
{{with .Friends}}
{{range .}}
	my friend name is {{.Fname}}
{{end}}
{{end}}`)
		p := Person{UserName: "Astaxie",
			Emails:  []string{"astaxie@begoo.me", "astaxie@gmail.com"},
			Friends: []*Friend{&f1, &f2},
		}
		_ = t.Execute(os.Stdout, p)
		fmt.Printf("\n==================================\n")
	}

}

func EmailDealWith(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}

	if !ok {
		s = fmt.Sprint(args...)
	}

	subStrs := strings.Split(s, "@")
	if len(subStrs) != 2 {
		return s
	}

	return subStrs[0] + " at " + subStrs[1]
}
