package common

import "fmt"

func output(c chan string, level string, format string, args ...interface{}) {
	fullFormat := fmt.Sprintf("[%-10v] %v", level, format)
	c <- fmt.Sprintf(fullFormat, args...)
}

func Info(c chan string, format string, args ...interface{}) {
	output(c, "INFO", format, args...)
}

func Debug(c chan string, format string, args ...interface{}) {
	output(c, "DEBUG", format, args...)
}

func Error(c chan string, format string, args ...interface{}) {
	output(c, "ERROR", format, args...)
}
