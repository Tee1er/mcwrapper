package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/fatih/color"
)

type Settings struct {
	WebhookUrl     string `json:"webhook_url`
	WebhookEnabled bool
}

func (s *Settings) Load(fname string) {
	file, err := ioutil.ReadFile(fname)
	check(err)

	json.Unmarshal(file, s)
	s.WebhookEnabled = s.WebhookUrl != ""
}

var (
	serverIO serverWrapper
	wh       = Webhook{}
	settings Settings

	exepath_str, _ = os.Executable()
	basedir        = path.Join(path.Base(exepath_str), "../")
	data_dir       = path.Join(basedir, "/data/")
	settings_path  = path.Join(data_dir, "/settings.json")
	tmpdir, _      = os.MkdirTemp("", "mcwrapper_tmp")
)

const def_settings = "{\n\t\"webhook_url\": \"\"\n}"

func main() {
	c := color.New(color.FgCyan, color.Bold)
	c.Println("-- MCWrapper v0.1-alpha CLI -- \n")

	os.MkdirAll(data_dir, 0664)
	os.MkdirAll(dataPath("/server"), 0664)

	if _, err := os.Stat(settings_path); os.IsNotExist(err) {
		os.WriteFile(settings_path, []byte(def_settings), 0664)
	}

	settings.Load(settings_path)

	// Connect the webhook
	if settings.WebhookEnabled {
		wh.Connect(settings.WebhookUrl)
	}

cmdloop:
	for {
		input := getInput("> ")

		switch input {
		case "init":
			url := getUrl("Enter URL for server download: ")
			data := getServer(url)
			getServer(url)
			unzip(data, dataPath("/server"))
		case "update":
			url := getUrl("Enter URL for server download: ")
			data := getServer(url)
			tmpdatadir := tmpPath("data_backup")

			fmt.Println("Moving data to temporary dir.")

			moveServer(dataPath("/server"), tmpdatadir)
			os.RemoveAll(dataPath("/server"))

			fmt.Println("Extracting new server files.")

			unzip(data, dataPath("/server"))

			fmt.Println("Moving data back.")

			moveServer(tmpdatadir, dataPath("/server"))
			os.RemoveAll(tmpdatadir)
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
			break cmdloop
		default:
			fmt.Println("\tInvalid command. Type 'help' for a list of commands.")
		}
	}

	os.RemoveAll(tmpdir)
}
