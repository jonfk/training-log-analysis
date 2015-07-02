package main

import (
	"encoding/csv"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"training-log/common"
)

func ExportCSV(c *cli.Context) {
	if !c.IsSet("input") {
		cli.ShowAppHelp(c)
		return
	}
	// create dir if it doesn't exist
	_, err := os.Stat(c.String("output"))
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(c.String("output"), 493)
		} else {
			log.Fatal(err)
		}
	}

	// Create output file and csv writer
	allFile, err := os.Create(c.String("output") + string(os.PathSeparator) + "all.csv")
	if err != nil {
		log.Fatal(err)
	}
	allCSVWriter := csv.NewWriter(allFile)

	allSquatsFile, err := os.Create(c.String("output") + string(os.PathSeparator) + "allSquats.csv")
	if err != nil {
		log.Fatal(err)
	}
	allSquatsCSVWriter := csv.NewWriter(allSquatsFile)

	compSquatsFile, err := os.Create(c.String("output") + string(os.PathSeparator) + "compSquats.csv")
	if err != nil {
		log.Fatal(err)
	}
	compSquatsCSVWriter := csv.NewWriter(compSquatsFile)

	benchFile, err := os.Create(c.String("output") + string(os.PathSeparator) + "bench.csv")
	if err != nil {
		log.Fatal(err)
	}
	benchCSVWriter := csv.NewWriter(benchFile)

	deadliftFile, err := os.Create(c.String("output") + string(os.PathSeparator) + "deadlift.csv")
	if err != nil {
		log.Fatal(err)
	}
	deadliftCSVWriter := csv.NewWriter(deadliftFile)

	// Write Headers
	allCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight"})
	allSquatsCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight"})
	compSquatsCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight"})
	benchCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight"})
	deadliftCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight"})

	// parse training logs
	traininglogs, err := common.ParseYamlDir(c.String("input"))
	if err != nil {
		log.Fatal("error parsing yaml: %s\n", err)
	}

	for i := range traininglogs {
		exercises := traininglogs[i].Workout
		var date string = traininglogs[i].SimpleTime()
		for _, ex := range exercises {
			err := allCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets), fmt.Sprintf("%d", ex.Reps), ex.Weight.String()})
			if err != nil {
				log.Fatal(err)
			}
		}

		squats := common.FilterVariation(common.SquatVariation, traininglogs[i])
		for _, ex := range squats {
			err := allSquatsCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets), fmt.Sprintf("%d", ex.Reps), ex.Weight.String()})
			if err != nil {
				log.Fatal(err)
			}
		}
		compSquats := common.Filter(traininglogs[i], "low bar squats", "belted low bar squats")
		for _, ex := range compSquats {
			err := compSquatsCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets), fmt.Sprintf("%d", ex.Reps), ex.Weight.String()})
			if err != nil {
				log.Fatal(err)
			}
		}
		bench := common.FilterVariation(common.BenchVariation, traininglogs[i])
		for _, ex := range bench {
			err := benchCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets), fmt.Sprintf("%d", ex.Reps), ex.Weight.String()})
			if err != nil {
				log.Fatal(err)
			}
		}
		deadlifts := common.FilterVariation(common.DeadliftVariation, traininglogs[i])
		for _, ex := range deadlifts {
			err := deadliftCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets), fmt.Sprintf("%d", ex.Reps), ex.Weight.String()})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	allCSVWriter.Flush()
	allSquatsCSVWriter.Flush()
	compSquatsCSVWriter.Flush()
	benchCSVWriter.Flush()
	deadliftCSVWriter.Flush()
}
