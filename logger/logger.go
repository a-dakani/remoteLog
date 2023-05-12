package logger

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	Red = iota + 31
	Green
	Yellow
)

func colorize(textColor int, text string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", textColor, text)
}

func now() time.Time {
	return time.Now()
}

func log(msg string, level string, color int) {
	preamble := colorize(color, fmt.Sprintf("%s %s", now().Format("2006-01-02 15:04:05"), level))
	fmt.Printf("%s: %s\n", preamble, msg)
}

func Fatal(msg string) {
	log(msg, "FATAL ERROR", Red)
	os.Exit(1)
}

func Warning(msg string) {
	log(msg, "WARNING", Yellow)
}

func Info(msg string) {
	log(msg, "INFO", Green)
}

func ProcessArgumentError() {
	flag.Usage()
	Fatal("Invalid arguments, Please check the arguments and try again")
}
