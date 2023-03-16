package main

import (
	"fmt"
	"log"

	"gn/config"
	"gn/issues"
)

func main() {
	conf, err := config.Get()
	if err != nil {
		log.Fatalln(err)
	}

	projectPath := "bachelorarbeitP3/P3"
	fmt.Printf("Looking for issues for '%s'\n", projectPath)
	issues.QueryAll(conf, projectPath)
	//querySingleIssue(config, projectPath, 2)
}
