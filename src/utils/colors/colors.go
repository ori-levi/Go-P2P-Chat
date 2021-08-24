package colors

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
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
