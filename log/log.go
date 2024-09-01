package log

import (
	"fmt"
	"log"
	"strings"

	"crypto/rand"
	"encoding/hex"
)

var (
	loglevel int
)

type Log struct {
	id string
}

func init(){
	log.SetFlags(5)
}

func New() *Log {
	LogInstance := &Log{}
	LogInstance.Start()
	return LogInstance
}

func (l *Log) Start(){
	bytes := make([]byte, 6) // 每两个十六进制字符表示一个字节
	_, err := rand.Read(bytes)
	if err != nil {
		l.Fatal(err)
	}
	l.id = hex.EncodeToString(bytes)
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

func (l *Log) print(level string, args ...any){
	var logs []any
	logs = append(logs, level + ": [" + l.id + "] ")
	logs = append(logs, args...)
	log.Print(logs...)
}

func (l *Log) printf(level string, format string, logs ...any){
	log.Printf(level + ": [" + l.id + "] " + format, logs...)
}

func (l *Log) Fatal(args ...any){
	log.Fatal(args...)
}

func (l *Log) Fatalf(format string, args ...any){
	log.Fatalf(format, args...)
}

func (l *Log) Trace(args ...any){
	if loglevel <= 1 && loglevel != 0 {
		l.print("Trace", args...)
	}
}

func (l *Log) Tracef(format string, logs ...any){
	if loglevel <= 1 && loglevel != 0 {
		l.printf("Trace", format, logs...)
	}
}

func (l *Log) Debug(args ...any){
	if loglevel <= 2 && loglevel != 0 {
		l.print("Debug", args...)
	}
}

func (l *Log) Debugf(format string, logs ...any){
	if loglevel <= 2 && loglevel != 0 {
		l.printf("Debug", format, logs...)
	}
}

func (l *Log) Info(args ...any){
	if loglevel <= 3 && loglevel != 0 {
		l.print("Info", args...)
	}
}

func (l *Log) Infof(format string, logs ...any){
	if loglevel <= 3 && loglevel != 0 {
		l.printf("Info", format, logs...)
	}
}

func (l *Log) Warn(args ...any){
	if loglevel <= 4 && loglevel != 0 {
		l.print("Warn", args...)
	}
}

func (l *Log) Warnf(format string, logs ...any){
	if loglevel <= 4 && loglevel != 0 {
		l.printf("Warn", format, logs...)
	}
}

func (l *Log) Error(args ...any){
	if loglevel <= 5 && loglevel != 0 {
		l.print("Error", args...)
	}
}

func (l *Log) Errorf(format string, logs ...any){
	if loglevel <= 5 && loglevel != 0 {
		l.printf("Error", format, logs...)
	}
}
