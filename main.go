package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jedthehumanoid/cardcabinet"
)

const defaultcommand = "list"
const gray = "\033[38;2;100;100;100m"
const darkgray = "\033[38;2;50;50;50m"
const yellow = "\033[38;2;250;189;47m"
const reset = "\033[0m"

type Config struct {
	Src    string            `toml:"src"`
	Colors map[string]string `toml:"colors"`
}

func loadConfig(file string) Config {
	var config Config
	loadToml(file, &config)
	return config
}

func main() {
	config := loadConfig("cabinet.toml")

	if config.Src == "" {
		config.Src = "."
	}
	config.Src = filepath.Clean(config.Src) + "/"
	/*
		piped := false
		fi, _ := os.Stdout.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			piped = true
		}
	*/

	cards := cardcabinet.ReadCards(config.Src)
	boards := cardcabinet.ReadBoards(config.Src)

	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "boards", "b":
		listBoards(boards, config)
	case "cards", "c":
		for _, card := range cards {
			fmt.Println(card.Name)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}
