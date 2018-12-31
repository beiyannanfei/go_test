package main

import (
	"errors"
	"fmt"
)

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	Remove(&arr, 5)
	fmt.Println(arr)
}

func Remove(arr *[]int, value int) error {
	for i, v := range *arr {
		if value == v {
			*arr = append((*arr)[:i], (*arr)[i+1:]...)
			return nil
		}
	}

	return errors.New("elem not exists")
}
