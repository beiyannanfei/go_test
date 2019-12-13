package main

import (
	"fmt"
	"sort"
)

type subLineData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type subMapData struct {
	LineType string         `json:"lineType"`
	LineData []*subLineData `json:"lineData"`
}

type subMap struct {
	MapType ProjectType   `json:"mapType"`
	MapData []*subMapData `json:"mapData"`
}

type subSummaryRsp struct {
	Dashboard map[string]int `json:"dashboard"`
	Maps      []*subMap      `json:"maps"`
}

type ProjectType int

const (
	_ ProjectType = iota
	Reach
	Bind
	Raise
	Huge
)

func main() {
	s := []string{
		"北京",
		"河北",
		"天津",
	}

	sort.Sort(sort.Reverse(sort.StringSlice(s)))

	fmt.Println(s)
}
