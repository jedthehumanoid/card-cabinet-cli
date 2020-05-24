package main

import (
	"fmt"
	"github.com/jedthehumanoid/card-cabinet"
	"strconv"
	"strings"
)

func listBoard(cards []cardcabinet.Card, board cardcabinet.Board, config Config) {
	i := 1

	for _, deck := range board.Decks {
		if deck.Name != "" {
			fmt.Println(deck.Name)
			fmt.Println(gray + fill("\u2500", len(deck.Name)) + reset)

		}

		//	fmt.Println(board.Cards(cards))
		for _, card := range deck.Get(board.Cards(cards)) {
			fmt.Printf("%d) ", i)
			listCard(card, config)
			i++
		}
		fmt.Println()
	}
}

func listBoards(boards []cardcabinet.Board, config Config) {
	for _, board := range boards {
		name := strings.TrimPrefix(board.Name, config.Src)
		if name == "" {
			name = "/"
		}
		fmt.Println(name)
	}
}

func catCards(cards []cardcabinet.Card, board cardcabinet.Board) {
	columns := getColumns()
	fmt.Println()
	for _, deck := range board.Decks {
		for _, card := range deck.Get(cards) {
			fmt.Println(darkgray + "\u250c" + fill("\u2500", columns-2) + reset)
			fmt.Println(darkgray + "\u2502 " + yellow + card.Name + reset)
			fmt.Println(darkgray + "\u251c" + fill("\u2500", columns-2) + reset)
			if card.MarshalFrontmatter(false) != "" {
				for _, line := range strings.Split(card.MarshalFrontmatter(false), "\n") {
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
			fmt.Println(darkgray + "\u2514" + fill("\u2500", columns-2) + reset)
		}
	}
}

func names(cards []cardcabinet.Card, board cardcabinet.Board, config Config) {
	for _, deck := range board.Decks {
		for _, name := range deck.Get(cards) {
			fmt.Printf("%s%s\n", config.Src, name)
		}
	}
}

func listCard(card cardcabinet.Card, config Config) {

	card.Name = strings.TrimPrefix(card.Name, config.Src)

	tokens := strings.Split(card.Name, "/")

	p := strings.Join(tokens[:len(tokens)-1], "/")
	if p != "" {
		p = gray + "[/" + p + "] " + reset
	}

	title := tokens[len(tokens)-1]

	title = FromSlug(title)
	title = strings.ToUpper(title[:1]) + title[1:]
	title = strings.TrimSuffix(title, ".md")
	fmt.Printf("%s%s", p, title)
	if card.Contents != "" {
		fmt.Print(yellow + " \u2261" + reset)
	}

	for _, label := range asStringSlice(card.Properties["labels"]) {
		c, hascolor := config.Colors[label]
		if hascolor {
			fmt.Printf(color(c))
		} else {
			fmt.Printf(gray)
		}
		fmt.Printf(" [%s]", label)
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
