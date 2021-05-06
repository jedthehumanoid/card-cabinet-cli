package cmd

import (

	"github.com/spf13/cobra"
	"github.com/jedthehumanoid/cardcabinet"
	"path/filepath"
)
var query string

func init(){
	rootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().StringVarP(&query, "query", "q", "", "Query")
}
var lsCmd = &cobra.Command{
	Use: "ls",
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
	if query != "" {
		cards = cardcabinet.QueryCards(cards, query)
	}
	for _, card := range cards {
		listCard(card, config)
	}
	
}
