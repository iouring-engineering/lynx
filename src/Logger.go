package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

type LogType string

const (
	INFO  LogType = "INFO"
	ERROR LogType = "ERROR"
	WARN  LogType = "WARN"
	DEBUG LogType = "DEBUG"
)

type writer struct {
	io.Writer
	timeFormat string
	level      string
	pid        string
}

func (w *writer) Write(b []byte) (n int, err error) {
	var value string = fmt.Sprintf("%s %s %s - ", time.Now().Format(w.timeFormat), w.level, w.pid)
	return w.Writer.Write(append([]byte(value), b...))
}

func InitLogging() {
	err := os.MkdirAll(libConfig.Log.Dir, os.ModePerm)

	filePath := libConfig.Log.File

	var iLogLevel LogType = DEBUG
	if libConfig.Log.Level != "" {
		iLogLevel = libConfig.Log.Level
		if iLogLevel != INFO && iLogLevel != DEBUG {
			log.Fatal("Invalid log level")
		}
	}

	if err == nil {
		filePath = filepath.Join(libConfig.Log.Dir, libConfig.Log.File)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(&writer{file, "2006-01-02 15:04:05,999", "INFO:", fmt.Sprintf("[%d]", os.Getpid())}, "", 0)
	WarningLogger = log.New(&writer{file, "2006-01-02 15:04:05,999", "WARN:", fmt.Sprintf("[%d]", os.Getpid())}, "", 0)
	ErrorLogger = log.New(&writer{file, "2006-01-02 15:04:05,999", "ERROR:", fmt.Sprintf("[%d]", os.Getpid())}, "", 0)
	DebugLogger = log.New(&writer{file, "2006-01-02 15:04:05,999", "DEBUG:", fmt.Sprintf("[%d]", os.Getpid())}, "", 0)

	if iLogLevel == INFO {
		DebugLogger.SetOutput(io.Discard)
	}
	log.SetOutput(file)
}
