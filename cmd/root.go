package cmd

import (
	"card-cabinet-cli/config"
	"fmt"
	"io/ioutil"
	"os"
	

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

var cfg config.Config
var recursive bool
var query string

var rootCmd = &cobra.Command{
	Use:   "card-cabinet-cli",
	Short: "",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func init() {
	//rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "print debug output")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", true, "Recurse into subdirectories")
	cobra.OnInitialize(readConfig)
}

func readConfig() {
	d, _ := ioutil.ReadFile("cabinet.toml")
	_, _ = toml.Decode(string(d), &cfg)
}
