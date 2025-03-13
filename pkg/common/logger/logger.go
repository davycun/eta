package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
)

var (
	Logger     *logger
	errTag     = fmt.Sprintf("%s %s %s", red, "error", reset)
	warningTag = fmt.Sprintf("%s %s %s", yellow, "warn", reset)
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func init() {
	//file, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0x755)
	//if err != nil {
	//	panic(err.Error())
	//}
	Logger = &logger{
		Logger: log.New(os.Stderr, "", log.Lmsgprefix|log.LstdFlags|log.Lshortfile),
	}
}

func Errorf(format string, v ...any) {
	outf(errTag, format, v...)
}
func Errorln(v ...any) {
	outln(errTag, v...)
}
func Error(v ...any) {
	out(errTag, v...)
}

func Warnf(format string, v ...any) {
	outf(warningTag, format, v...)
}
func Warnln(v ...any) {
	outln(warningTag, v...)
}
func Warn(v ...any) {
	out(warningTag, v...)
}

func Infof(format string, v ...any) {
	outf("info", format, v...)
}
func Infoln(v ...any) {
	outln("info", v...)
}
func Info(v ...any) {
	out("info", v...)
}

func Debugf(format string, v ...any) {
	outf("debug", format, v...)
}
func Debugln(v ...any) {
	outln("debug", v...)
}
func Debug(v ...any) {
	out("debug", v...)
}

func outf(tag string, format string, v ...any) {
	Logger.Output(3, fmt.Sprintf("[%s] %s", tag, fmt.Sprintf(format, v...)))
}
func outln(tag string, v ...any) {
	Logger.Output(3, fmt.Sprintf("[%s] %s", tag, fmt.Sprintln(v...)))
}
func out(tag string, v ...any) {
	Logger.Output(3, fmt.Sprintf("[%s] %s", tag, fmt.Sprint(v...)))
}

func OutputPanic(v ...any) {
	bf := bytes.NewBuffer(make([]byte, 0))
	idx := 3
	//如果想打印更深的栈改大10即可
	for i := 3; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			bf.WriteString(file)
			bf.WriteByte(':')
			bf.WriteString(strconv.Itoa(line))
			bf.WriteByte('\n')
			idx = i
		} else {
			break
		}
	}
	Logger.Output(idx, fmt.Sprint(v...))
	Logger.Writer().Write(bf.Bytes())
}

func Writer() io.Writer {
	return Logger.Writer()
}
