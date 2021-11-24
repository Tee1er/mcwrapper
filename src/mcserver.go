package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

type ServerWrapper struct {
	cmd    *exec.Cmd
	stdin  *io.WriteCloser
	stdout *io.ReadCloser
}

func startServer(result *ServerWrapper) {
	color.Yellow("Starting server.")

	var mcServer *exec.Cmd
	if runtime.GOOS == "windows" {
		mcServer = exec.Command(dataPath("/server/bedrock_server.exe"))
	} else {
		mcServer = exec.Command(
			fmt.Sprintf("LD_LIBRARY_PATH=%s %s", dataPath("/server/bedrock_server"), dataPath("/server/bedrock_server")),
		)
	}

	stdout, _ := mcServer.StdoutPipe()
	stdin, _ := mcServer.StdinPipe()

	mcServer.Start()
	color.Blue("Server started.")
	*result = ServerWrapper{
		cmd:    mcServer,
		stdin:  &stdin,
		stdout: &stdout,
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

}

func stopServer(serverIO *ServerWrapper) {
	color.Red("Stopping server.")
	writer := bufio.NewWriter(*serverIO.stdin)
	writer.WriteString("stop\n")
}
