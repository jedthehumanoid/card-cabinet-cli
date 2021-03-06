package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jedthehumanoid/cardcabinet"
	"card-cabinet-cli/tools"
)

func listBoard(cards []cardcabinet.Card, board cardcabinet.Board, config Config) {
	for _, deck := range board.Decks {
		if deck.Name != "" {
			fmt.Println(deck.Name)
			fmt.Println(gray + tools.Fill("\u2500", len(deck.Name)) + reset)

		}

		for _, card := range deck.Get(board.Cards(cards)) {
			listCard(card, config)
		}
		fmt.Println()
	}
}

func listBoards(boards []cardcabinet.Board, config Config) {
	for _, board := range boards {
		fmt.Println(board.Name)
	}
}

func catCards(cards []cardcabinet.Card, board cardcabinet.Board) {
	columns := tools.GetColumns()
	fmt.Println()
	for _, deck := range board.Decks {
		for _, card := range deck.Get(cards) {
			fmt.Println(darkgray + "\u250c" + tools.Fill("\u2500", columns-2) + reset)
			fmt.Println(darkgray + "\u2502 " + yellow + card.Name + reset)
			fmt.Println(darkgray + "\u251c" + tools.Fill("\u2500", columns-2) + reset)
			if card.MarshalFrontmatter(false) != "" {
				for _, line := range strings.Split(card.MarshalFrontmatter(false), "\n") {
					fmt.Println(darkgray + "\u2502 " + gray + line + reset)
				}
			}
			if card.Contents != "" {
				for _, line := range strings.Split(card.Contents, "\n") {
					for _, line := range tools.Splitlen(line, columns-2) {
						fmt.Println(darkgray + "\u2502 " + reset + line)
					}
				}
			}
			fmt.Println(darkgray + "\u2514" + tools.Fill("\u2500", columns-2) + reset)
		}
	}
}

func names(cards []cardcabinet.Card, board cardcabinet.Board, config Config) {
	for _, deck := range board.Decks {
		for _, card := range deck.Get(cards) {
			fmt.Printf("%s%s\n", config.Src, card.Name)
		}
	}
}

func listCard(card cardcabinet.Card, config Config) {
	tokens := strings.Split(card.Name, "/")

	p := strings.Join(tokens[:len(tokens)-1], "/")
	if p != "" {
		p = gray + "" + p + " " + reset
	}

	title := tokens[len(tokens)-1]

	title = tools.FromSlug(title)
	title = strings.ToUpper(title[:1]) + title[1:]
	title = strings.TrimSuffix(title, ".md")
	fmt.Printf("%s%s", p, title)
	if card.Contents != "" {
		fmt.Print(yellow + " \u2261" + reset)
	}

	for _, label := range tools.AsStringSlice(card.Properties["labels"]) {
		c, hascolor := config.Colors[label]
		if hascolor {
			fmt.Printf(color(c))
		} else {
			fmt.Printf(gray)
		}
		fmt.Printf(" #%s", label)
	}
	fmt.Println(reset)
}

func color(hex string) string {
	r := hex[0:2]
	g := hex[2:4]
	b := hex[4:6]

	red, err := strconv.ParseInt(r, 16, 64)
	if err != nil {
		panic(err)
	}
	green, err := strconv.ParseInt(g, 16, 64)
	blue, err := strconv.ParseInt(b, 16, 64)

	fmt.Printf("\033[38;2;%d;%d;%dm", red, green, blue)
	return ""
}
