package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/fatih/color"
)

type serverWrapper struct {
	cmd    *exec.Cmd
	stdin  *io.WriteCloser
	stdout *io.ReadCloser
}

func startServer(result *serverWrapper) {
	color.Yellow("Starting server.")
	mcServer := exec.Command("../server/bedrock_server.exe")

	stdout, _ := mcServer.StdoutPipe()
	stdin, _ := mcServer.StdinPipe()

	mcServer.Start()
	color.Blue("Server started.")
	*result = serverWrapper{
		stdin:  &stdin,
		stdout: &stdout,
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

}

func stopServer(serverIO *serverWrapper) {
	color.Red("Stopping server.")
	writer := bufio.NewWriter(*serverIO.stdin)
	writer.WriteString("stop\n")
}
