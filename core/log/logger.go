package log

import (
	"fmt"
	"log"
	"os"
)

var Logger *logManager

type logManager struct {
	logHandler *log.Logger
}

func InitLogger() {
	logHandler := log.New(os.Stderr, "[dpm]", log.Ldate|log.Ltime|log.Llongfile)
	Logger = &logManager{logHandler: logHandler}
}

func (l *logManager) Info(msg interface{}) {
	l.logHandler.Output(2, fmt.Sprintln(msg))
}

func (l *logManager) Warn(msg interface{}) {
	l.logHandler.Output(2, fmt.Sprintln(msg))
}

func (l *logManager) Error(msg interface{}) {
	l.logHandler.Output(2, fmt.Sprintln(msg))
}

func (l *logManager) Fatal(msg interface{}) {
	l.logHandler.Output(2, fmt.Sprintln(msg))
	os.Exit(1)
}

func (l *logManager) Fatalf(format string, msg interface{}) {
	l.logHandler.Fatalf(format, msg)
	l.logHandler.Output(2, fmt.Sprintf(format, msg))
	os.Exit(1)
}
