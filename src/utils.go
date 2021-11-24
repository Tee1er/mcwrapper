package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

func check(err error) {
	if err != nil {
		// Do cleanup
		os.RemoveAll(tmpDir)
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
	return path.Join(dataDir, p)
}

func tmpPath(p string) string {
	resolved := path.Join(tmpDir, p)
	os.MkdirAll(resolved, 0666)
	return resolved
}

/// Returns a string:string map of struct field names and stringified values
/// Strings are surrounded in quotation marks
func getStrKeyValues(s interface{}) map[string]string {
	m := make(map[string]string)
	rt := reflect.TypeOf(settings)
	rv := reflect.ValueOf(settings)

	for fi := 0; fi < rt.NumField(); fi++ {
		ft := rt.Field(fi)
		fv := rv.FieldByName(ft.Name)
		if ft.Type.Kind() == reflect.String {
			m[ft.Name] = fmt.Sprintf("\"%v\"", fv)
		} else {
			m[ft.Name] = fmt.Sprintf("%v", fv)
		}
	}

	return m
}
