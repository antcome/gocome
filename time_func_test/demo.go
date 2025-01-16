package main

import (
	"fmt"
	"github.com/antcome/gocome/time_func"
	"time"
)

func main() {
	time_func.Filter = "demo.go"
	{
		startTime := time.Now()
		defer time_func.PrintStack(startTime)
	}

	fmt.Println("hello world")
}
