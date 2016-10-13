package main

import (
	"fmt"
)

func main() {
	FZMPerformacCounters := GetCounters()
	slice := FZMPerformacCounters[0:4]
	fmt.Println(slice)
}
