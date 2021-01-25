package cmd

import (
"fmt"
	"github.com/spf13/cobra"
		"github.com/jedthehumanoid/cardcabinet"
)

func init(){
	rootCmd.AddCommand(cardsCmd)
}
var cardsCmd = &cobra.Command{
	Use: "cards",
	Short: "Get cards",
	Run: func(cmd *cobra.Command, args []string) {
		cards(args)
	},
}
func cards(args []string) {
	cards := cardcabinet.ReadCards(config.Src)
	for _, card := range cards {
		fmt.Println(card.Name)
	}
}
