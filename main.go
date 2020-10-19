package main

import (
	"fmt"
	"github.com/jedthehumanoid/card-cabinet"
	"os"
	"strings"
	"path/filepath"
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
	board := "/"

	arguments, flags := extractPrefix(os.Args[1:], "-")

	if len(arguments) > 0 {
		command = arguments[0]
		arguments = arguments[1:]
	}

	if len(arguments) > 0 {
		board = arguments[0]
		arguments = arguments[1:]
	}

	//	fmt.Println(command, board, arguments, flags)
	return command, board, arguments, flags
}

func loadConfig(file string) Config {
	var config Config
	loadToml(file, &config)
	return config
}

func main() {
	config := loadConfig("cabinet.toml")
	command, boardname, _, _ := getArguments()

	if config.Src == "" {
		config.Src = "."
	}

	fmt.Printf("Looking for cards in: %s\n", config.Src)

	fi, _ := os.Stdout.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		command = "f"
	}

	config.Src = filepath.Clean(config.Src) + "/"

	cards := cardcabinet.ReadCards(config.Src)
	boards := cardcabinet.ReadBoards(config.Src)

	for _, card := range cards {
		fmt.Println(card.Name)

		if card.Frontmatter != "" {
			fmt.Println(card.Properties)
			fmt.Println(card.Frontmatter)
		}
		if card.Contents != "" {
		if len(card.Contents) > 80 {
			fmt.Println(strings.Replace(card.Contents, "\n", "", -1)[:80])
		} else {
			fmt.Println(strings.Replace(card.Contents, "\n", "", -1))
		}
		}
		fmt.Println()
	}

	for _, board := range boards {
		fmt.Println(board.Name)
	}


	var board cardcabinet.Board

	for _, b := range boards {
		if b.Name == boardname {
			board = b
		}
	}

	switch command {
	case "boards", "b":
		listBoards(boards, config)
	case "list", "ls":
		listBoard(cards, board, config)
	case "filename", "f":
		cards = board.Cards(cards)
		for _, deck := range board.Decks {
			for _, card := range deck.Get(cards) {
				fmt.Printf("%s%s%s\n", config.Src, board.Path(), card.Name)
			}
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
