package lock

import (
	"log"
	"sync"
	"time"
)

var Lock sync.RWMutex

var DepartmentMap = make(map[int]string)

func GetSubTreeByDepartmentId(departmentId int) string {
	Lock.RLock()
	defer Lock.RUnlock()

	log.Println("begin read")

	node, ok := DepartmentMap[departmentId]
	if !ok {
		log.Println("begin null")
		return ""
	}

	log.Println("end read")
	return node
}

func do() {
	Lock.Lock()
	defer Lock.Unlock()

	log.Println("begin write")

	time.Sleep(5 * time.Second)

	DepartmentMap[10] = "aaa"

	log.Println("end write")
}

func BuildTree() {
	for {
		do()
		time.Sleep(1 * 60 * time.Second)
	}
}
