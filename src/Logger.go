package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
)

type Writer struct {
	io.Writer
	f              *os.File
	prevRotateTime time.Time
	filename       string
	path           string
	mu             *sync.Mutex
}

type LoggerObj struct {
	level      LogLevel
	timeFormat string
	pid        string
	logger     *log.Logger
	filename   string
	path       string
}

var Logger LoggerObj

type LogInput struct {
	filename string
	path     string
	level    LogLevel
}

func (w *Writer) getDate() string {
	var previousDateTime = w.prevRotateTime
	formattedPreviousDate := previousDateTime.Format("2006-01-02")
	return formattedPreviousDate
}

func (w *Writer) isSatisfied(prevTime time.Time) bool {
	curTime := time.Now()
	return !(prevTime.Year() == curTime.Year() && prevTime.Month() == curTime.Month() &&
		prevTime.Day() == curTime.Day())
}

func (w *Writer) rotate() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.f.Sync()
	w.f.Close()

	filePath := filepath.Join(w.path, w.filename)
	newFileName := fmt.Sprintf("%s.%s", filePath, w.getDate())

	err := os.Rename(filePath, newFileName)
	if err != nil {
		Logger.Error(err)
		return
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Logger.Error(err)
		return
	}

	info, err := file.Stat()
	if err != nil {
		Logger.Error(err)
		return
	}
	w.Writer = file
	w.f = file
	w.prevRotateTime = info.ModTime()
}

func (w *Writer) checkAndRotate() {
	if w.isSatisfied(w.prevRotateTime) {
		if w.f == nil {
			return
		}
		w.rotate()
	}
}

func (w *Writer) Write(b []byte) (n int, err error) {
	w.checkAndRotate()
	return w.Writer.Write(b)
}

func (f *LoggerObj) Info(v ...any) {
	f.log(string(INFO), v...)
}

func (f *LoggerObj) Warn(v ...any) {
	f.log(string(WARN), v...)
}

func (f *LoggerObj) Error(v ...any) {
	f.log(string(ERROR), v...)
}

func (f *LoggerObj) Debug(v ...any) {
	if !strings.EqualFold(string(f.level), string(DEBUG)) {
		return
	}
	f.log(string(DEBUG), v...)
}

func (f *LoggerObj) log(logType string, v ...any) {
	const nanoLength = 9
	logMsgType := fmt.Sprintf("%s:", logType)
	var t = time.Now()
	nanoStr := strconv.Itoa(t.Nanosecond())
	nanoStr = strings.Repeat("0", nanoLength-len(nanoStr)) + nanoStr
	tframed := fmt.Sprintf("%s.%s", f.timeFormat, nanoStr[:6])
	value := fmt.Sprintf("%s %s %s - ", tframed, logMsgType, f.pid)
	f.logger.Print(string(fmt.Appendln([]byte(value), v...)))
}

func updateLastModifiedTime(file *os.File, prevRotateTime *time.Time) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	*prevRotateTime = info.ModTime()
	return nil
}

func InitLogging() {
	var input LogInput = LogInput{
		filename: libConfig.Log.File,
		path:     libConfig.Log.Dir,
		level:    libConfig.Log.Level,
	}
	os.MkdirAll(input.path, os.ModePerm)
	var filePath = filepath.Join(input.path, input.filename)

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	var prevRotateTime = time.Now()
	updateLastModifiedTime(file, &prevRotateTime)

	logObj := log.New(&Writer{f: file, Writer: file, prevRotateTime: prevRotateTime,
		filename: input.filename, mu: &sync.Mutex{}, path: input.path}, "", 0)

	Logger = LoggerObj{
		timeFormat: "2006-01-02 15:04:05",
		pid:        fmt.Sprintf("[%d]", os.Getpid()),
		logger:     logObj,
		level:      DEBUG,
		filename:   input.filename,
		path:       input.path,
	}
	log.SetOutput(file)
}
