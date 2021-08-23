package common

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

type Color string

const (
	ResetColor  Color = "\033[0m"
	Black             = "\033[0;30m"
	Red               = "\033[0;31m"
	Green             = "\033[0;32m"
	Gold              = "\033[0;33m"
	Blue              = "\033[0;34m"
	Purple            = "\033[0;35m"
	Cyan              = "\033[0;36m"
	LightGray         = "\033[0;37m"
	DarkGray          = "\033[1;30m"
	LightRed          = "\033[1;31m"
	LightGreen        = "\033[1;32m"
	Yellow            = "\033[1;33m"
	LightBlue         = "\033[1;34m"
	LightPurple       = "\033[1;35m"
	LightCyan         = "\033[1;36m"
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
	return fmt.Fprintln(w, c, a, ResetColor)
}

func RandomColor() Color {
	index := randObj.Intn(len(colors))
	return colors[index]
}
