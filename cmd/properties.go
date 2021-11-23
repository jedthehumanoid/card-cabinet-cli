package cmd

import (
	"card-cabinet-cli/ansi"
	"encoding/json"
	"fmt"
	"github.com/jedthehumanoid/cardcabinet"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var expandall bool

const gray = "\033[38;2;100;100;100m"

func init() {
	rootCmd.AddCommand(propsCmd)
	//propsCmd.PersistentFlags().StringVarP(&board, "expand", "e", "", "Show keys for property")
	propsCmd.PersistentFlags().BoolVarP(&expandall, "expand-all", "e", true, "Show keys for all properties")
	propsCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", true, "Recurse into subdirectories")

}

var propsCmd = &cobra.Command{
	Use:   "properties [DIR]",
	Short: "List properties",
	Run: func(cmd *cobra.Command, args []string) {
		props(args)
	},
}

func props(args []string) {
	if len(args) == 0 {
		args = []string{filepath.Clean(".") + "/"}
	}

	fi, err := os.Stat(args[0])
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	cards := []cardcabinet.Card{}

	if fi.Mode().IsDir() {
		cards = cardcabinet.ReadCards(args[0], recursive)
	} else if strings.HasSuffix(args[0], ".board.toml") {
		b, err := cardcabinet.ReadBoard(args[0])
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		boardcards := cardcabinet.ReadCards(b.Path(), recursive)
		for _, deck := range b.Decks {
			cards = append(cards, deck.Get(boardcards)...)
		}
	}

	properties := map[string][]string{}

	for _, card := range cards {
		for key, value := range card.Properties {
			switch value.(type) {
			case string, int, bool:
				properties[key] = append(properties[key], toJSON(value))
			case []interface{}:
				for _, v := range value.([]interface{}) {
					properties[key] = append(properties[key], toJSON(v))
				}
			}
		}
	}

	for key, value := range properties {
		fmt.Println(key)
		if expandall {
			for _, v := range unique(value) {
				fmt.Print(gray + v + ansi.Reset + " ")
			}

			fmt.Println()
			fmt.Println()
		}

	}
}

func unique(ss []string) []string {
	ret := []string{}
	values := map[string]bool{}
	for _, val := range ss {
		_, exists := values[val]
		if !exists {
			values[val] = true
			ret = append(ret, val)
		}
	}
	return ret

}

// ToJSON return JSON representation of interface
func toJSON(in interface{}) string {
	b, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return string(b)
}
