package utils

import (
	"flag"
	"fmt"
	"time"
)

const (
	Red = iota + 31
	Green
	Yellow
)

func Colorize(text string, textColor int) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", textColor, text)
}

func now() time.Time {
	return time.Now()
}

func log(msg string, level string, color int) {
	preamble := Colorize(fmt.Sprintf("%s %s", now().Format("2006-01-02 15:04:05"), level), color)
	fmt.Printf("%s: %s\n", preamble, msg)
}

func Fatal(msg string) {
	log(msg, "FATAL", Red)
}

func Warning(msg string) {
	log(msg, "WARNING", Yellow)
}

func Info(msg string) {
	log(msg, "INFO", Green)
}

func ProcessArgumentError() {
	Fatal("Invalid arguments, Please check the arguments and try again")
	flag.Usage()
}
