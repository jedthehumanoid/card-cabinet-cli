package main

import (
	"card-cabinet"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultcommand = "list"
const gray = "\033[38;2;127;127;127m"
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

	//fmt.Println("filter:", filter, "command:", command, "arguments:", arguments, "selected:", selected, "flags:", flags)
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

	dir := "."
	if config.Src != "" {
		dir = config.Src
	}
	dir = filepath.Clean(dir) + "/"
	boards := cardcabinet.ReadDir(dir)

	if selected != 0 {
		for i, board := range boards {
			if board.Title == filter {
				id := 1
				for _, deck := range board.Decks {
					for _, card := range deck.Cards {
						if id == selected {
							deck.Cards = []cardcabinet.Card{card}
						}
						id++
					}
					board.Decks = []cardcabinet.Deck{deck}
				}
				boards[i] = board
			}
		}
	}

	switch command {
	case "boards", "b":
		listBoards(boards)
	case "list":
		listBoard(getBoard(boards, filter))
	case "cat", "c":
		catCards(getBoard(boards, filter))
	case "filename", "f":
		names(getBoard(boards, filter), dir)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}

}

func getBoard(boards []cardcabinet.Board, board string) cardcabinet.Board {
	for _, b := range boards {
		if b.Title == board {
			return b
		}
	}
	return cardcabinet.Board{}
}

func listBoard(board cardcabinet.Board) {
	i := 1
	for _, deck := range board.Decks {
		if deck.Title != "" {
			fmt.Println(deck.Title)
			fmt.Println(dash(len(deck.Title)))

		}
		for _, card := range deck.Cards {
			fmt.Printf("%d) ", i)
			listCard(card)
			i++
		}
		fmt.Println()
	}
}

func listBoards(boards []cardcabinet.Board) {
	for _, board := range boards {
		if board.Title != "" {
			fmt.Println(board.Title)
		}
	}
}

func dash(len int) string {
	ret := ""
	for i := 0; i < len; i++ {
		ret += "-"
	}
	return ret
}

func catCards(board cardcabinet.Board) {
	for _, deck := range board.Decks {
		for _, card := range deck.Cards {
			fmt.Println("\n" + card.Title)
			fmt.Println(dash(len(card.Title)))
			fmt.Println(gray + cardcabinet.MarshalFrontmatter(card) + reset)
			fmt.Println(card.Contents)
		}
	}
}

func names(board cardcabinet.Board, dir string) {
	for _, deck := range board.Decks {
		for _, card := range deck.Cards {
			fmt.Printf("%s%s\n", dir, card.Title)
		}
	}
}

func listCard(card cardcabinet.Card) {
	fmt.Printf("%s", card.Title)
	if card.Contents != "" {
		fmt.Print(" []")
	}
	fmt.Print(gray)
	for _, label := range asStringSlice(card.Properties["labels"]) {
		fmt.Printf(" [%s]", label)
	}
	fmt.Println(reset)

}

func toJSON(i interface{}) string {
	b, _ := json.MarshalIndent(i, " ", "   ")
	return string(b)
}
