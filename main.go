package main

import (
	"card-cabinet"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultcommand = "list"
const gray = "\033[38;2;100;100;100m"
const darkgray = "\033[38;2;50;50;50m"
const yellow = "\033[38;2;250;189;47m"
const reset = "\033[0m"

type Config struct {
	Src string `toml:"src"`
}

var cards []cardcabinet.Card

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
	var boards []cardcabinet.Board
	cards, boards = cardcabinet.ReadDir(dir)

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
		listBoard(board)
	case "cat", "c":
		catCards(board)
	case "filename", "f":
		names(board, dir)
	case "addlabel":
		fmt.Println("add label")
	case "removelabel":
		fmt.Println("remove label")
	case "edit", "e":
		edit(board, dir)
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func listBoard(board cardcabinet.Board) {
	i := 1
	for _, deck := range board.Decks {
		if deck.Title != "" {
			fmt.Println(deck.Title)
			fmt.Println(gray + fill("\u2500", len(deck.Title)) + reset)

		}
		for _, card := range deck.Cards {
			fmt.Printf("%d) ", i)
			listCard(getCard(cards, card))
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

func edit(board cardcabinet.Board, dir string) {
	args := []string{}
	for _, deck := range board.Decks {
		for _, card := range deck.Cards {
			args = append(args, dir+card)
		}
	}

	cmd := exec.Command("emacs", args...)
	cmd.Start()
}

func catCards(board cardcabinet.Board) {
	columns := getColumns()
	fmt.Println()
	for _, deck := range board.Decks {
		for _, title := range deck.Cards {
			card := getCard(cards, title)
			fmt.Println(darkgray + "\u250c" + fill("\u2500", columns-2) + "\u2510" + reset)
			fmt.Println(darkgray + "\u2502 " + reset + card.Title)
			fmt.Println(darkgray + "\u251c" + fill("\u2500", columns-2) + "\u2524" + reset)
			if cardcabinet.MarshalFrontmatter(card, false) != "" {
				for _, line := range strings.Split(cardcabinet.MarshalFrontmatter(card, false), "\n") {
					fmt.Println(darkgray + "\u2502 " + gray + line + reset)
				}
			}
			if card.Contents != "" {
				for _, line := range strings.Split(card.Contents, "\n") {
					for _, line := range splitlen(line, columns-2) {
						fmt.Println(darkgray + "\u2502 " + reset + line)
					}
				}
			}
			fmt.Println(darkgray + "\u2514" + fill("\u2500", columns-2) + "\u2518" + reset)
		}
	}
}

func names(board cardcabinet.Board, dir string) {
	for _, deck := range board.Decks {
		for _, card := range deck.Cards {
			fmt.Printf("%s%s\n", dir, card)
		}
	}
}

func getCard(cards []cardcabinet.Card, title string) cardcabinet.Card {
	for _, card := range cards {
		if card.Title == title {
			return card
		}
	}
	return cardcabinet.Card{}
}

func listCard(card cardcabinet.Card) {

	fmt.Printf("%s", card.Title)
	if card.Contents != "" {
		fmt.Print(yellow + " \u2261" + reset)
	}
	fmt.Print(gray)
	for _, label := range asStringSlice(card.Properties["labels"]) {
		fmt.Printf(" [%s]", label)
	}
	fmt.Println(reset)
}
