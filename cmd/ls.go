package cmd

import (
	"fmt"
	"github.com/jedthehumanoid/cardcabinet"
	"github.com/spf13/cobra"
	"path/filepath"
)

var query string
var board string

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "Query")
	lsCmd.PersistentFlags().StringVarP(&board, "board", "b", "", "Show cards in board")
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List cards, with optional filter",
	Run: func(cmd *cobra.Command, args []string) {
		ls(args)
	},
}

func ls(args []string) {
	src := "."
	if len(args) > 0 {
		src = args[0]

	}
	src = filepath.Clean(src) + "/"
	cards := cardcabinet.ReadCards(src)
	boards := cardcabinet.ReadBoards(src)

	if board != "" {
		for _, b := range boards {
			if b.Name != board {
				continue
			}

			for _, deck := range b.Decks {
				fmt.Printf("%s:\n", deck.Name)
				if query == "" {
					for _, card := range deck.Get(cards) {
						listCard(card, config)
					}
				} else {
					for _, card := range cardcabinet.QueryCards(deck.Get(cards), query) {
						listCard(card, config)
					}
				}
				fmt.Println()
			}
		}
	}

	if board == "" && query != "" {
		cards = cardcabinet.QueryCards(cards, query)
		for _, card := range cards {
			listCard(card, config)
		}
	}
	if board == "" && query == "" {
		for _, card := range cards {
			listCard(card, config)
		}
	}

}
