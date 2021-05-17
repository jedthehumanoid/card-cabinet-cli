package cmd

import (
	"fmt"
	"github.com/jedthehumanoid/cardcabinet"
	"github.com/spf13/cobra"
	"path/filepath"
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
		config.Src = args[0]
		config.Src = filepath.Clean(config.Src) + "/"

	}
	
	boards := cardcabinet.ReadBoards(config.Src, recursive)
	for _, board := range boards {
		fmt.Println(board.Name)

	}
}
