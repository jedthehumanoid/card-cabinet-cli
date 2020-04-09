package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Card is a text with properties, meant for displaying on a board.
type Card struct {
	Title      string                 `json:"title"`
	Contents   string                 `json:"contents"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// ReadCardFile takes a file path, reading file in to a card.
func ReadCardFile(path string) (Card, error) {
	var card Card

	card.Title = ToSlug(strings.TrimSuffix(path, ".md")) + ".md"

	contents, err := ioutil.ReadFile(filepath.FromSlash(path))
	if err != nil {
		return card, err
	}
	frontmatter := string(GetYAMLFrontmatter(contents))
	card.Contents = strings.TrimPrefix(string(contents), frontmatter)
	card.Contents = strings.TrimSpace(card.Contents)
	frontmatter = strings.TrimSpace(frontmatter)
	frontmatter = strings.Trim(frontmatter, "---")

	err = yaml.Unmarshal([]byte(frontmatter), &card.Properties)
	return card, err
}

func MarshalFrontmatter(card Card) string {
	b, _ := yaml.Marshal(card.Properties)
	ret := ""
	frontmatter := strings.TrimSpace(string(b))
	if frontmatter != "{}" {
		ret = "---\n" + frontmatter + "\n---"
	}
	return ret
}

// WriteCardFile writes a card to file
func WriteCardFile(path string, card Card) error {
	y, err := yaml.Marshal(card.Properties)
	if err != nil {
		panic(err)
	}

	filedata := "---\n" + string(y) + "---\n"
	filedata += card.Contents

	err = ioutil.WriteFile(path, []byte(filedata), 0644)

	if err != nil {
		panic(err)
	}

	return nil
}

func IsCard(file string) bool {
	return strings.HasSuffix(file, ".md")
}

func ReadCardDir(dir string) []Card {
	cards := []Card{}

	for _, file := range FindFiles(dir) {
		if !IsCard(file) {
			continue
		}
		card, err := ReadCardFile(file)
		if err != nil {
			panic(err)
		}
		card.Title = strings.TrimPrefix(card.Title, dir)
		cards = append(cards, card)
	}

	return cards
}
