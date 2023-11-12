package main

import (
	"fmt"
	"sync"
)

func main() {
	var newMap sync.Map

	newMap.Store(1, "eva")
	fmt.Println(newMap.Load(1))
}
