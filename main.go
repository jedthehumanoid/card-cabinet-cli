package main

import (
	"fmt"
	"github.com/jedthehumanoid/card-cabinet"
	"os"
	"path/filepath"
	"strconv"
)

const defaultcommand = "list"
const gray = "\033[38;2;100;100;100m"
const darkgray = "\033[38;2;50;50;50m"
const yellow = "\033[38;2;250;189;47m"
const reset = "\033[0m"

type Config struct {
	Src string `toml:"src"`
}

func getArguments() (string, string, []string, int, []string) {
	command := defaultcommand
	filter := ""
	selected := 0

	arguments, flags := extractPrefix(os.Args[1:], "-")
	arguments, filters := extractPrefix(arguments, ":")

	if len(filters) > 0 {
		filter = filters[0][1:]
	}

	// Extract selected
	temp := []string{}
	for _, arg := range arguments {
		num, err := strconv.Atoi(arg)
		if err == nil {
			selected = num
		} else {
			temp = append(temp, arg)
		}
	}
	arguments = temp

	if len(arguments) > 0 {
		command = arguments[0]
		arguments = arguments[1:]
	}

	return filter, command, arguments, selected, flags
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
	filter, command, _, selected, _ := getArguments()

	if config.Src == "" {
		config.Src = "."
	}
	config.Src = filepath.Clean(config.Src) + "/"

	cards, boards := cardcabinet.ReadDir(config.Src)

	if selected != 0 {
		for i, board := range boards {
			if board.Title == filter {
				id := 1
				for _, deck := range board.Decks {
					for _, card := range deck.Cards {
						if id == selected {
							deck.Cards = []string{card}
							board.Decks = []cardcabinet.Deck{deck}
						}
						id++
					}
				}
				boards[i] = board
			}
		}
	}

	board := cardcabinet.GetBoard(boards, filter)

	switch command {
	case "boards", "b":
		listBoards(boards)
	case "list":
		listBoard(cards, board)
	case "cat", "c":
		catCards(cards, board)
	case "filename", "f":
		names(board, config)
	case "addlabel":
		fmt.Println("add label")
	case "removelabel":
		fmt.Println("remove label")
	case "edit", "e":
		edit(board, config)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
