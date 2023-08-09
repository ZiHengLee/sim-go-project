package logger

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// logger define
var (
	Global Logger
)

func init() {
	Global = NewDefaultLogger(DEBUG)
}

func Init(opt *Option) (err error) {
	var lvl level
	if len(opt.Level) > 0 {
		err = lvl.Parse(opt.Level)
		if err != nil {
			fmt.Printf("invalid global level:%v", opt.Level)
			return
		}
	}
	Global = make(Logger)
	if opt.Stdout != nil {
		slvl := lvl
		if len(opt.Stdout.Level) > 0 {
			err = slvl.Parse(opt.Stdout.Level)
			if err != nil {
				fmt.Printf("invalid stdout level:%v", opt.Stdout.Level)
				return
			}
		}
		w := NewStdoutLogWriter()
		Global.AddFilter("stdout", slvl, w)
	}

	if opt.Files != nil {
		for k, v := range opt.Files {
			flvl := lvl
			if len(v.Level) > 0 {
				err = flvl.Parse(v.Level)
				if err != nil {
					fmt.Printf("invalid file:%v level:%v", k, v.Level)
					return
				}
			}
			w, err := newFileLogWriter(v)
			if err != nil {
				fmt.Printf("new file:%#v log err:%v", v, err)
				return err
			}
			Global.AddFilter("file-"+k, flvl, w)
		}
	}

	if opt.Dingdings != nil {
		for k, v := range opt.Dingdings {
			dlvl := lvl
			if len(v.Level) > 0 {
				err = dlvl.Parse(v.Level)
				if err != nil {
					fmt.Printf("invalid dingding:%v level:%v", k, v.Level)
					return
				}
			}
			w, err := NewDingDingWriter(v)
			if err != nil {
				fmt.Printf("new dingding:%#v log err:%v", v, err)
				return err
			}
			Global.AddFilter("dingding-"+k, dlvl, w)
		}
	}

	return
}

func GetCurrPath() (string, string) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	p1 := path[:index]
	p2 := path[index:]
	return p1, p2
}

func newFileLogWriter(opt *FileOption) (w *FileLogWriter, err error) {
	file := strings.Trim(opt.FileName, " \r\n")
	if file == "" {
		p1, p2 := GetCurrPath()
		index := strings.LastIndex(p1, string(os.PathSeparator))
		path := p1[:index] + "/log"
		os.MkdirAll(path, 0777)
		file = path + p2
	} else {
		index := strings.LastIndex(file, string(os.PathSeparator))
		if index >= 0 {
			path := file[:index] + "/"
			if _, err := os.Stat(path); err != nil {
				if os.IsNotExist(err) {
					os.MkdirAll(path, 0777)
				}
			}
		}
	}

	format := strings.Trim(opt.Format, " \r\n")
	if len(format) == 0 {
		format = "[%D %T] [%L] (%S) %M"
	}

	w = NewFileLogWriter(file, !opt.NoRotate)
	w.SetFormat(format)
	w.SetRotateLines(opt.MaxLines)
	w.SetRotateSize(opt.MaxSize)
	w.SetRotateDaily(!opt.NoDaily)
	return
}

// AddFilter : Wrapper for (*Logger).AddFilter
func AddFilter(name string, lvl level, writer LogWriter) {
	Global.AddFilter(name, lvl, writer)
}

// Close : Wrapper for (*Logger).Close (closes and removes all logwriters)
func Close() {
	Global.Close()
}

// Crash _
func Crash(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	panic(args)
}

// Crashf : Logs the given message and crashes the program
func Crashf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	panic(fmt.Sprintf(format, args...))
}

// Exit : Compatibility with `log`
func Exit(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Exitf : Compatibility with `log`
func Exitf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
	Global.Close() // so that hopefully the messages get logged
	os.Exit(0)
}

// Stderr : Compatibility with `log`
func Stderr(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(ERROR, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Stderrf : Compatibility with `log`
func Stderrf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
}

// Stdout : Compatibility with `log`
func Stdout(args ...interface{}) {
	if len(args) > 0 {
		Global.intLogf(INFO, strings.Repeat(" %v", len(args))[1:], args...)
	}
}

// Stdoutf : Compatibility with `log`
func Stdoutf(format string, args ...interface{}) {
	Global.intLogf(INFO, format, args...)
}

// Log : Send a log message manually
// Wrapper for (*Logger).Log
func Log(lvl level, source, message string) {
	Global.Log(lvl, source, message)
}

// Logf : Send a formatted log message easily
// Wrapper for (*Logger).Logf
func Logf(lvl level, format string, args ...interface{}) {
	Global.intLogf(lvl, format, args...)
}

// Logc : Send a closure log message
// Wrapper for (*Logger).Logc
func Logc(lvl level, closure func() string) {
	Global.intLogc(lvl, closure)
}

// Debug : Utility for debug log messages
// When given a string as the first argument, this behaves like Logf but with the DEBUG log level (e.g. the first argument is interpreted as a format for the latter arguments)
// When given a closure of type func()string, this logs the string returned by the closure iff it will be logged.  The closure runs at most one time.
// When given anything else, the log message will be each of the arguments formatted with %v and separated by spaces (ala Sprint).
// Wrapper for (*Logger).Debug
func Debug(arg0 interface{}, args ...interface{}) {
	var (
		lvl = DEBUG
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Trace : Utility for trace log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Trace
func Trace(arg0 interface{}, args ...interface{}) {
	var (
		lvl = TRACE
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Info : Utility for info log messages (see Debug() for parameter explanation)
// Wrapper for (*Logger).Info
func Info(arg0 interface{}, args ...interface{}) {
	var (
		lvl = INFO
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		Global.intLogc(lvl, first)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(arg0)+strings.Repeat(" %v", len(args)), args...)
	}
}

// Warn : Utility for warn log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Warn
func Warn(arg0 interface{}, args ...interface{}) error {
	var (
		lvl = WARNING
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return fmt.Errorf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
}

// Error : Utility for error log messages (returns an error for easy function returns) (see Debug() for parameter explanation)
// These functions will execute a closure exactly once, to build the error message for the return
// Wrapper for (*Logger).Error
func Error(arg0 interface{}, args ...interface{}) error {
	var (
		lvl = ERROR
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return fmt.Errorf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
}

func Alert(arg0 interface{}, args ...interface{}) error {
	var (
		lvl = ALERT
	)
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		Global.intLogf(lvl, first, args...)
		return fmt.Errorf(first, args...)
	case func() string:
		// Log the closure (no other arguments used)
		str := first()
		Global.intLogf(lvl, "%s", str)
		return errors.New(str)
	default:
		// Build a format string so that it will be similar to Sprint
		Global.intLogf(lvl, fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
		return errors.New(fmt.Sprint(first) + fmt.Sprintf(strings.Repeat(" %v", len(args)), args...))
	}
}

func Fatal(args ...interface{}) {
	Crash(args)
}

func Fatalf(format string, args ...interface{}) {
	Crashf(format, args...)
}

func Fatalln(args ...interface{}) {
	Crash(args)
}

func Infoln(args ...interface{}) {
	var (
		lvl = INFO
	)
	Global.intLogf(lvl, strings.Repeat(" %v", len(args)), args...)
}

func Traceln(args ...interface{}) {
	var (
		lvl = TRACE
	)
	Global.intLogf(lvl, strings.Repeat(" %v", len(args)), args...)
}

func Recover() {
	err := recover()
	if err != nil {
		buf := make([]byte, 10240)
		runtime.Stack(buf, false)
		Alert("panic: %v,\n%s", err, string(buf))
	}
}

func PrintStack(err interface{}) {
	buf := make([]byte, 10240)
	runtime.Stack(buf, false)
	Error("stack: %v,\n%s", err, string(buf))
}
