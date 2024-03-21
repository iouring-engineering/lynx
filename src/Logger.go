package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	var value string = fmt.Sprintf("%s %s %s - ", time.Now().Format(f.timeFormat), "INFO:", f.pid)
	f.logger.Print(string(fmt.Appendln([]byte(value), v...)))
}

func (f *LoggerObj) Warn(v ...any) {
	var value string = fmt.Sprintf("%s %s %s - ", time.Now().Format(f.timeFormat), "WARN:", f.pid)
	f.logger.Print(string(fmt.Appendln([]byte(value), v...)))
}

func (f *LoggerObj) Error(v ...any) {
	var value string = fmt.Sprintf("%s %s %s - ", time.Now().Format(f.timeFormat), "ERROR:", f.pid)
	f.logger.Print(string(fmt.Appendln([]byte(value), v...)))
}

func (f *LoggerObj) Debug(v ...any) {
	if !strings.EqualFold(string(f.level), string(DEBUG)) {
		return
	}
	var value string = fmt.Sprintf("%s %s %s - ", time.Now().Format(f.timeFormat), "DEBUG:", f.pid)
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
		timeFormat: "2006-01-02 15:04:05,999",
		pid:        fmt.Sprintf("[%d]", os.Getpid()),
		logger:     logObj,
		level:      DEBUG,
		filename:   input.filename,
		path:       input.path,
	}
	log.SetOutput(file)
}
