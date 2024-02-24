package util

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func Infof(format string, args ...any) {
	printLevel(2, "Info", format, args...)
}

func Debugf(format string, args ...any) {
	printLevel(2, "Debug", format, args...)
}

func currentTime() string {
	return time.Now().Format("15:04:05.000")
}

func getCallerName(skip int) (string, string, error) {
	pc, filePath, line, ok := runtime.Caller(skip)
	if !ok {
		return "", "", fmt.Errorf("getCallerName(): could not retrieve pc")
	}
	function := runtime.FuncForPC(pc)
	if function == nil {
		return "", "", fmt.Errorf("getCallerName(): function object is nil")
	}
	tmp := strings.Split(filePath, "/")
	fileName := tmp[len(tmp)-1]
	tmp = strings.Split(function.Name(), ".")
	funcName := tmp[len(tmp)-1]

	return fmt.Sprintf(`%s:%d`, fileName, line), funcName, nil
}

func printLevel(skip int, level, format string, args ...any) {
	log := fmt.Sprintf(format, args...)

	fileName, funcName, err := getCallerName(skip + 1)
	if err != nil {
		fileName = ""
		funcName = ""
	}
	fmt.Printf("%s %-7s %-20s %-20s %s\n", currentTime(), "["+level+"]", fileName, funcName, log)
}
