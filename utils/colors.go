package utils

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

type Color string

const (
	ResetColor  Color = "\033[0m"
	Black       Color = "\033[0;30m"
	Red         Color = "\033[0;31m"
	Green       Color = "\033[0;32m"
	Gold        Color = "\033[0;33m"
	Blue        Color = "\033[0;34m"
	Purple      Color = "\033[0;35m"
	Cyan        Color = "\033[0;36m"
	LightGray   Color = "\033[0;37m"
	DarkGray    Color = "\033[1;30m"
	LightRed    Color = "\033[1;31m"
	LightGreen  Color = "\033[1;32m"
	Yellow      Color = "\033[1;33m"
	LightBlue   Color = "\033[1;34m"
	LightPurple Color = "\033[1;35m"
	LightCyan   Color = "\033[1;36m"
)

var (
	colors  = []Color{Red, Green, Purple, LightBlue, Yellow}
	randObj = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func ColorSprintf(c Color, format string, args ...interface{}) string {
	fullFormat := fmt.Sprintf("%v%v%v", c, format, ResetColor)
	return fmt.Sprintf(fullFormat, args...)
}

func ColorFprintln(w io.Writer, c Color, a ...interface{}) (n int, err error) {
	buff := []interface{}{c}
	buff = append(buff, a...)
	buff = append(buff, ResetColor)

	str := fmt.Sprint(buff...)
	return fmt.Fprintln(w, strings.Trim(str, " "))
}

func RandomColor() Color {
	index := randObj.Intn(len(colors))
	return colors[index]
}
