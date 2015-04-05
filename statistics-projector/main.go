package main

import (
	"fmt"
	"github.com/jonfk/training-log-analysis/common"
	"log"
	"os"
)

type Statistics struct {
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("usage: statistics <directory>\n")
		os.Exit(0)
	}

	traininglogs, err := common.ParseYaml(args[0], true)
	if err != nil {
		log.Fatal("Error parsing yaml: %s\n", err)
	}

	calculateStatistics(traininglogs)
}

func calculateStatistics(logs []common.TrainingLog) Statistics {
	return Statistics{}
}
