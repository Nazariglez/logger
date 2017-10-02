// Created by nazarigonzalez on 1/10/17.

package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	LOG
	WARN
	ERROR
	FATAL
)

type Logger struct {
	mu sync.Mutex

	logger    *log.Logger
	loggerErr *log.Logger

	level LogLevel

	fileLevel  LogLevel
	fileTime   time.Time
	filePath   string
	fileName   string
	fileLogger *log.Logger
	logFile    *os.File

	isAsync bool
	queue   chan func()
}

func New() *Logger {
	flags := log.Flags()

	return &Logger{
		level:     LOG,
		logger:    log.New(os.Stdout, "", flags),
		loggerErr: log.New(os.Stderr, "", flags),
	}
}

func NewAsync() *Logger {
	l := New()
	l.isAsync = true
	l.queue = make(chan func())
	go l.readQueue()
	return l
}

func (l *Logger) readQueue() {
	for f := range l.queue {
		f()
	}
}

func (l *Logger) SetLevel(level LogLevel) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.level = level
	return l
}

func (l *Logger) GetLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *Logger) EnableFileOutput(name, directory string, level LogLevel) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.fileName = name
	l.filePath = directory
	l.fileLevel = level

	err := l.checkCurrentFile()
	if err != nil {
		return err
	}

	return nil
}
func (l *Logger) DisableFileOutput() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logFile != nil {
		err := l.logFile.Close()
		if err != nil {
			return err
		}

		l.logFile = nil
	}

	return nil
}

func (l *Logger) checkCurrentFile() error {
	now := time.Now()

	if l.logFile != nil {

		//create a new file for the new day
		if now.Day() != l.fileTime.Day() {
			err := l.logFile.Close()
			if err != nil {
				return err
			}

			err = l.initLogFile()
			if err != nil {
				return err
			}
		}

		return nil
	}

	//initialize logFile
	err := l.initLogFile()
	if err != nil {
		return err
	}

	return nil
}

func (l *Logger) initLogFile() error {
	now := time.Now()

	name := l.fileName + "." + now.Format("20060102") + ".log"
	f, err := os.OpenFile(
		path.Join(l.filePath, name),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0640,
	)

	if err != nil {
		return err
	}

	l.logFile = f
	l.fileTime = now

	l.fileLogger = log.New(f, "", log.Flags())

	return nil
}

func (l *Logger) msg(text string, level LogLevel, isErr bool) {
	if l.isAsync {
		l.queue <- func(){
			l.sendMsg(text, level, isErr)
		}

		return
	}

	l.sendMsg(text, level, isErr)
}

func (l *Logger) sendMsg(text string, level LogLevel, isErr bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.level <= level {
		if isErr {
			l.loggerErr.Println(getPrefix(level, false) + text)
		} else {
			l.logger.Println(getPrefix(level, false) + text)
		}
	}

	if l.logFile != nil && l.fileLevel <= level {
		err := l.checkCurrentFile()

		if err != nil {
			l.logger.Println(getPrefix(WARN, false) + "CANNOT SAVE LOG TO FILE - " + err.Error())
			return
		}

		l.fileLogger.Println(getPrefix(level, true) + text)
	}
}

func (l *Logger) Trace(text string) *Logger {
	l.msg(text, TRACE, false)
	return l
}

func (l *Logger) Tracef(format string, v ...interface{}) *Logger {
	l.msg(fmt.Sprintf(format, v...), TRACE, false)
	return l
}

func (l *Logger) Debug(text string) *Logger {
	l.msg(text, DEBUG, false)
	return l
}

func (l *Logger) Debugf(format string, v ...interface{}) *Logger {
	l.msg(fmt.Sprintf(format, v...), DEBUG, false)
	return l
}

func (l *Logger) Info(text string) *Logger {
	l.msg(text, INFO, false)
	return l
}

func (l *Logger) Infof(format string, v ...interface{}) *Logger {
	l.msg(fmt.Sprintf(format, v...), INFO, false)
	return l
}

func (l *Logger) Log(text string) *Logger {
	l.msg(text, LOG, false)
	return l
}

func (l *Logger) Logf(format string, v ...interface{}) *Logger {
	l.msg(fmt.Sprintf(format, v...), LOG, false)
	return l
}

func (l *Logger) Warn(text string) *Logger {
	l.msg(text, WARN, false)
	return l
}

func (l *Logger) Warnf(format string, v ...interface{}) *Logger {
	l.msg(fmt.Sprintf(format, v...), WARN, false)
	return l
}

func (l *Logger) Error(text string) *Logger {
	l.msg(text, ERROR, true)
	return l
}

func (l *Logger) Errorf(format string, v ...interface{}) *Logger {
	l.msg(fmt.Sprintf(format, v...), ERROR, true)
	return l
}

func (l *Logger) Fatal(text string) *Logger {
	l.msg(text, FATAL, true)
	return l
}

func (l *Logger) Fatalf(format string, v ...interface{}) *Logger {
	l.msg(fmt.Sprintf(format, v...), FATAL, true)
	return l
}

//
var levelPrefix = map[LogLevel]string{
	TRACE: "Trace: ",
	DEBUG: "Debug: ",
	INFO:  "Info: ",
	LOG:   "Log: ",
	WARN:  "Warn: ",
	ERROR: "Error: ",
	FATAL: "Fatal: ",
}

func getPrefix(level LogLevel, file bool) string {
	prefix := levelPrefix[level]
	if file {
		return prefix
	}

	switch level {
	case TRACE:
		prefix = "\033[35m" + prefix + "\033[39m"
	case DEBUG:
		prefix = "\033[96m" + prefix + "\033[39m"
	case INFO:
		prefix = "\033[32m" + prefix + "\033[39m"
	case WARN:
		prefix = "\033[93m" + prefix + "\033[39m"
	case ERROR:
		prefix = "\033[31m" + prefix + "\033[39m"
	case FATAL:
		prefix = "\033[0;41m" + prefix[:len(prefix)-1] + "\033[0;39m "
	}

	return prefix
}
