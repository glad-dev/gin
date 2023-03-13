package main

import (
	"log"
)

func main() {
	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	projectPath := "bachelorarbeitP3/P3"
	queryAllIssues(config, projectPath)
	//querySingleIssue(config, projectPath, //20)
}
