package logger

import (
	"fmt"
	"runtime"
	"strings"
)

type location struct {
	function string
	file     string
	line     int
}

func (l location) String() string {
	return fmt.Sprintf("%s:%d(%s)", l.file, l.line, l.function)
}

func getLocation(skipCalls int) *location {
	pc, file, line, ok := runtime.Caller(skipCalls)
	if !ok {
		return nil
	}
	fn := getFuncName(pc)
	return &location{
		function: fn,
		file:     file,
		line:     line,
	}
}

func getFuncName(pc uintptr) string {
	fn := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fn, "/")
	return parts[len(parts)-1]
}
