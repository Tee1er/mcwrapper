package main

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// Not-So-Nicely display an error
func check(err error) {
	if err != nil {
		// Do cleanup
		os.RemoveAll(tmpDir)
		panic(err)
	}
}

// Nicely display an error
func printErr(err error) {
	color.HiRed("Error: %s", err.Error())
}

func contains(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

// Global stdin reader instance, ensured to be non-nil
var inputReader = bufio.NewReader(os.Stdin)

// A tee-ed copy of inputReader, that can be assigned to the server's stdin
var miscReader = io.TeeReader(inputReader, os.Stdin)

func getInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := inputReader.ReadString('\n')
	return strings.Trim(input, "\n\r")
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
	for k := range s {
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

func relayIf(in io.Reader, out io.Writer, cond *bool, exit chan bool, alternate io.Writer) {
	reader := bufio.NewReader(in)
	for {
		if len(exit) != 0 {
			return
		}

		data, noeol, _ := reader.ReadLine()
		if *cond {
			out.Write(data)
			if !noeol {
				out.Write([]byte("\n"))
			}
		} else if alternate != nil {
			alternate.Write(data)
			if !noeol {
				alternate.Write([]byte("\n"))
			}
		}
	}
}

func relayWhile(in io.Reader, out io.Writer, cond *bool) {
	reader := bufio.NewReader(in)
	for {
		data, noeol, _ := reader.ReadLine()
		if *cond {
			out.Write(data)
			if !noeol {
				out.Write([]byte("\n"))
			}
		} else {
			return
		}
	}
}
