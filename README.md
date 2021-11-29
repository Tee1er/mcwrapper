# McWrapper
Interface/wrapper for the Bedrock Dedicated Server written in Go.

## Usage
### Installation
This program requires [Go](https://go.dev/) to run. The latest version is preferred, but it may work with older Go versions (McWrapper has only been tested on Go 1.17).
McWrapper runs on Windows and most non-arm, Debian-based Linux distributions (Ubuntu, Linux Mint, Kali). 
### Useage
McWrapper serves as a sort of intermediary interface between the minecraft bedrock dedicated server and the user, abstracting some common tasks. Starting McWrapper is easy, on windows, run `start.bat` (from the command line or double click), on linux, `chmod +x ./start.cmd` and run it with `./start.cmd` (it might throw a soft error, but it *will* work on POSIX compliant shells). First of all, to do anything useful, you will need to have the dedicated server installed, luckily, McWrapper has a built-in command for installing the server, `init`. When you run `init` in the McWrapper CLI, it will ask you for a URL, go to the [Bedrock dedicated server download page](https://www.minecraft.net/en-us/download/server/bedrock) and copy the URL of the download button for your OS, paste it into McWrapper's URL prompt and hit enter (this may b removed in the fuure because its quite clumsy). McWrapper will try to download and install the server in the `data/server/` subdirectory of your install, if it fails, check the URL and filesystem permissions. After the server is downloaded, starting up the server is as simple as typing run into McWrapper. McWrapper also provides some other commands listed below

 - `stop`: Stops the server (duh)
 - `exit`: Quits McWrapper (and the server if its still running)
 - `clear`: Clears the McWrapper console
 - `help`: Shows some additional help on various commands
 - `settings`: Show or modify the server configuration file (`data/server/server.properties`)
 - `update`: Update the bedrock server from a different download URL, preserving files
 - `server`: Enters a prompt directly connected to the dedicated server, everything you do in the prompt is sent directly to the server. Type `exit` to exit the prompt and go back to the normal McWrapper prompt.


## Source structure
`filemgr.go`      - Utilities for managing server files. <br/>
`utils.go`        - Misc. utility functions. <br/>
`main.go`         - Command loop and user interaction. <br/>
`mcserver.go`     - Abstraction over the dedicated server child process. <br/>
`serverprops.go`  - Modification of server configuration. <br/>
`webhook.go`      - Basic discord communication with webhooks. <br/>

## Features
- [x] Safe server updating
- [x] Basic configuration editor
- [ ] Whitelist editor
- [ ] Player list, info, etc.
- [x] Raw server cli
- [ ] Auto backups and versioning with git
- [ ] Mapping with [Bedrock-Viz](https://github.com/bedrock-viz/bedrock-viz), [PapyrusCS](https://github.com/papyrus-mc/papyruscs) or a [custom LevelDB reader](https://github.com/syndtr/goleveldb)
- [ ] Server status/info REST API?
- [ ] Better webhook interaction
- [ ] Chat relay?
