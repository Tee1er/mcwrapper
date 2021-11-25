package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/fatih/color"
)

type ServerWrapper struct {
	serverpath string
	cmd        *exec.Cmd
	stdin      *bufio.Writer
	stdout     *bufio.Reader
}

func (sw *ServerWrapper) Start() {
	color.Yellow("Starting server.")

	sw.serverpath = dataPath("/server")
	if runtime.GOOS == "windows" {
		sw.cmd = exec.Command(path.Join(sw.serverpath, "/bedrock_server.exe"))
	} else if runtime.GOOS == "linux" {
		fmt.Printf("%s\n", path.Join(sw.serverpath, "/bedrock_server"))
		sw.cmd = exec.Command("./bedrock_server")
		sw.cmd.Env = append(make([]string, 0), "LD_LIBRARY_PATH=.")
	} else {
		color.Red("Error, your system doesn't appear to support the minecraft bedrock dedicated server.")
	}

	// Set CWD to inside server data dir
	sw.cmd.Dir = sw.serverpath

	// Get pipes and make r/w IO for them
	stdout, _ := sw.cmd.StdoutPipe()
	stdin, _ := sw.cmd.StdinPipe()
	sw.stdout = bufio.NewReader(stdout)
	sw.stdin = bufio.NewWriter(stdin)

	// Hook server stdout to go's stdout
	sw.cmd.Stdout = os.Stdout

	sw.cmd.Start()
	color.Blue("Server started.")
}

func (sw *ServerWrapper) Stop() {
	color.Red("Stopping server.")
	if sw.stdin != nil {
		sw.Send("stop")
		//sw.cmd.Process.Kill()
	}
}

func (sw *ServerWrapper) Send(cmd string) {
	sw.stdin.WriteString("stop\n")
	sw.stdin.Flush()
}
