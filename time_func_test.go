package gocome

import (
	"fmt"
	"github.com/antcome/gocome/time_func"
	"testing"
)

func TestName(t *testing.T) {
	err := time_func.AstRewrite("./time_func_test")
	fmt.Println(err)
}
