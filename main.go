package main

import (
	"card-cabinet"
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

func getArguments() (string, string, []string, string, []string) {
	filter := ""
	command := ""
	arguments := []string{}
	flags := []string{}
	selected := ""

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
		_, err := strconv.ParseInt(arg, 10, 32)
		if err == nil {
			selected = arg
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

func getLabels(cards []cardcabinet.Card) []cardcabinet.Board {
	boards := []cardcabinet.Board{}
	for _, label := range cardcabinet.GetLabels(cards) {
		deck := cardcabinet.Deck{}
		deck.Title = label
		deck.Labels = []string{label}

		board := cardcabinet.Board{}
		board.Title = "+" + label
		board.Decks = append(board.Decks, deck)

		boards = append(boards, board)
	}
	return boards
}

func main() {
	config := loadConfig("cabinet.toml")
	filter, command, _, selected, _ := getArguments()

	dir := "."
	if config.Src != "" {
		dir = config.Src
	}
	dir = filepath.Clean(dir) + "/"

	cards := cardcabinet.ReadCards(dir)
	boards := cardcabinet.ReadBoards(dir)

	boards = append(boards, getLabels(cards)...)

	deck := cardcabinet.Deck{}
	deck.Cards = cards

	board := cardcabinet.Board{}
	board.Decks = append(board.Decks, deck)
	boards = append(boards, board)

	for i, board := range boards {
		for i, deck := range board.Decks {
			if len(deck.Labels) > 0 {
				deck.Cards = cardcabinet.FilterLabels(cards, deck.Labels)
			}
			board.Decks[i] = deck
		}
		boards[i] = board
	}

	if selected != "" {
		i, err := strconv.Atoi(selected)
		if err == nil && i <= len(cards) {
			cards = cards[i-1 : i]
		}
	}

	switch command {
	case "boards", "b":
		listBoards(boards)
	case "list":
		listBoard(boards, filter)
	case "cat", "c":
		catCards(cards)
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

func catCards(cards []cardcabinet.Card) {
	for _, card := range cards {
		fmt.Println("\n\n\n" + card.Title)
		fmt.Println(cardcabinet.MarshalFrontmatter(card))
		fmt.Println(card.Contents)
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
