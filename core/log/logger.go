package log

import (
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
	l.logHandler.Println(msg)
}

func (l *logManager) Warn(msg interface{}) {
	l.logHandler.Println(msg)
}

func (l *logManager) Error(msg interface{}) {
	l.logHandler.Println(msg)
}

func (l *logManager) Fatal(msg interface{}) {
	l.logHandler.Fatalln(msg)
}

func (l *logManager) Fatalf(format string, msg interface{}) {
	l.logHandler.Fatalf(format, msg)
}
