package main

import (
	"fmt"
	"github.com/jedthehumanoid/card-cabinet"
	"os"
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
	board := ""

	arguments, flags := extractPrefix(os.Args[1:], "-")

	if len(arguments) > 0 {
		command = arguments[0]
		arguments = arguments[1:]
	}

	if len(arguments) > 0 {
		board = arguments[0]
		arguments = arguments[1:]
	}
	return command, board, arguments, flags
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
	command, b, _, _ := getArguments()

	if config.Src == "" {
		config.Src = "."
	}

	config.Src = filepath.Clean(config.Src) + "/"
	files := cardcabinet.FindFiles(config.Src)

	cards := cardcabinet.ReadCards(files)
	boards := cardcabinet.ReadBoards(files)

	for _, folder := range cardcabinet.GetFolders(cards) {
		var board cardcabinet.Board
		board.Name = folder + "/"
		deck := cardcabinet.Deck{}
		board.Decks = []cardcabinet.Deck{deck}
		boards = append(boards, board)
	}

	var board cardcabinet.Board

	if b == "" || b == "." {
		board = cardcabinet.Board{}
		deck := cardcabinet.Deck{}
		for _, card := range cards {
			deck.Names = append(deck.Names, card.Name)
		}
		board.Decks = []cardcabinet.Deck{deck}
	} else {
		board = cardcabinet.GetBoard(boards, config.Src+b)
	}

	fmt.Println(ToJSON(boards))

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
