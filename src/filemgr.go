package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/mholt/archiver"
	"github.com/otiai10/copy"
)

//TODO: move to settings.json file.
var copyFiles = []string{"server.properties", "permissions.json", "whitelist.json"}
var copyDirs = []string{"worlds"}

// Credit to var23rav: https://gist.github.com/var23rav/23ae5d0d4d830aff886c3c970b8f6c6b
func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func getServer(url string) []byte {
	//"https://minecraft.azureedge.net/bin-win/bedrock-server-1.17.41.01.zip"
	// Latest ZIP file for server release ... unsure how to get the latest version since version numbers are not always in order. Could ask user for version #.
	// Alternatively, could use GH to get the latest version
	resp, err := http.Get(url)

	check(err)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		panic("HTTP error " + string(resp.StatusCode))
	}

	// Store the response in a file, so that archiver can unzip it.
	// TODO - it would be better if we could just pass the response and not have to store it.

	bytes, err := ioutil.ReadAll(resp.Body)

	check(err)

	resp.Body.Close()

	return bytes
}

func unzip(data []byte, spath string) {
	os.WriteFile(dataPath("/server.zip"), data, 0644)

	color.Green("Downloaded successfully! \nExtracting ZIP archive. ")

	archiver.Unarchive(dataPath("/server.zip"), spath)

	os.Remove(dataPath("/server.zip"))
}

func moveServer(serverPath string, destPath string) {
	// TODO - this is a bit of a hack. Need to find a better way to do this.
	// Copy all files from temp to server.
	files, err := ioutil.ReadDir(serverPath)
	check(err)

	// Copy all files from source to destination. TODO: clean this up.

	for _, file := range files {
		if !file.IsDir() && contains(copyFiles, file.Name()) {
			err := moveFile(serverPath+"/"+file.Name(), destPath+"/"+file.Name())
			check(err)
			fmt.Printf("Copied file %s to %s\n", file.Name(), destPath+"/"+file.Name())
		} else if file.IsDir() && contains(copyDirs, file.Name()) {
			err := copy.Copy(serverPath+"/"+file.Name(), destPath+"/"+file.Name())
			check(err)
			fmt.Printf("Copied directory %s to %s\n", file.Name(), destPath+"/"+file.Name())

		}
	}
}
