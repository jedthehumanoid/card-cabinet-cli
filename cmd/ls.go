package cmd

import (
	"fmt"
	"github.com/jedthehumanoid/cardcabinet"
	"github.com/spf13/cobra"
	"strings"
)

var query string

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", true, "Recurse into subdirectories")
	lsCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "Query")
}

var lsCmd = &cobra.Command{
	Use:   "ls [DIR]",
	Short: "List cards",
	Run: func(cmd *cobra.Command, args []string) {
		ls(args)
	},
}

func ls(args []string) {
	if len(args) > 0 {
		config.Src = args
	}

	for _, src := range config.Src {

		// Add folder
		cards := cardcabinet.ReadCards(src, recursive)

		if len(cards) > 0 {
			if len(config.Src) > 1 {
				fmt.Printf("%s:\n", src)
			}
			for _, card := range cardcabinet.QueryCards(cards, query) {
				listCard(card, config)
			}
			fmt.Println()
			continue
		}

		// Add board
		src = strings.TrimSuffix(src, ".board.toml")
		src = strings.TrimSuffix(src, ".board")
		b, err := cardcabinet.ReadBoard(src + ".board.toml")
		if err == nil {
			if len(config.Src) > 1 {
				fmt.Printf("%s:\n", src)
			}
			cards = cardcabinet.ReadCards(b.Path(), recursive)

			for _, deck := range b.Decks {
				fmt.Printf("[%s]\n", deck.Name)
				for _, card := range cardcabinet.QueryCards(deck.Get(cards), query) {
					listCard(card, config)
				}

				fmt.Println()
			}
			continue
		}
		
		// Add file
		card, err := cardcabinet.ReadCard(src)
		if err == nil {
			listCard(card, config)
			continue
		}
		
		fmt.Printf("%s: No such file or directory\n", src)

	}

}
