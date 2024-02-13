package main

import (
	"github.com/glad-dev/gin/cmd"
	"github.com/glad-dev/gin/logger"
	"github.com/glad-dev/gin/style"
)

func main() {
	logger.Init()

	err := cmd.Execute()
	if err != nil {
		style.PrintErrAndExit("Execute failed: " + err.Error())
	}
}
