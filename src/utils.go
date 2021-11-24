package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func check(err error) {
	if err != nil {
		// Do cleanup
		os.RemoveAll(tmpdir)
		panic(err)
	}
}

func contains(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.Replace(strings.Replace(input, "\n", "", -1), "\r", "", -1)
}

var matchCtrlChars = regexp.MustCompile("\x1B(?:[@-Z\\-_]|\\[[0-?]*[ -/]*[@-~])")

func cleanCtrlChars(data string) string {
	return matchCtrlChars.ReplaceAllString(data, "")
}

func getUrl(prompt string) string {
	result := getInput(prompt)
	_, err := url.ParseRequestURI(result)
	check(err)
	color.Green("Valid URL.")
	return result
}

func dataPath(p string) string {
	return path.Join(data_dir, p)
}

func tmpPath(p string) string {
	resolved := path.Join(tmpdir, p)
	os.MkdirAll(resolved, 0666)
	return resolved
}
