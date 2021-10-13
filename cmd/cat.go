package cmd

import (
	"card-cabinet-cli/config"
	"card-cabinet-cli/view"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedthehumanoid/cardcabinet"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(catCmd)
	catCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", true, "Recurse into subdirectories")
	catCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "Query")
}

var catCmd = &cobra.Command{
	Use:   "cat [DIR]",
	Short: "Cat cards",
	Run: func(cmd *cobra.Command, args []string) {
		cat(args)
	},
}

func cat(args []string) {
	if len(args) == 0 {
		args = []string{filepath.Clean(".") + "/"}
	}

	fi, err := os.Stat(args[0])
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	if fi.Mode().IsDir() {
		cards := cardcabinet.ReadCards(args[0], recursive)
		for _, card := range cardcabinet.QueryCards(cards, query) {
			view.List(card, cfg)
		}
		fmt.Println()
	} else if strings.HasSuffix(args[0], ".board.toml") {
		b, err := cardcabinet.ReadBoard(args[0])
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		catboard(b, cfg)
	}
}

func catboard(b cardcabinet.Board, config config.Config) {
	cards := cardcabinet.ReadCards(b.Path(), recursive)
	for _, deck := range b.Decks {
		fmt.Printf("[%s]\n", deck.Name)
		for _, card := range cardcabinet.QueryCards(deck.Get(cards), query) {
			view.Cat(card, config)
		}
		fmt.Println()
	}
}
