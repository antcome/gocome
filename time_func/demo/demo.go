package demo

import (
	"github.com/antcome/gocome/time_func"
	"time"
)

func calTime() {
	{
		startTime := time.Now()
		defer time_func.PrintStack(startTime)
	}
}
