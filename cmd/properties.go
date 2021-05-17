package cmd

import (
	"card-cabinet-cli/tools"
	"fmt"
	"github.com/jedthehumanoid/cardcabinet"
	"github.com/spf13/cobra"
	"path/filepath"
)

var expandall bool

func init() {
	rootCmd.AddCommand(propsCmd)
	propsCmd.PersistentFlags().StringVarP(&board, "board", "b", "", "Show properties in board")
	propsCmd.PersistentFlags().StringVarP(&board, "expand", "e", "", "Show keys for property")
	propsCmd.PersistentFlags().BoolVarP(&expandall, "expand-all", "E", true, "Show keys for all properties")
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
	if len(args) > 0 {
		config.Src = args[0]
		config.Src = filepath.Clean(config.Src) + "/"
	}

	cards := cardcabinet.ReadCards(config.Src, recursive)
	properties := map[string][]string{}

		for _, card := range cards {
			for key, value := range card.Properties {
				switch value.(type) {
				case string, int, bool:
					properties[key] = append(properties[key], tools.ToJSON(value))
				case []interface{}:
					for _,v := range value.([]interface{}) {
						properties[key] = append(properties[key], tools.ToJSON(v))
					}
				}		
			}
		}
	
		for key, value := range properties {
			fmt.Println(key)
			if expandall {
			for _, v := range unique(value){
				fmt.Print(gray + v + reset + " ")
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