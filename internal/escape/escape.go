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

//type Logger

func NewLogger() *log.Logger {
	logFile, err := os.OpenFile("/var/log/ont.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Error")
	}
	defer logFile.Close()

	logger := log.New(logFile, "PREFIX:", log.Ldate|log.Ltime|log.Lshortfile)

	return logger

}

func LogPrint(value ...any) {
	logger := NewLogger()

	logger.Print(value...)
}

func LogPrintf(format string, value ...any) {
	logger := NewLogger()
	logger.Printf(format, value)

}

func LogFatal(value ...any) {
	LogPrint(value)
	os.Exit(1)
}
