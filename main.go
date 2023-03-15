package main

import (
	"fmt"
	"log"
)

func main() {
	config, err := handleConfig()
	if err != nil {
		log.Fatalln(err)
	}

	projectPath := "bachelorarbeitP3/P3"
	fmt.Printf("Looking for issues for '%s'\n", projectPath)
	queryAllIssues(config, projectPath)
	//querySingleIssue(config, projectPath, 2)
}
