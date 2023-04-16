package main

import (
	"fmt"
	"os"

	"gn/cmd"
	"gn/logger"
)

func main() {
	logger.Init()

	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Execute failed: %s", err)
		os.Exit(1)
	}
}
