# McWrapper
Interface/wrapper for the Bedrock Dedicated Server written in Go.

## Usage
### Installation
This program requires [Go](https://go.dev/) to run. The latest version is preferred, but it may work with older Go versions (McWrapper has only been tested on Go 1.17). In the future we may try distributing a compiled binary for MCWrapper; check the Releases sidebar. 
McWrapper runs on Windows and most non-arm, Debian-based Linux distributions (Ubuntu, Linux Mint, Kali). However, the Bedrock server *itself* is only compiled for Windows and Ubuntu.
### Useage
McWrapper serves as an intermediary between the MCBE Bedrock Dedicated Server and the user, abstracting and automating some common tasks. It uses a CLI (command line interface), which can be intimidating but is really quite simple. Starting McWrapper is easy; on Windows, run `start.bat` (from the command line or by double- clicking); on Linux, `chmod +x ./start.cmd` and run it with `./start.cmd` (it might throw a soft error, but it *will* work on POSIX compliant shells). 

To be able to do anything useful, you will need to have the dedicated server software installed. Luckily, McWrapper has a built-in command for installing the server, `init`. When you run `init` in the McWrapper CLI, it will ask you for a URL, go to the [Bedrock Dedicated Server download page](https://www.minecraft.net/en-us/download/server/bedrock) and copy the URL of the download button for your OS, paste it into McWrapper's URL prompt and hit enter (This may be changed in the future since it *is* quite clumsy. ). McWrapper will try to download and install the server in the `data/server/` subdirectory of your install, if it fails, check the URL and filesystem permissions. After the server is downloaded, starting up the server is as simple as typing run into McWrapper. McWrapper also provides some other commands listed below

 - `stop`: Stops the server
 - `exit`: Quits McWrapper (and the server if it's still running)
 - `clear`: Clears the McWrapper console
 - `help`: Shows some additional help on various commands
 - `settings`: Show or modify the server configuration file (`data/server/server.properties`)
 - `update`: Update the bedrock server from a different download URL, preserving files


## Source structure
`filemgr.go`      - Utilities for managing server files. <br/>
`utils.go`        - Misc. utility functions. <br/>
`main.go`         - Command loop and user interaction. <br/>
`mcserver.go`     - Abstraction over the dedicated server child process. <br/>
`serverprops.go`  - Modification of server configuration. <br/>
`webhook.go`      - Basic discord communication with webhooks. <br/>

## Features
- [x] Safe server updating & installation.
- [x] Basic configuration editor
- [ ] Whitelist editor
- [ ] Player list, info, etc.
- [ ] Raw server cli
- [ ] Auto backups and versioning with git
- [ ] Mapping with [Bedrock-Viz](https://github.com/bedrock-viz/bedrock-viz), [PapyrusCS](https://github.com/papyrus-mc/papyruscs) or a [custom LevelDB reader](https://github.com/syndtr/goleveldb)
- [ ] Server status/info REST API?
- [ ] Better webhook interaction
- [ ] Chat relay?
