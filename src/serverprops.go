package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

//var valuesRegex = regexp.MustCompile(`(?gim)^(?!( |\t)*(#.*)).*$`)

func parseProperties(filepath string) (map[string]string, error) {
	props := make(map[string]string)

	fileContents, err := ioutil.ReadFile(filepath)
	if os.IsNotExist(err) {
		return nil, errors.New("properties file not found, do you have the dedicated server set up?")
	}

	propertiesStr := string(fileContents)
	propertiesValues := strings.Split(propertiesStr, "\n") //valuesRegex.FindAllString(propertiesStr, -1)

	for _, prop := range propertiesValues {
		if !strings.HasPrefix(strings.Trim(prop, " \t"), "#") && prop != "" {
			split := strings.Split(prop, "=")
			val := ""

			if len(split) >= 2 {
				val = split[1]
			}

			props[split[0]] = val
		}
	}

	return props, nil
}
