package lib

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "API: ", log.Ldate|log.Ltime|log.Lshortfile)

func Infof(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	Logger.Printf("ERROR: "+format, v...)
}

func Fatalf(format string, v ...interface{}) {
	Logger.Fatalf("FATAL: "+format, v...)
}
