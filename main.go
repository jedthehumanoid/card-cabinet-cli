package main

import (
	"card-cabinet"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
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
	filter := ""
	command := ""
	arguments := []string{}
	flags := []string{}
	selected := 0

	// Extract flags
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			arguments = append(arguments, arg)
		}
	}

	// Extract filter
	temp := []string{}
	for _, arg := range arguments {
		if strings.HasPrefix(arg, "@") ||
			strings.HasPrefix(arg, "+") ||
			strings.HasPrefix(arg, "/") ||
			strings.HasPrefix(arg, ":") {
			filter = arg[1:]
		} else {
			temp = append(temp, arg)
		}
	}
	arguments = temp

	// Extract selected
	temp = []string{}
	for _, arg := range arguments {
		num, err := strconv.Atoi(arg)
		if err == nil {
			selected = num
		} else {
			temp = append(temp, arg)
		}
	}
	arguments = temp

	if len(arguments) == 0 {
		command = defaultcommand
	} else {
		command = arguments[0]
		arguments = arguments[1:]
	}

	return filter, command, arguments, selected, flags
}

func loadConfig(file string) Config {
	var config Config

	d, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("no configuration file")
		return config
	}

	_, err = toml.Decode(string(d), &config)

	if err != nil {
		panic(err)
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
		listBoard(boards, filter)
	case "cat", "c":
		catCards(boards, filter)
	case "filename", "f":
		names(boards, filter, dir)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}

}

func listBoard(boards []cardcabinet.Board, filter string) {
	for _, board := range boards {
		if board.Title == filter {
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

func catCards(boards []cardcabinet.Board, filter string) {
	for _, board := range boards {
		if board.Title == filter {
			for _, deck := range board.Decks {
				for _, card := range deck.Cards {
					fmt.Println("\n" + card.Title)
					fmt.Println(dash(len(card.Title)))
					fmt.Println(gray + cardcabinet.MarshalFrontmatter(card) + reset)
					fmt.Println(card.Contents)
				}
			}
		}
	}
}

func names(boards []cardcabinet.Board, filter string, dir string) {
	for _, board := range boards {
		if board.Title == filter {
			for _, deck := range board.Decks {
				for _, card := range deck.Cards {
					fmt.Printf("%s%s\n", dir, card.Title)
				}
			}
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
