package cmd

import (
	"fmt"
	"github.com/jedthehumanoid/cardcabinet"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(boardsCmd)
	boardsCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", true, "Recurse into subdirectories")

}

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "List boards",
	Run: func(cmd *cobra.Command, args []string) {
		boards(args)
	},
}

func boards(args []string) {
	if len(args) > 0 {
		config.Src = args
	}
	
	boards := []cardcabinet.Board{}
	for _, src := range config.Src {
		boards = append(boards, cardcabinet.ReadBoards(src, recursive)...)
	}
	//boards := cardcabinet.ReadBoards(config.Src, recursive)
	for _, board := range boards {
		fmt.Println(board.Name)

	}
}
