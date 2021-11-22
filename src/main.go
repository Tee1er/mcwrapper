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
)

func check(err error) {
	if err != nil {
		panic(err)
	}
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
		if input == "exit" {
			color.Red("Exiting ...")
			os.Exit(0)
		}
		// if input == "update" {
		// 	url := getUrl("Enter URL for server download: ")

		// }
		if input == "init" {
			url := getUrl("Enter URL for server download: ")
			download(url)
		}
		if input == "help" {
			fmt.Println("\tCOMMANDS:\n")
			fmt.Println("\tinit\tDownloads and extracts the latest server version.")
			fmt.Println("\tupdate\tDownloads, extracts, and updates the server to the latest version.")
			fmt.Println("\texit\tExits the program.")
		}
	}
}

func download(url string) {
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
	
		os.WriteFile("../server.zip", bytes, 0644)
	
		color.Green("Downloaded successfully! \nExtracting ZIP archive. ")
	
		archiver.Unarchive("../server.zip", "./server")
	
		os.Remove("../server.zip")
	
		resp.Body.Close()
}