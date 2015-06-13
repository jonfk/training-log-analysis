// statistics-projector saves various statistics from training logs
// to influxdb for later analysis and display.
//
// statistics-projector uses INFLUX_USER and INFLUX_PWD environment variables
// if they are set for username and password of the influxdb database.
//
//  Usage: statistics-projector <directory>
//  <directory> is a directory with training log files
//
// statistics-projector currently does the following projections:
//  - BarLifts: number of reps per exercise per day
//  - Intensity: intensity of highest working set for given exercise per day
//  - Tonnage: total weight lifted by multiplying the weight, reps and sets per exercise per day
//  - TrainingDuration: training duration per session
//  - Frequency: frequency of exercise, a tic metric of value 1 per day
//  - Bodyweight: bodyweight per day
//
package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"os"
	"strings"
	"training-log/common"
)

const (
	// Hard coded for me but should be changed eventually
	User         = "jonathan"
	influxDBHost = "localhost"
	influxDBPort = 8086
	influxDB     = "traininglog"
)

var (
	influxUser = "jonfk"
	influxPwd  = "password"
)

func init() {
	if os.Getenv("INFLUX_USER") == "" || os.Getenv("INFLUX_PWD") == "" {
		log.Println("Env INFLUX_USER and INFLUX_PWD are unset, using default user and password")
	} else {
		influxUser = os.Getenv("INFLUX_USER")
		influxPwd = os.Getenv("INFLUX_PWD")
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("usage: statistics <directory>\n")
		os.Exit(0)
	}

	traininglogs, err := common.ParseYamlDir(args[0])
	if err != nil {
		log.Fatal("error parsing yaml: %s\n", err)
	}

	connection := initInfluxdB(influxDBHost, influxDBPort)

	exercisesToProject := []string{"low bar squats", "bench press", "sumo deadlift", "conventional deadlift", "overhead press"}
	// Exercise Projections
	for _, exercise := range exercisesToProject {
		ProjectExerciseIntensity(connection, exercise, traininglogs)
		ProjectExerciseTonnage(connection, exercise, traininglogs)
		ProjectExerciseBarLifts(connection, exercise, traininglogs)
		ProjectExerciseFrequency(connection, exercise, traininglogs)
	}

	// Training stats
	ProjectTrainingDuration(connection, traininglogs)
	ProjectBodyWeight(connection, traininglogs)
}

// ------------ Projections -------------

// ProjectExerciseIntensity takes the name of a valid exercise, a connection to an influxdb database
// and a list of []common.TrainingLog. It saves the intensity of the highest working set of that exercise
// per day.
func ProjectExerciseIntensity(conn *client.Client, name string, logs []common.TrainingLog) {
	var metricsToBeInserted []ExerciseMetric
	for _, trainLog := range logs {
	Loop1:
		for _, exercise := range trainLog.Workout {
			if exercise.Name == name {
				// Intensity metric
				metric := ExerciseMetric{
					Name:      strings.Replace(exercise.Name, " ", "_", -1) + "_intensity",
					Username:  User,
					Value:     exercise.Weight.Value,
					Unit:      exercise.Weight.Unit,
					Timestamp: trainLog.Timestamp,
				}
				// If 2 values with same timestamp are inserted the first one will
				// be overwritten by influx
				for _, mtrs := range metricsToBeInserted {
					if mtrs.Timestamp.Equal(metric.Timestamp) && mtrs.Value > metric.Value {
						continue Loop1
					}
				}
				log.Printf("[Intensity]Projecting %v for %s", exercise, trainLog.Timestamp)
				metricsToBeInserted = append(metricsToBeInserted, metric)
			}
		}
	}
	err := WritePoints(conn, metricsToBeInserted)
	if err != nil {
		log.Fatal(err)
	}
}

// ProjectExerciseTonnage takes the name of a valid exercise, a connection to influxdb and
// a []common.TrainingLog and saves the total weight of the exercise per day.
// It assumes the exercise has the same unit on the same day.
func ProjectExerciseTonnage(conn *client.Client, name string, logs []common.TrainingLog) {
	var metricsToBeInserted []ExerciseMetric
	var unit string
	for _, trainLog := range logs {
		var projectDay = false
		var tonnagePerDay float64
		for _, exercise := range trainLog.Workout {
			if exercise.Name == name {
				projectDay = true
				tonnagePerDay += (exercise.Weight.Value * float64(exercise.Reps) * float64(exercise.Sets))
				unit = exercise.Weight.Unit
			}
		}
		if projectDay {
			// Tonnage metric
			metric := ExerciseMetric{
				Name:      strings.Replace(name, " ", "_", -1) + "_tonnage",
				Username:  User,
				Value:     tonnagePerDay,
				Unit:      unit,
				Timestamp: trainLog.Timestamp,
			}
			log.Printf("[Tonnage]Projecting %v for %s", name, trainLog.Timestamp)
			metricsToBeInserted = append(metricsToBeInserted, metric)
		}
	}
	err := WritePoints(conn, metricsToBeInserted)
	if err != nil {
		log.Fatal(err)
	}
}

// ProjectExerciseBarLifts takes a valid exercise name, a connection to influxdb and a
// []common.TrainingLog. It saves the total number of reps done on the exercise per day.
func ProjectExerciseBarLifts(conn *client.Client, name string, logs []common.TrainingLog) {
	var metricsToBeInserted []ExerciseMetric
	var unit string
	for _, trainLog := range logs {
		var projectDay = false
		var barliftsPerDay = 0
		for _, exercise := range trainLog.Workout {
			if exercise.Name == name {
				projectDay = true
				barliftsPerDay += (exercise.Reps * exercise.Sets)
				unit = exercise.Weight.Unit
			}
		}
		if projectDay {
			// Barlifts metric
			metric := ExerciseMetric{
				Name:      strings.Replace(name, " ", "_", -1) + "_barlifts",
				Username:  User,
				Value:     float64(barliftsPerDay),
				Unit:      unit,
				Timestamp: trainLog.Timestamp,
			}
			log.Printf("[Barlifts]Projecting %v for %s", name, trainLog.Timestamp)
			metricsToBeInserted = append(metricsToBeInserted, metric)
		}
	}
	err := WritePoints(conn, metricsToBeInserted)
	if err != nil {
		log.Fatal(err)
	}
}

// ProjectTrainingDuration takes a connection to influxdb and a
// []common.TrainingLog. It saves training duration per day.
func ProjectTrainingDuration(conn *client.Client, logs []common.TrainingLog) {
	var metricsToBeInserted []ExerciseMetric
	for _, trainLog := range logs {
		metric := ExerciseMetric{
			Name:      "training_duration",
			Username:  User,
			Value:     float64(trainLog.Duration.Hours()),
			Unit:      "hours",
			Timestamp: trainLog.Timestamp,
		}
		log.Printf("[Training Duration]Projecting %s", trainLog.Timestamp)
		metricsToBeInserted = append(metricsToBeInserted, metric)
	}
	err := WritePoints(conn, metricsToBeInserted)
	if err != nil {
		log.Fatal(err)
	}
}

// ProjectExerciseFrequency takes a valid exercise name, a connection to influxdb and a
// []common.TrainingLog. It saves the frequency of training the exercise.
func ProjectExerciseFrequency(conn *client.Client, name string, logs []common.TrainingLog) {
	var metricsToBeInserted []ExerciseMetric
	for _, trainLog := range logs {
		for _, exercise := range trainLog.Workout {
			if exercise.Name == name {
				metric := ExerciseMetric{
					Name:      strings.Replace(name, " ", "_", -1) + "_frequency",
					Username:  User,
					Value:     1,
					Timestamp: trainLog.Timestamp,
				}
				log.Printf("[Exercise Frequency]Projecting %v for %s", name, trainLog.Timestamp)
				metricsToBeInserted = append(metricsToBeInserted, metric)
			}
		}
	}
	err := WritePoints(conn, metricsToBeInserted)
	if err != nil {
		log.Fatal(err)
	}
}

func ProjectBodyWeight(conn *client.Client, logs []common.TrainingLog) {
	var metricsToBeInserted []ExerciseMetric
	for _, trainLog := range logs {
		metric := ExerciseMetric{
			Name:      "bodyweight",
			Username:  User,
			Value:     float64(trainLog.Bodyweight.Value),
			Unit:      trainLog.Bodyweight.Unit,
			Timestamp: trainLog.Timestamp,
		}
		log.Printf("[Bodyweight]Projecting %v on %s", trainLog.Bodyweight.Value, trainLog.Timestamp)
		metricsToBeInserted = append(metricsToBeInserted, metric)
	}
	err := WritePoints(conn, metricsToBeInserted)
	if err != nil {
		log.Fatal(err)
	}
}
