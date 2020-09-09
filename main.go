package main

import (
	"fmt"
	"github.com/jedthehumanoid/card-cabinet"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const defaultcommand = "list"
const gray = "\033[38;2;100;100;100m"
const darkgray = "\033[38;2;50;50;50m"
const yellow = "\033[38;2;250;189;47m"
const reset = "\033[0m"

type Config struct {
	Src    string            `toml:"src"`
	Colors map[string]string `toml:"colors"`
}

func getArguments() (string, string, []string, []string) {
	command := defaultcommand
	board := "."

	arguments, flags := extractPrefix(os.Args[1:], "-")

	if len(arguments) > 0 {
		board = arguments[0]
		arguments = arguments[1:]
	}

	if len(arguments) > 0 {
		command = arguments[0]
		arguments = arguments[1:]
	}
	return board, command, arguments, flags
}

func loadConfig(file string) Config {
	var config Config
	err := loadToml(file, &config)
	if err != nil {
		fmt.Println("Couldn't load configuration file")
	}
	return config
}

func main() {
	config := loadConfig("cabinet.toml")
	b, command, args, _ := getArguments()

	if config.Src == "" {
		config.Src = "."
	}

	config.Src = filepath.Clean(config.Src) + "/"

	cards := cardcabinet.ReadCards(config.Src)
	boards := cardcabinet.ReadBoards(config.Src)

	for _, card := range cards {
		fmt.Println(card.Name)
	}

	fmt.Println(boards)

	if b == "." {
		b = ""
	}

	board := cardcabinet.GetBoard(boards, config.Src+b)

	if command == "search" || command == "s" {
		re := regexp.MustCompile("(?i)" + strings.Join(args[:len(args)-1], ".*"))
		command = args[len(args)-1]

		temp := []cardcabinet.Card{}
		for _, card := range cards {
			if re.MatchString(card.Name) {
				temp = append(temp, card)
			}
		}
		cards = temp
	}

	number, err := strconv.Atoi(command)
	if err == nil {
		cards = cards[number-1 : number]
		command = args[0]
	}

	switch command {
	case "boards", "b":
		listBoards(boards, config)
	case "list", "ls":
		listBoard(cards, board, config)
	case "cat", "c":
		catCards(cards, board)
	case "filename", "f":
		names(cards, board, config)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
