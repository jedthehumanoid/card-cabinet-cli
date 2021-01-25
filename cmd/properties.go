package cmd

import (
	"card-cabinet-cli/tools"
	"fmt"
	"github.com/jedthehumanoid/cardcabinet"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(propsCmd)
}

var propsCmd = &cobra.Command{
	Use:   "properties",
	Short: "List properties",
	Run: func(cmd *cobra.Command, args []string) {
		props(args)
	},
}

func props(args []string) {
	cards := cardcabinet.ReadCards(config.Src)
	values := []string{}

	if len(args) > 0 {
		for _, card := range cards {
			for key, value := range card.Properties {
				if key == args[0] {
					if !tools.ContainsString(values, tools.ToJSON(value)) {
						values = append(values, tools.ToJSON(value))
					}
				}
			}
		}
	} else {
		for _, card := range cards {
			for key, _ := range card.Properties {
				if !tools.ContainsString(values, key) {
					values = append(values, key)
				}
			}
		}
	}
	for _, value := range values {
		fmt.Println(value)
	}
}
