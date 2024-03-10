package escape

import (
	"fmt"
	"log"
	"os"
)

func Error(err string) {
	fmt.Println(err)
	os.Exit(1)
}

/*
func ErrorWithZeroRC(err string) {
	fmt.Println(err)
	os.Exit(0)
}*/

func Init() *os.File {
	logFile, err := os.OpenFile("/var/log/ont.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Error(err.Error())
	}

	return logFile
}

func NewLogger() (*log.Logger, *os.File) {

	// logFile, err := os.OpenFile("/var/log/ont.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	// if err != nil {
	// 	Error(err.Error())
	// }

	logFile := Init()

	logger := log.New(logFile, "", log.Ldate|log.Ltime)

	return logger, logFile

}

func LogPrint(value ...any) {
	logger, logFile := NewLogger()
	logger.Println(value...)
	logFile.Close()
}

func LogPrintf(format string, value ...any) {
	logger, logFile := NewLogger()
	logger.Printf(format, value...)
	logFile.Close()
}

func LogFatal(value ...any) {
	LogPrint(value...)
	os.Exit(1)
}
