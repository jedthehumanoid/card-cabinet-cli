package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	dir = filepath.Clean(dir) + "/"

	fmt.Println(dir)
	fmt.Println()
	cards := ReadCardDir(dir)

	for i, card := range cards {
		fmt.Printf("%d) %s", i+1, prettyTitle(card.Title))
		if card.Properties["labels"] != nil {
			fmt.Print("\033[1;30m")
			for _, label := range card.Properties["labels"].([]interface{}) {
				fmt.Printf(" [%s]", label.(string))
			}
		}
		fmt.Println("\033[0m")
		//fmt.Printf("%+v\n", card)
	}
}

func prettyTitle(s string) string {
	s = strings.ToUpper(s[:1]) + s[1:]
	s = strings.Replace(s, "-", " ", -1)
	return strings.TrimSuffix(s, ".md")
}
