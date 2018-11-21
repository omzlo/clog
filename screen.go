package clog

import (
	"log"
	"os"
)

var (
	PlainTerminal = NewScreenLogWriter(false)
	ColorTerminal = NewScreenLogWriter(true)
)

type ScreenLogWriter struct {
	useColors bool
}

type FileLogWriter struct {
	filelog *log.Logger
}

var color_tags = [...]string{
	"\033[90mDEBUG++\033[0m",
	"\033[36mDEBUG+\033[0m",
	"\033[96mDEBUG\033[0m",
	"\033[92mINFO\033[0m",
	"\033[93mWARNING\033[0m",
	"\033[91mERROR\033[0m",
}

var plain_tags = [...]string{
	"DEBUG++",
	"DEBUG+",
	"DEBUG",
	"INFO",
	"WARNING",
	"ERROR",
}

func NewScreenLogWriter(with_colors bool) *ScreenLogWriter {
	return &ScreenLogWriter{useColors: with_colors}
}

func (slw *ScreenLogWriter) LogWrite(level LogLevel, text string) {
	var tag string

	if slw.useColors {
		tag = color_tags[level] + " "
	} else {
		tag = plain_tags[level] + " "
	}

	log.Println(tag + text)
}

func NewFileLogWriter(fname string) *FileLogWriter {
	logfile, err := os.Create(fname)
	if err != nil {
		return nil
	}
	filelog := log.New(logfile, "", log.LstdFlags)
	filelog.Println("**************** START *****************")
	return &FileLogWriter{filelog: filelog}
}

func (flw *FileLogWriter) LogWrite(level LogLevel, text string) {
	flw.filelog.Printf(plain_tags[level] + " " + text)
}
