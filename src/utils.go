package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"sort"
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
	m := getKeyValuePairs(s)
	nm := make(map[string]string)

	for k, v := range m {
		if v.Type().Kind() == reflect.String {
			nm[k] = fmt.Sprintf("\"%v\"", v)
		} else {
			nm[k] = fmt.Sprintf("%v", v)
		}
	}

	return nm
}

func getKeyValuePairs(s interface{}) map[string]reflect.Value {
	m := make(map[string]reflect.Value)
	rt := reflect.TypeOf(settings)
	rv := reflect.ValueOf(settings)

	for fi := 0; fi < rt.NumField(); fi++ {
		m[rt.Field(fi).Name] = rv.Field(fi)
	}

	return m
}

func prettyPrintStruct(s interface{}) {
	kvp := getKeyValuePairs(s)
	for k, v := range kvp {
		val := ""
		switch v.Type().Kind() {
		// Strings
		case reflect.String:
			val = color.GreenString(fmt.Sprintf("\"%v\"", v))

		// Numbers
		case reflect.Int, reflect.Float64, reflect.Float32, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			val = color.GreenString(fmt.Sprintf("\"%v\"", v))

		// Bools, etc
		default:
			val = color.CyanString(fmt.Sprintf("%v", v))
		}
		fmt.Printf("\t%s%s %s\n", color.RedString(k), ":", val)
	}
}

// Pretty-Prints a string[string] map n alphabetical order, colorized
func prettyPrintMap(s map[string]string) {
	keysAlpha := make([]string, len(s))
	i := 0
	// Collate keys and sort them
	for k, _ := range s {
		keysAlpha[i] = k
		i++
	}
	sort.Strings(keysAlpha)

	for _, k := range keysAlpha {
		vStr := fmt.Sprintf("%v", s[k])
		if vStr == "" {
			vStr = `""`
		}

		fmt.Printf(
			"\t%s%s %s\n",
			color.RedString(fmt.Sprintf("%v", k)),
			":",
			color.GreenString(vStr),
		)
	}
}
