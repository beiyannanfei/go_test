package t

import "log"

var index = 0

func init() {
	log.Println("==================== init ====================")
}

func Show() {
	log.Println("---- index:", index)
}

func Add() {
	index++
	if index >= 5 {
		panic("********* index bigger than five *********")
	}
}
