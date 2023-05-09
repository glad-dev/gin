package main

import (
	"gn/cmd"
	"gn/logger"
	"gn/style"
)

func main() {
	logger.Init()
	style.Init()

	err := cmd.Execute()
	if err != nil {
		style.PrintErrAndExit("Execute failed: " + err.Error())
	}
}
