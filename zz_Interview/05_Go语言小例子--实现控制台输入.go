package main

import (
	"bufio"
	"os"
	"strconv"
	"log"
)

func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)

	for input.Scan() {
		i, err := strconv.Atoi(input.Text())
		if err != nil {
			log.Println("input err: ", err)
			break
		}

		if i <= 0 {
			log.Println("input negative and break.")
			break
		}

		counts[input.Text()]++

	}

	for line, n := range counts {
		log.Println(line, n)
	}
}
