package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func init() {
	// file, err := os.OpenFile("myLOG.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Create a new log.Logger that writes to the standard output
	// ConsoleLog = log.New(os.Stdout, "", 0)
}
func Info(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)

	var formattedParams []interface{}
	for _, param := range v {
		if obj, ok := param.(map[string]interface{}); ok {
			jsonParam, err := json.Marshal(obj)
			if err != nil {
				log.Println(err)
				continue
			}
			formattedParams = append(formattedParams, string(jsonParam))
		} else {
			formattedParams = append(formattedParams, param)
		}
	}

	InfoLog.Printf("%s:%d: %s", filepath.Base(file), line, fmt.Sprint(formattedParams...))
}
func Warning(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	WarningLog.Printf("%s:%d: %s", filepath.Base(file), line, fmt.Sprint(v...))
}

func Error(v ...any) {
	_, file, line, _ := runtime.Caller(1)
	ErrorLog.Printf("%s:%d: %s", filepath.Base(file), line, fmt.Sprint(v...))
}
