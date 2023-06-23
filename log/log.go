package log

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
)

type OutputType string

const (
	FILE   OutputType = "file"
	STDOUT OutputType = "out"
	STDERR OutputType = "error"
)

var (
	debugLogger   *log.Logger
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
)

func init() {
	if err := SetOutputType(STDOUT, ""); err != nil {
		panic(err)
	}
}

func SetOutputType(outputType OutputType, filePath string) error {
	switch outputType {
	case FILE:
		if len(filePath) == 0 {
			filePath = "./runtime.log"
		}
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return errors.Wrap(err, "create log file error")
		}
		editOutputType(file)
	case STDERR:
		editOutputType(os.Stderr)
	case STDOUT:
		editOutputType(os.Stdout)
	default:
		return fmt.Errorf("unknown output type: %v", outputType)
	}
	return nil
}

func editOutputType(out io.Writer) {
	debugLogger = log.New(out, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(out, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(out, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debug(format string, v ...any) {
	debugLogger.Printf(format, v...)
}

func Info(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

func Warn(format string, v ...any) {
	warningLogger.Printf(format, v...)
}

func Fatal(format string, v ...any) {
	errorLogger.Printf(format, v...)
}
