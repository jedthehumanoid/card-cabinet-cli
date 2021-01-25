package cmd

import (

	"os"
	"path/filepath"

	"github.com/jedthehumanoid/cardcabinet"
)

func main() {
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

	_ = cardcabinet.ReadCards(config.Src)
	_ = cardcabinet.ReadBoards(config.Src)

	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	//case "boards", "b":
//		listBoards(boards, config)
//	case "cards", "c":
	//	for _, card := range cards {
//			fmt.Println(card.Name)
//		}
//	default:
	//	fmt.Printf("Unknown command: %s\n", command)
	}
}
