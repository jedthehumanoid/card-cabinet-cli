package ansi

import (
	"strconv"
	"fmt"
	"os/exec"
	"strings"
)

const Reset = "\033[0m"

func Color(hex string) string {
	r := hex[0:2]
	g := hex[2:4]
	b := hex[4:6]

	red, err := strconv.ParseInt(r, 16, 64)
	if err != nil {
		panic(err)
	}
	green, err := strconv.ParseInt(g, 16, 64)
	if err != nil {
		panic(err)
	}
	blue, err := strconv.ParseInt(b, 16, 64)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\033[38;2;%d;%d;%dm", red, green, blue)
	return ""
}

func GetColumns() int {
	cmd := exec.Command("tput", "cols")
	columns, _ := cmd.Output()
	ret, _ := strconv.Atoi(strings.TrimSpace(string(columns)))
	return ret
}
