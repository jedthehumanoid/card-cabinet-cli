package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"card-cabinet-cli/tools"

)

type Config struct {
	Src    string            `toml:"src"`
	Colors map[string]string `toml:"colors"`
}

var config Config

const gray = "\033[38;2;100;100;100m"
const darkgray = "\033[38;2;50;50;50m"
const yellow = "\033[38;2;250;189;47m"
const reset = "\033[0m"

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
	cobra.OnInitialize(readConfig)
}

func readConfig() {
	tools.LoadToml("cabinet.toml", &config)
}
