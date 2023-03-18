package main

import (
	"fmt"
	"gn/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Execute failed: %s", err)
	}
}
