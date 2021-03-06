package tools

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ToJSON return JSON representation of interface
func ToJSON(in interface{}) string {
	b, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// ToJSONPretty return indented JSON representation of interface
func ToJSONPretty(in interface{}) string {
	b, err := json.MarshalIndent(in, "", "   ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func GetColumns() int {

	cmd := exec.Command("tput", "cols")
	columns, _ := cmd.Output()
	ret, _ := strconv.Atoi(strings.TrimSpace(string(columns)))
	return ret
}

func Splitlen(s string, length int) []string {
	ret := []string{}
	for len(s) > length {
		ret = append(ret, s[:length])
		s = s[length:]
	}

	ret = append(ret, s)
	return ret
}

func Fill(c string, len int) string {
	ret := ""
	for i := 0; i < len; i++ {
		ret += c
	}
	return ret
}

func AsStringSlice(i interface{}) []string {
	ret := []string{}
	if i == nil {
		return ret
	}
	for _, v := range i.([]interface{}) {
		ret = append(ret, v.(string))
	}
	return ret
}

func extractPrefix(list []string, prefix string) ([]string, []string) {
	normal := []string{}
	prefixed := []string{}

	for _, arg := range list {
		if strings.HasPrefix(arg, prefix) {
			prefixed = append(prefixed, arg)
		} else {
			normal = append(normal, arg)
		}
	}
	return normal, prefixed
}

func LoadToml(file string, i interface{}) error {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(d), i)
	return err
}

// ContainsString searches slice for string
func ContainsString(list []string, s string) bool {
	for _, val := range list {
		if val == s {
			return true
		}
	}
	return false
}

// FromSlug returns "this format" from "this-format"
func FromSlug(s string) string {
	return strings.Replace(s, "-", " ", -1)
}

// ToSlug returns a "page slug" out of string
// Makes lowercase, removes international accents, turns spaces and
// other nonalphanumeric characters to dashes, but keeps slashes
func ToSlug(s string) string {
	var re = regexp.MustCompile("[^a-z0-9/]+")

	s = strings.ToLower(s)
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)
	s, _, _ = transform.String(t, s)
	s = re.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func ToMap(i interface{}) map[string]interface{} {
	var ret map[string]interface{}

	js, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(js, &ret)
	if err != nil {
		panic(err)
	}
	return ret
}
