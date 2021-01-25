package cmd

import (
"fmt"
	"github.com/spf13/cobra"
		"github.com/jedthehumanoid/cardcabinet"
)

func init(){
	rootCmd.AddCommand(boardsCmd)
}
var boardsCmd = &cobra.Command{
	Use: "boards",
	Short: "Get boards",
	Run: func(cmd *cobra.Command, args []string) {
		boards(args)
	},
}
func boards(args []string) {
	boards := cardcabinet.ReadBoards(config.Src)
	for _, board := range boards {
		fmt.Println(board.Name)
	}
}
