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

	projectPath := "glad.dev/testing-repo"
	fmt.Printf("Looking for issues for '%s'\n", projectPath)
	//*
	_, err = issues.QueryAll(conf, projectPath, "https://gitlab.com")
	if err != nil {
		log.Fatalf("QueryAll failed: %s", err)
	} //*/

	_, err = issues.QuerySingle(conf, projectPath, "https://gitlab.com", "1")
	if err != nil {
		log.Fatalf("QuerySingle failed: %s", err)
	}
}
