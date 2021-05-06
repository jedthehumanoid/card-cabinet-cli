package cmd

import (
	"fmt"
	"github.com/jedthehumanoid/cardcabinet"
	"github.com/spf13/cobra"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(boardsCmd)
}

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Get boards",
	Run: func(cmd *cobra.Command, args []string) {
		boards(args)
	},
}

func boards(args []string) {

	src := "."
	if len(args) > 0 {
		src = args[0]

	}
	src = filepath.Clean(src) + "/"
	cards := cardcabinet.ReadCards(src)
	boards := cardcabinet.ReadBoards(src)
	for _, board := range boards {
		fmt.Println(board.Name)

		fmt.Println()

		fmt.Printf("%+v\n", board)
		for _, deck := range board.Decks {
			fmt.Println(deck.Name)
			fmt.Println("----")
			fmt.Println()
			for _, card := range deck.Get(cards) {
				listCard(card, config)
			}
		}
	}
}
