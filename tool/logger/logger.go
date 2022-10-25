package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	LevelAll   Level = 0
	LevelDebug Level = 1
	LevelInfo  Level = 2
	LevelError Level = 3
	LevelFatal Level = 4
	LevelNone  Level = 5
)

var prefixes = map[Level]string{
	LevelDebug: "DEBUG ",
	LevelInfo:  "INFO  ",
	LevelError: "ERROR ",
	LevelFatal: "FATAL ",
}

type (
	Level   int
	Options struct {
		MinLevel Level
		Output   io.Writer
	}
	Logger interface {
		Print(level Level, msg ...string)
		Debug(msg ...string)
		Info(msg ...string)
		Error(msg ...string)
		Fatal(msg ...string)
		DebugF(msg string, values ...any)
		InfoF(msg string, values ...any)
		ErrorF(msg string, values ...any)
		FatalF(msg string, values ...any)
		ErrorE(err error)
		FatalE(err error)
		WithStackTrace() Logger
		AddPrefix(str string) Logger
	}
	logger struct {
		output         io.Writer
		minLevel       Level
		msgPrefixes    []string
		withStackTrace bool
	}
)

func New() Logger {
	return WithOptions(Options{})
}

func WithOptions(opts Options) Logger {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	return &logger{
		minLevel: opts.MinLevel,
		output:   opts.Output,
	}
}

func (l logger) Print(level Level, msg ...string) {
	l.print(level, msg...)
}

func (l logger) Debug(msg ...string) {
	l.print(LevelDebug, msg...)
}

func (l logger) Info(msg ...string) {
	l.print(LevelInfo, msg...)
}

func (l logger) Error(msg ...string) {
	l.print(LevelError, msg...)
}

func (l logger) Fatal(msg ...string) {
	l.print(LevelFatal, msg...)
}

func (l logger) DebugF(msg string, values ...any) {
	l.print(LevelDebug, fmt.Sprintf(msg, values...))
}

func (l logger) InfoF(msg string, values ...any) {
	l.print(LevelInfo, fmt.Sprintf(msg, values...))
}

func (l logger) ErrorF(msg string, values ...any) {
	l.print(LevelError, fmt.Sprintf(msg, values...))
}

func (l logger) FatalF(msg string, values ...any) {
	l.print(LevelFatal, fmt.Sprintf(msg, values...))
}

func (l logger) ErrorE(err error) {
	l.print(LevelError, err.Error())
}

func (l logger) FatalE(err error) {
	l.print(LevelFatal, err.Error())
}

func (l logger) WithStackTrace() Logger {
	l.withStackTrace = true
	return l
}

func (l logger) AddPrefix(prefix string) Logger {
	l.msgPrefixes = append(l.msgPrefixes, prefix)
	return l
}

func (l logger) print(level Level, msg ...string) {
	if level < l.minLevel {
		return
	}

	var builder strings.Builder

	builder.WriteString(prefixes[level])
	builder.WriteRune(' ')

	builder.WriteString(time.Now().Format("2006-01-02 15:04:05 -0700"))
	builder.WriteRune(' ')

	if len(l.msgPrefixes) != 0 {
		writeStringArray(&builder, l.msgPrefixes, '/')
		builder.WriteRune(' ')
	}

	writeStringArray(&builder, msg, ' ')

	if l.withStackTrace {
		writeStackTrace(&builder, 4)
	}

	builder.WriteRune('\n')

	_, err := l.output.Write([]byte(builder.String()))
	if err != nil {
		panic(err)
	}
	if level == LevelFatal {
		os.Exit(1)
	}
}

func writeStackTrace(builder *strings.Builder, skipCalls int) {
	var l *location
	for ; ; skipCalls++ {
		l = getLocation(skipCalls)
		if l == nil {
			break
		}

		builder.WriteString("\n\tat ")
		builder.WriteString(l.String())

		if l.function == "main.main" {
			break
		}
	}
}

func writeStringArray(builder *strings.Builder, msg []string, sep rune) {
	builder.WriteString(msg[0])
	for i := 1; i < len(msg); i++ {
		builder.WriteRune('/')
		builder.WriteString(msg[i])
	}
}
