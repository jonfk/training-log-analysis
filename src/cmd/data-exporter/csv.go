package main

import (
	"encoding/csv"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
	"time"
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
	bodyweightFile, err := os.Create(c.String("output") + string(os.PathSeparator) + "bodyweight.csv")
	if err != nil {
		log.Fatal(err)
	}
	bodyweightCSVWriter := csv.NewWriter(bodyweightFile)

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
	bodyweightCSVWriter.Write([]string{"date", "bodyweight", "unit"})
	allCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight", "unit"})
	allSquatsCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight", "unit"})
	compSquatsCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight", "unit"})
	benchCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight", "unit"})
	deadliftCSVWriter.Write([]string{"date", "exercise", "sets", "reps", "weight", "unit"})

	// parse training logs
	traininglogs, err := common.ParseYamlDir(c.String("input"))
	if err != nil {
		log.Fatal("error parsing yaml: %s\n", err)
	}

	for i := range traininglogs {
		exercises := traininglogs[i].Workout
		var date string = traininglogs[i].Timestamp.Format(time.RFC3339)

		err := bodyweightCSVWriter.Write([]string{date, fmt.Sprintf("%.2f",
			traininglogs[i].Bodyweight.Value), traininglogs[i].Bodyweight.Unit})
		if err != nil {
			log.Fatal(err)
		}

		for _, ex := range exercises {
			err := allCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets),
				fmt.Sprintf("%d", ex.Reps), fmt.Sprintf("%.2f", ex.Weight.Value), ex.Weight.Unit})
			if err != nil {
				log.Fatal(err)
			}
		}

		squats := common.FilterVariation(common.SquatVariation, traininglogs[i])
		for _, ex := range squats {
			err := allSquatsCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets),
				fmt.Sprintf("%d", ex.Reps), fmt.Sprintf("%.2f", ex.Weight.Value), ex.Weight.Unit})
			if err != nil {
				log.Fatal(err)
			}
		}
		compSquats := common.Filter(traininglogs[i], "low bar squats", "belted low bar squats")
		for _, ex := range compSquats {
			err := compSquatsCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets),
				fmt.Sprintf("%d", ex.Reps), fmt.Sprintf("%.2f", ex.Weight.Value), ex.Weight.Unit})
			if err != nil {
				log.Fatal(err)
			}
		}
		bench := common.FilterVariation(common.BenchVariation, traininglogs[i])
		for _, ex := range bench {
			err := benchCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets),
				fmt.Sprintf("%d", ex.Reps), fmt.Sprintf("%.2f", ex.Weight.Value), ex.Weight.Unit})
			if err != nil {
				log.Fatal(err)
			}
		}
		deadlifts := common.FilterVariation(common.DeadliftVariation, traininglogs[i])
		for _, ex := range deadlifts {
			err := deadliftCSVWriter.Write([]string{date, ex.Name, fmt.Sprintf("%d", ex.Sets),
				fmt.Sprintf("%d", ex.Reps), fmt.Sprintf("%.2f", ex.Weight.Value), ex.Weight.Unit})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	bodyweightCSVWriter.Flush()
	allCSVWriter.Flush()
	allSquatsCSVWriter.Flush()
	compSquatsCSVWriter.Flush()
	benchCSVWriter.Flush()
	deadliftCSVWriter.Flush()
}
