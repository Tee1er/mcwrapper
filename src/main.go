package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var serverIO serverWrapper

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
			os.RemoveAll("../server")

			fmt.Println("Extracting new server files.")

			unzip(data, "../server")

			fmt.Println("Moving data back.")

			moveServer("../temp", "../server")
			os.RemoveAll("../temp")

		case "help":
			fmt.Println("\tCOMMANDS:\n")

			fmt.Println("\tinit\tDownloads and extracts the latest server version.")
			fmt.Println("\tupdate\tDownloads, extracts, and updates the server to the latest version. Preserves worlds and some other config files.")
			fmt.Println("\trun\tRuns the server.")
			fmt.Println("\texit\tExits the program.")
		case "run":
			go startServer(&serverIO)
		case "stop":
			stopServer(&serverIO)
		case "exit":
			os.Exit(0)
		default:
			fmt.Println("\tInvalid command. Type 'help' for a list of commands.")
		}
	}
}
