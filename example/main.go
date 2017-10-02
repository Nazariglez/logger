// Created by nazarigonzalez on 2/10/17.

package main

import (
	"fmt"
	"github.com/nazariglez/logger"
)

func main() {
	//loggerExample()
	asyncLoggerExample()
}

func loggerExample() {
	l := logger.New()
	l.SetLevel(logger.DEBUG)                       //print in the terminal info level logs
	l.EnableFileOutput("test", "./", logger.TRACE) //print in the file trace level logs

	l.Trace("Find me in the log file")

	l.Debug("Hey i'm a debug message")
	l.Info("Hey i'm a debug message")
	l.Log("Hey i'm a normal message")
	l.Warn("Hey i'm a warn message")
	l.Error("Hey i'm a error message")
	l.Fatal("Hey i'm a fatal message")
}

func asyncLoggerExample() {
	l := logger.NewAsync() //async logger
	l.SetLevel(logger.TRACE)
	l.EnableFileOutput("test", "./", logger.TRACE)

	l.Trace("Find me in the log file")

	fmt.Println("I'm the first message in your terminal not the 2nd!")

	l.Debug("Hey i'm a debug message")
	l.Info("Hey i'm a debug message")
	l.Log("Hey i'm a normal message")
	l.Warn("Hey i'm a warn message")
	l.Error("Hey i'm a error message")
	l.Fatal("Hey i'm a fatal message")
}
