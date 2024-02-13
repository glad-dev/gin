package main

import (
	"github.com/glad-dev/gin/cmd"
	"github.com/glad-dev/gin/log"
	"github.com/glad-dev/gin/style"
)

func main() {
	log.Init()

	err := cmd.Execute()
	if err != nil {
		style.PrintErrAndExit("Execute failed: " + err.Error())
	}
}
