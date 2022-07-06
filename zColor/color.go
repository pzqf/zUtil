package zColor

import (
	"fmt"
	"strings"
)

func Green(str string) string {
	return cliColorRender(str, 32, 0)
}

func LightGreen(str string) string {
	return cliColorRender(str, 32, 1)
}

func Cyan(str string) string {
	return cliColorRender(str, 36, 0)
}

func LightCyan(str string) string {
	return cliColorRender(str, 36, 1)
}

func Red(str string) string {
	return cliColorRender(str, 31, 0)
}

func LightRed(str string) string {
	return cliColorRender(str, 31, 1)
}

func Yellow(str string) string {
	return cliColorRender(str, 33, 0)
}

func Black(str string) string {
	return cliColorRender(str, 30, 0)
}

func DarkGray(str string) string {
	return cliColorRender(str, 30, 1)
}

func LightGray(str string) string {
	return cliColorRender(str, 37, 0)
}

func White(str string) string {
	return cliColorRender(str, 37, 1)
}

func Blue(str string) string {
	return cliColorRender(str, 34, 0)
}

func LightBlue(str string) string {
	return cliColorRender(str, 34, 1)
}

func Purple(str string) string {
	return cliColorRender(str, 35, 0)
}

func LightPurple(str string) string {
	return cliColorRender(str, 35, 1)
}

func Brown(str string) string {
	return cliColorRender(str, 33, 0)
}

func cliColorRender(str string, color int, weight int) string {
	var mo []string

	if weight > 0 {
		mo = append(mo, fmt.Sprintf("%d", weight))
	}
	if len(mo) <= 0 {
		mo = append(mo, "0")
	}
	return fmt.Sprintf("\033[%s;%dm"+str+"\033[0m", strings.Join(mo, ";"), color)
}
