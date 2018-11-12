package main

import (
	"github.com/omzlo/clog"
)

func main() {
	clog.AddWriter(clog.ColorTerminal)
	clog.SetLogLevel(clog.DEBUGXX)

	clog.Log(clog.DEBUGXX, "0")
	clog.Log(clog.DEBUGX, "1")
	clog.Log(clog.DEBUG, "2")
	clog.Log(clog.INFO, "3")
	clog.Log(clog.WARNING, "4")
	clog.Log(clog.ERROR, "5")

	clog.Terminate()
}
