package time_func

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

func AstRewrite(dir ...string) (err error) {
	files, funcErr := getAllDirs(dir...)
	if funcErr != nil {
		err = funcErr
		return
	}

	fileFilter := make(map[string]bool)
	for _, file := range files {
		err = astRewriteFile(file, func(path string, n ast.Node) bool {
			if fileFilter[path] {
				return false
			}
			fileFilter[path] = true
			importDemoFile(n)
			return true
		})
		if err != nil {
			return err
		}
	}
	return
}

var (
	rootDir, _  = os.Getwd()
	processName = ""
	Filter      = "/data_config/router"
)

func init() {
	tmp := strings.Split(rootDir, "/")
	processName = tmp[len(tmp)-1]
	Filter = "/" + processName + "/"
}

func PrintStack(startTime time.Time) {
	//time.Sleep(1 * time.Microsecond)
	mil := time.Since(startTime).Milliseconds()
	if mil <= 3 {
		return
	}

	stack := string(debug.Stack())
	if !strings.Contains(stack, Filter) {
		return
	}

	method := ""
	{
		pc, file, _, _ := runtime.Caller(1)
		file = processName + strings.ReplaceAll(file, rootDir, "")
		method = runtime.FuncForPC(pc).Name()
		method = filepath.Ext(method) + "()"
		method = file + method
	}

	sourceFileLine := ""
	{
		_, file, line, _ := runtime.Caller(2)
		//method2 := runtime.FuncForPC(pc2).Name()
		if strings.Contains(file, processName) {
			sourceFileLine = fmt.Sprintf(":%d", line)
		}
	}

	count := strings.Count(stack, rootDir)
	spaceCnt := make([]string, count)
	spaceCnt[count-1] = fmt.Sprint(count)
	spaceStr := strings.Join(spaceCnt, "  ")

	//tip := fmt.Sprintf("@@%d %s_%s 耗时执行:%d", startTime.UnixNano(), spaceStr, method, mil)
	tip := fmt.Sprintf("@@%d 耗时执行:%-4d %s_%s%s", startTime.UnixNano(), mil, spaceStr, method, sourceFileLine)
	//fmt.Println(tip)
	defaultSafeArray.Add(int(startTime.UnixNano()), count, tip)
	if count <= 2 {
		fmt.Println()
		fmt.Println(rootDir)
		fmt.Println(strings.Join(defaultSafeArray.ReserveAndClear(), "\n"))
		fmt.Println()
	}
}
