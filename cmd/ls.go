package cmd

import (
"fmt"
	"github.com/spf13/cobra"
	"github.com/jedthehumanoid/cardcabinet"
	"strings"
)

func init(){
	rootCmd.AddCommand(lsCmd)
}
var lsCmd = &cobra.Command{
	Use: "ls",
	Short: "List cards, with optional filter",
	Run: func(cmd *cobra.Command, args []string) {
		ls(args)
	},
}
func ls(args []string) {
	cards := cardcabinet.ReadCards(config.Src)
	ret := cardcabinet.QueryCards(cards, strings.Join(args, " "))
	for _, card := range ret {
		fmt.Println(card.Name)
	}
}
