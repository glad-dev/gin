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
	//*
	_, err = issues.QueryAll(conf, projectPath)
	if err != nil {
		log.Fatalf("QueryAll failed: %s", err)
	} //*/

	_, err = issues.QuerySingle(conf, projectPath, "1")
	if err != nil {
		log.Fatalf("QuerySingle failed: %s", err)
	}
}
