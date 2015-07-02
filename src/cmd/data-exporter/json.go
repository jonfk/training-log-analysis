package main

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"training-log/common"
	"training-log/projections"
)

type DataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

func ExportJson(c *cli.Context) {
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

	// Parse Training logs
	traininglogs, err := common.ParseYamlDir(c.String("input"))
	if err != nil {
		log.Fatal("error parsing yaml: %s\n", err)
	}

	// Project Bodyweights
	bodyweights := projections.ProjectBodyWeight(traininglogs)
	bodyweightjson, err := projectionToJson(bodyweights)
	if err != nil {
		log.Fatal(err)
	}

	// Project training duration
	trainingDurations := projections.ProjectTrainingDuration(traininglogs)
	trainingDurationjson, err := projectionToJson(trainingDurations)
	if err != nil {
		log.Fatal(err)
	}

	// Project Squats
	beltedLowBarSquatsIntensity := projections.ProjectExerciseIntensity(traininglogs, "belted low bar squats")
	beltedLowBarSquatsIntensityJson, err := projectionToJson(beltedLowBarSquatsIntensity)
	if err != nil {
		log.Fatal(err)
	}
	lowBarSquatsIntensity := projections.ProjectExerciseIntensity(traininglogs, "low bar squats")
	lowBarSquatsIntensityJson, err := projectionToJson(lowBarSquatsIntensity)
	if err != nil {
		log.Fatal(err)
	}

	squatsTonnage := projections.ProjectExerciseTonnage(traininglogs, "low bar squats", "high bar squats", "front squats", "belted low bar squats")
	squatsTonnageJson, err := projectionToJson(squatsTonnage)
	if err != nil {
		log.Fatal(err)
	}

	// Project Bench
	benchIntensity := projections.ProjectExerciseIntensity(traininglogs, "bench press", "close grip bench press")
	benchIntensityJson, err := projectionToJson(benchIntensity)
	if err != nil {
		log.Fatal(err)
	}

	benchTonnage := projections.ProjectExerciseTonnage(traininglogs, "bench press", "close grip bench press")
	benchTonnageJson, err := projectionToJson(benchTonnage)
	if err != nil {
		log.Fatal(err)
	}

	// Project Deadlift
	deadliftIntensity := projections.ProjectExerciseIntensity(traininglogs, "sumo deadlift", "conventional deadlift")
	deadliftIntensityJson, err := projectionToJson(deadliftIntensity)
	if err != nil {
		log.Fatal(err)
	}

	deadliftTonnage := projections.ProjectExerciseTonnage(traininglogs, "sumo deadlift", "conventional deadlift")
	deadliftTonnageJson, err := projectionToJson(deadliftTonnage)
	if err != nil {
		log.Fatal(err)
	}

	// Write Json to file
	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"bodyweight.json", bodyweightjson, 493)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"training_duration.json", trainingDurationjson, 493)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"belted_low_bar_squats_intensity.json", beltedLowBarSquatsIntensityJson, 493)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"low_bar_squats_intensity.json", lowBarSquatsIntensityJson, 493)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"squats_tonnage.json", squatsTonnageJson, 493)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"bench_intensity.json", benchIntensityJson, 493)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"bench_tonnage.json", benchTonnageJson, 493)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"deadlift_intensity.json", deadliftIntensityJson, 493)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"deadlift_tonnage.json", deadliftTonnageJson, 493)
	if err != nil {
		log.Fatal(err)
	}

}

func projectionToJson(dataPoints []projections.DataPoint) ([]byte, error) {
	var serializablePoints []DataPoint
	for _, point := range dataPoints {
		serializablePoints = append(serializablePoints, DataPoint{
			Date:  point.Timestamp.Format(common.SimpleTimeRef),
			Value: point.Value,
		})
	}
	return json.Marshal(serializablePoints)
}
