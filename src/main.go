package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/mholt/archiver"
	"github.com/otiai10/copy"
)

//TODO: move to settings.json file.
var copyFiles = []string{"server.properties", "permissions.json", "whitelist.json"}
var copyDirs = []string{"worlds"}

func check(err error) {
	if err != nil {
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

func getUrl(prompt string) string {
	result := getInput(prompt)
	_, err := url.ParseRequestURI(result)
	check(err)
	color.Green("Valid URL.")
	return result
}

func main() {
	c := color.New(color.FgCyan, color.Bold)
	c.Println("-- MCWrapper v0.1-alpha CLI -- \n")
	for {
		input := getInput("> ")

		switch input {
			case "init":
				url := getUrl("Enter URL for server download: ")
				data := getServer(url)
				getServer(url)
				unzip(data, "../server")
			case "update":
				url := getUrl("Enter URL for server download: ")
				data := getServer(url)
				fmt.Println("Moving data to temporary dir.")
				moveServer("../server", "../temp")
				err := os.RemoveAll("../server")
				check(err)
				fmt.Println("Extracting new server files.")
				unzip(data, "../server")
				fmt.Println("Moving data back.")
				moveServer("../temp", "../server")


			case "help":
				fmt.Println("\tCOMMANDS:\n")
				fmt.Println("\tinit\tDownloads and extracts the latest server version.")
				fmt.Println("\tupdate\tDownloads, extracts, and updates the server to the latest version.")
				fmt.Println("\texit\tExits the program.")
			case "exit":
				os.Exit(0)
			default:
				fmt.Println("\tInvalid command. Type 'help' for a list of commands.")
			}
	}
}

func getServer(url string) []byte {
	//"https://minecraft.azureedge.net/bin-win/bedrock-server-1.17.41.01.zip"
	// Latest ZIP file for server release ... unsure how to get the latest version since version numbers are not always in order. Could ask user for version #.
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

func unzip(data []byte, path string) {
	os.WriteFile("../server.zip", data, 0644)

	color.Green("Downloaded successfully! \nExtracting ZIP archive. ")

	archiver.Unarchive("../server.zip", path)

	os.Remove("../server.zip")

}

func moveServer(serverPath string, destPath string) {
	// TODO - this is a bit of a hack. Need to find a better way to do this.
	// Copy all files from temp to server.
	files, err := ioutil.ReadDir(serverPath)
	check(err)

	for _, file := range files {
		if !file.IsDir() && contains(copyFiles, file.Name())  {
			err := os.Rename(serverPath + "/" + file.Name(), destPath + "/" + file.Name())
			check(err)
			fmt.Printf("Copied file %s to %s\n", file.Name(), destPath + "/" + file.Name() )
		} else if file.IsDir() && contains(copyDirs, file.Name()) {
			err := copy.Copy(serverPath + "/" + file.Name(), destPath + "/" + file.Name())
			check(err)
			fmt.Printf("Copied directory %s to %s\n", file.Name(), destPath + "/" + file.Name())

		}
	}
}
