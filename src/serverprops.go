package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
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
		if !strings.HasPrefix(strings.TrimSpace(prop), "#") && prop != "" {
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

func setProperty(propfile string, propname string, propval string) error {
	fileContents, err := ioutil.ReadFile(propfile)
	if os.IsNotExist(err) {
		return errors.New("properties file not found, do you have the dedicated server set up?")
	}

	propertiesStr := string(fileContents)

	propRegexp := regexp.MustCompile(fmt.Sprintf(`(?m)^%s=.*$`, propname))
	if propRegexp.MatchString(propertiesStr) {
		propertiesStr = propRegexp.ReplaceAllString(propertiesStr, fmt.Sprintf("%s=%s", propname, propval))
		color.Green("Set server property '%s' to '%s'. These changes will be reflected when you next restart the server", propname, propval)
	} else {
		return fmt.Errorf("could not find specified property \"%s\"", propname)
	}

	err = os.WriteFile(propfile, []byte(propertiesStr), 0664)
	if err != nil {
		return fmt.Errorf("error writing modified properties to file '%s'", propfile)
	}

	return nil
}
