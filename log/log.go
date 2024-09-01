package log

import (
	"fmt"
	"log"
	"strings"
)

var (
	loglevel int
)

func init(){
	log.SetFlags(5)
}

func SetLevel(level string) error {
	switch strings.ToLower(level) {
	case "trace":
		loglevel = 1
	case "debug":
		loglevel = 2
	case "info":
		loglevel = 3
	case "warn":
		loglevel = 4
	case "error":
		loglevel = 5
	default:
		return fmt.Errorf("Unknown log level: %s", level)
	}
	return nil
}

func Fatal(args ...any){
	log.Fatal(args...)
}

func Fatalf(format string, args ...any){
	log.Fatalf(format, args...)
}

func print(level string, args ...any){
	var logs []any
	logs = append(logs, level + ": ")
	logs = append(logs, args...)
	log.Print(logs...)
}

func printf(level string, format string, logs ...any){
	log.Printf(level + ": " + format, logs...)
}

func Trace(args ...any){
	if loglevel <= 1 && loglevel != 0 {
		print("Trace", args...)
	}
}

func Tracef(format string, logs ...any){
	if loglevel <= 1 && loglevel != 0 {
		printf("Trace", format, logs...)
	}
}

func Debug(args ...any){
	if loglevel <= 2 && loglevel != 0 {
		print("Debug", args...)
	}
}

func Debugf(format string, logs ...any){
	if loglevel <= 2 && loglevel != 0 {
		printf("Debug", format, logs...)
	}
}

func Info(args ...any){
	if loglevel <= 3 && loglevel != 0 {
		print("Info", args...)
	}
}

func Infof(format string, logs ...any){
	if loglevel <= 3 && loglevel != 0 {
		printf("Info", format, logs...)
	}
}

func Warn(args ...any){
	if loglevel <= 4 && loglevel != 0 {
		print("Warn", args...)
	}
}

func Warnf(format string, logs ...any){
	if loglevel <= 4 && loglevel != 0 {
		printf("Warn", format, logs...)
	}
}

func Error(args ...any){
	if loglevel <= 5 && loglevel != 0 {
		print("Error", args...)
	}
}

func Errorf(format string, logs ...any){
	if loglevel <= 5 && loglevel != 0 {
		printf("Error", format, logs...)
	}
}
