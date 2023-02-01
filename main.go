package main

import (
	"fmt"
	"runtime"
)

func main() {
	jobdistributor.Init()
	go jobdistributor.Instance.Run(10)

	runtime.Goexit()

	fmt.Println("Exit")
}
