package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

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
	mcServer        ServerWrapper
	wh              = Webhook{}
	settings        Settings
	defaultSettings = Settings{
		WebhookUrl: "",
	}

	exepathStr, _ = os.Executable()
	baseDir       = path.Join(path.Base(exepathStr), "../")
	dataDir       = path.Join(baseDir, "/data/")
	settingsPath  = path.Join(dataDir, "/settings.json")
	tmpDir, _     = os.MkdirTemp("", "mcwrapper_tmp")
)

func handleSignal(sig os.Signal) {
	os.RemoveAll(tmpDir) // Cleanup, no temp file left behind
	mcServer.Stop()
}

func handleCommand(cmdHistory []string, input string) bool {
	cmdHistory = append(cmdHistory, input)
	inputs := strings.Split(input, " ")
	args := inputs[1:]

	switch inputs[0] {
	case "init":
		url := getUrl("Enter URL for server download: ")
		data, err := getServer(url)
		if err != nil {
			printErr(err)
			return false
		}
		getServer(url)
		unzip(data, dataPath("/server"))

	case "update":
		url := getUrl("Enter URL for server download: ")
		data, err := getServer(url)
		if err != nil {
			printErr(err)
			return false
		}
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
		fmt.Println("\tclear\tClears the console.")
		fmt.Println("\tserver\tEnters into another sub-console, with input passed directly to the dedicated server.")
		fmt.Println("\tsettings\tPrints the currently loaded settings (from data/settings.json).")
		fmt.Println("\texit\tExits the program.")

	case "server":
		// FIXME:

		mcServer.pipeStdout = true
		mcServer.pipeStdin = true
		if mcServer.overflow.Len() != 0 {
			fmt.Print(mcServer.overflow.String())
			mcServer.overflow.Reset()
		}

		for {
			rawInp := getInput("")
			if strings.TrimSpace(rawInp) == "exit" {
				break
			}

			mcServer.Send(rawInp)
		}

		// Restore pipework
		mcServer.pipeStdout = false
		mcServer.pipeStdin = false

	case "run":
		err := mcServer.Start()
		if err != nil {
			color.Yellow(err.Error())
		}

	case "clear":
		fmt.Print("\033[H\033[2J") // Should work

	case "settings":
		props, err := parseProperties(dataPath("/server/server.properties"))
		if err != nil {
			printErr(err)
			return false
		}

		switch len(args) {
		case 0:
			fmt.Println("\nMcWrapper settings:")
			prettyPrintStruct(settings)

			fmt.Println("\nServer properties:")
			prettyPrintMap(props)

		case 1:
			if val, contains := props[args[0]]; contains {
				fmt.Printf("\nServer property '%s' is '%s'\n", args[0], val)
			} else {
				color.HiRed("Error, property '%s' not found in server.properties", args[0])
			}

		case 2:
			err = setProperty(dataPath("/server/server.properties"), args[0], args[1])
			if err != nil {
				printErr(err)
			}

		default:
			color.HiRed("Error, invalid arguments for 'settings'")
		}

	case "stop":
		mcServer.Stop()

	case "exit":
		mcServer.Stop()
		return true

	default:
		fmt.Println("\tInvalid command. Type 'help' for a list of commands.")
	}
	return false
}

func main() {
	os.MkdirAll(dataDir, 0664)
	os.MkdirAll(dataPath("/server"), 0664)

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		// MarshalIndent for pretty-print
		defSettingsStr, _ := json.MarshalIndent(defaultSettings, "", "\t")
		os.WriteFile(settingsPath, []byte(defSettingsStr), 0664)
	}

	settings.Load(settingsPath)

	// Connect the webhook
	if settings.WebhookEnabled {
		wh.Connect(settings.WebhookUrl)
	}

	// Handle signals (Ctrl-C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChan // So the variable is unused, sigh
		handleSignal(sig)
		os.Exit(0)
	}()

	cmdHistory := make([]string, 20)
	for {
		if handleCommand(cmdHistory, getInput("> ")) {
			break
		}
	}

	os.RemoveAll(tmpDir)
}
