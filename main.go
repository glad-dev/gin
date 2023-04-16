package main

import (
	"fmt"
	"os"

	"gn/cmd"
	"gn/logger"
	"gn/style"
)

func main() {
	logger.Init()
	style.Init()

	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Execute failed: %s", err)
		os.Exit(1)
	}
}
