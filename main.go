package main

import (
	"fmt"
	"os"

	"gn/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Execute failed: %s", err)
		os.Exit(1)
	}
}
