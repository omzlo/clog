package clog

import (
	"fmt"
	"os"
)

type LogWriter interface {
	LogWrite(level LogLevel, text string)
}

type LogLevel uint

const (
	DEBUGXX LogLevel = iota
	DEBUGX
	DEBUG
	INFO
	WARNING
	ERROR
	NONE
)

type logEntry struct {
	level LogLevel
	text  string
}

type bucket struct {
	writer LogWriter
	next   *bucket
}

type LogManager struct {
	logQueue    chan *logEntry
	logLevel    LogLevel
	firstBucket *bucket
	lastBucket  *bucket
}

var DefaultLogManager *LogManager = NewLogManager()

func NewLogManager() *LogManager {
	return &LogManager{
		logLevel: INFO,
	}
}

func (lm *LogManager) AddWriter(writer LogWriter) {
	b := &bucket{writer: writer, next: nil}
	if lm.firstBucket == nil {
		lm.firstBucket = b
		lm.lastBucket = b
		lm.logQueue = make(chan *logEntry, 32)
		go lm.processLogQueue()
	} else {
		lm.lastBucket.next = b
		lm.lastBucket = b
	}
}

func AddWriter(writer LogWriter) {
	DefaultLogManager.AddWriter(writer)
}

func (lm *LogManager) SetLogLevel(level LogLevel) {
	lm.logLevel = level
}

func SetLogLevel(level LogLevel) {
	DefaultLogManager.SetLogLevel(level)
}

func (lm *LogManager) Log(level LogLevel, format string, v ...interface{}) {
	if level >= lm.logLevel {
		lm.logQueue <- &logEntry{level, fmt.Sprintf(format, v...)}
	}
}

func Log(level LogLevel, format string, v ...interface{}) {
	DefaultLogManager.Log(level, format, v...)
}

func (lm *LogManager) processLogQueue() {
	for {
		entry := <-lm.logQueue

		for iter := lm.firstBucket; iter != nil; iter = iter.next {
			iter.writer.LogWrite(entry.level, entry.text)
		}
	}
}

func (lm *LogManager) Sync() {
	for len(lm.logQueue) > 0 {
	}
}

func Sync() {
	DefaultLogManager.Sync()
}

func Terminate() {
	Log(INFO, "Terminating.")
	DefaultLogManager.Sync()
	os.Exit(1)
}

func Warning(format string, v ...interface{}) {
	Log(WARNING, format, v...)
}

func Error(format string, v ...interface{}) {
	Log(ERROR, format, v...)
}
func Fatal(format string, v ...interface{}) {
	Log(ERROR, format, v...)
	Terminate()
}

func Info(format string, v ...interface{}) {
	Log(INFO, format, v...)
}

func Debug(format string, v ...interface{}) {
	Log(DEBUG, format, v...)
}

func DebugX(format string, v ...interface{}) {
	Log(DEBUGX, format, v...)
}

func DebugXX(format string, v ...interface{}) {
	Log(DEBUGXX, format, v...)
}