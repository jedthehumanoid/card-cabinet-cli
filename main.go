package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultcommand = "list"
const gray = "\033[38;2;127;127;127m"
const reset = "\033[0m"

type Config struct {
	Src string `toml:"src"`
}

func getArguments() (string, string, []string, string, []string) {
	filter := ""
	command := ""
	arguments := []string{}
	flags := []string{}
	selected := ""

	// Extract flags
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			arguments = append(arguments, arg)
		}
	}

	// Extract filter
	temp := []string{}
	for _, arg := range arguments {
		if strings.HasPrefix(arg, "@") ||
			strings.HasPrefix(arg, "+") ||
			strings.HasPrefix(arg, "/") {
			filter = arg
		} else {
			temp = append(temp, arg)
		}
	}
	arguments = temp

	// Extract selected
	temp = []string{}
	for _, arg := range arguments {
		_, err := strconv.ParseInt(arg, 10, 32)
		if err == nil {
			selected = arg
		} else {
			temp = append(temp, arg)
		}
	}
	arguments = temp

	if len(arguments) == 0 {
		command = defaultcommand
	} else {
		command = arguments[0]
		arguments = arguments[1:]
	}

	return filter, command, arguments, selected, flags
}

func loadConfig(file string) Config {
	var config Config

	d, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("no configuration file")
		return config
	}

	_, err = toml.Decode(string(d), &config)

	if err != nil {
		panic(err)
	}

	return config
}

func getLabels(cards []Card) []string {
	labels := []string{}
	for _, card := range cards {
		for _, label := range asStringSlice(card.Properties["labels"]) {
			if !ContainsString(labels, "+"+label) {
				labels = append(labels, "+"+label)
			}
		}
	}
	return labels
}

func filterLabel(cards []Card, label string) []Card {
	ret := []Card{}
	for _, card := range cards {
		labels := asStringSlice(card.Properties["labels"])

		if ContainsString(labels, label) {
			ret = append(ret, card)
		}
	}
	return ret
}

func main() {

	config := loadConfig("cabinet.toml")
	filter, command, _, selected, _ := getArguments()

	dir := "."

	if config.Src != "" {
		dir = config.Src
	}

	dir = filepath.Clean(dir) + "/"

	cards := ReadCardDir(dir)
	boards := ReadBoards(dir)
	boards = append(boards, getLabels(cards)...)

	if filter != "" {
		switch filter[0] {
		case '+':
			cards = filterLabel(cards, filter[1:])
		}
	}

	if selected != "" {
		i, err := strconv.Atoi(selected)
		if err == nil && i <= len(cards) {
			cards = cards[i-1 : i]
		}
	}

	switch command {
	case "boards":
		for _, board := range boards {
			fmt.Println(board)
		}

	case "list":
		for i, card := range cards {
			fmt.Printf("%d) %s", i+1, prettyTitle(card.Title))
			if card.Contents != "" {
				fmt.Print(" []")
			}
			fmt.Print(gray)
			for _, label := range asStringSlice(card.Properties["labels"]) {
				fmt.Printf(" [%s]", label)
			}
			fmt.Println(reset)
		}
	case "cat", "c":
		for _, card := range cards {
			fmt.Println("\n\n\n" + card.Title)
			fmt.Println(MarshalFrontmatter(card))
			fmt.Println(card.Contents)
		}
	case "filename", "f":
		for _, card := range cards {
			fmt.Printf("%s%s\n", dir, card.Title)
		}
	}
}

func prettyTitle(s string) string {
	s = strings.ToUpper(s[:1]) + s[1:]
	s = strings.Replace(s, "-", " ", -1)
	return strings.TrimSuffix(s, ".md")
}
