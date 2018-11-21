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

var LogLevelStrings = [...]string{
	"DEBUGXX",
	"DEBUGX",
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
	"NONE",
}

func (ll *LogLevel) Set(s string) error {
	for i, level := range LogLevelStrings {
		if s == level {
			*ll = LogLevel(i)
			return nil
		}
	}
	return fmt.Errorf("Unrecognized loglevel '%s'", s)
}

func (ll LogLevel) String() string {
	if ll > NONE {
		return "NONE"
	}
	return LogLevelStrings[ll]
}

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
	if level >= lm.logLevel && lm.firstBucket != nil {
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

func Fatal(format string, v ...interface{}) {
	Log(ERROR, format, v...)
	Terminate()
}

func (lm *LogManager) Warning(format string, v ...interface{}) {
	lm.Log(WARNING, format, v...)
}

func Warning(format string, v ...interface{}) {
	Log(WARNING, format, v...)
}

func (lm *LogManager) Error(format string, v ...interface{}) {
	lm.Log(ERROR, format, v...)
}

func Error(format string, v ...interface{}) {
	Log(ERROR, format, v...)
}

func (lm *LogManager) Info(format string, v ...interface{}) {
	lm.Log(INFO, format, v...)
}

func Info(format string, v ...interface{}) {
	Log(INFO, format, v...)
}

func (lm *LogManager) Debug(format string, v ...interface{}) {
	lm.Log(DEBUG, format, v...)
}

func Debug(format string, v ...interface{}) {
	Log(DEBUG, format, v...)
}

func (lm *LogManager) DebugX(format string, v ...interface{}) {
	lm.Log(DEBUGX, format, v...)
}

func DebugX(format string, v ...interface{}) {
	Log(DEBUGX, format, v...)
}

func (lm *LogManager) DebugXX(format string, v ...interface{}) {
	lm.Log(DEBUGXX, format, v...)
}

func DebugXX(format string, v ...interface{}) {
	Log(DEBUGXX, format, v...)
}
