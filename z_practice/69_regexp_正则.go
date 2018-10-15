package main

// cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 69_regexp_正则.go aaa

import (
	"regexp"
	"fmt"
	"os"
)

func IsIp(ip string) bool {
	regexpReguler := "^[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}$"
	if m, _ := regexp.MatchString(regexpReguler, ip); !m {
		return false
	}

	return true
}

func main() {
	ip := "127.0.0.1"
	r := IsIp(ip)
	fmt.Printf("%v is ip? %v\n", ip, r)
	fmt.Println(os.Args)

	if len(os.Args) == 1 {
		fmt.Println("Usage: regexp [string]")
		os.Exit(1)
	}

	if m, _ := regexp.MatchString("^[0-9]+$", os.Args[1]); m {
		fmt.Println("is number")
	} else {
		fmt.Println("not number")
	}
}
