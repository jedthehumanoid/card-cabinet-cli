package view

import (
	"card-cabinet-cli/slug"
	"fmt"
	"strings"

	"card-cabinet-cli/ansi"
	"card-cabinet-cli/config"
	"github.com/jedthehumanoid/cardcabinet"
)

const yellow = "\033[38;2;250;189;47m"
const gray = "\033[38;2;100;100;100m"
const darkgray = "\033[38;2;50;50;50m"

func List(card cardcabinet.Card, config config.Config) {
	title := slug.From(strings.TrimPrefix(card.Name, card.Path()))
	title = strings.ToUpper(title[:1]) + title[1:]
	title = strings.TrimSuffix(title, ".md")

	if card.Path() != "" {
		fmt.Printf("%s%s%s ", gray, card.Path(), ansi.Reset)
	}
	fmt.Printf("%s", title)

	if card.Contents != "" {
		fmt.Printf("%s%s%s", yellow, " \u2261", ansi.Reset)
	}

	for _, label := range asStringSlice(card.Properties["labels"]) {
		c, hascolor := config.Colors[label]
		if hascolor {
			fmt.Printf(ansi.Color(c))
		} else {
			fmt.Printf(gray)
		}
		fmt.Printf(" #%s", label)
	}
	fmt.Println(ansi.Reset)
}

func Cat(card cardcabinet.Card, config config.Config) {
	columns := ansi.GetColumns()

	fmt.Println(darkgray + "\u250c" + fill("\u2500", columns-2) + ansi.Reset)
	fmt.Println(darkgray + "\u2502 " + yellow + card.Name + ansi.Reset)
	fmt.Println(darkgray + "\u251c" + fill("\u2500", columns-2) + ansi.Reset)
	if card.MarshalFrontmatter(false) != "" {
		for _, line := range strings.Split(card.MarshalFrontmatter(false), "\n") {
			fmt.Println(darkgray + "\u2502 " + gray + line + ansi.Reset)
		}
	}
	if card.Contents != "" {
		for _, line := range strings.Split(card.Contents, "\n") {
			for _, line := range splitlen(line, columns-2) {
				fmt.Println(darkgray + "\u2502 " + ansi.Reset + line)
			}
		}
	}
	fmt.Println(darkgray + "\u2514" + fill("\u2500", columns-2) + ansi.Reset)

}

func asStringSlice(i interface{}) []string {
	ret := []string{}
	if i == nil {
		return ret
	}
	for _, v := range i.([]interface{}) {
		ret = append(ret, v.(string))
	}
	return ret
}

// Fill returns a len long string filled with char
func fill(char string, len int) string {
	ret := ""
	for i := 0; i < len; i++ {
		ret += char
	}
	return ret
}

// Splitlen splits string by length
func splitlen(s string, length int) []string {
	ret := []string{}
	for len(s) > length {
		ret = append(ret, s[:length])
		s = s[length:]
	}

	ret = append(ret, s)
	return ret
}
