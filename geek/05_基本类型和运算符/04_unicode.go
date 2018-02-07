package main

import (
	"fmt"
	"unicode"
)

func main() {
	var chA = 'A'
	var chB = '3'
	var chC = ' '

	fmt.Printf("chA: %t, chB: %t, chC: %t\n", unicode.IsLetter(chA), unicode.IsDigit(chB), unicode.IsSpace(chC))	//chA: true, chB: true, chC: true
	fmt.Printf("chA: %t, chB: %t, chC: %t\n", unicode.IsLetter(chB), unicode.IsDigit(chC), unicode.IsSpace(chA))	//chA: false, chB: false, chC: false
}
