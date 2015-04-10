package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"github.com/jonfk/training-log-analysis/common"
	"log"
	"os"
	"strings"
)

const (
	User         = "jonathan"
	InfluxDBHost = "localhost"
	InfluxDBPort = 8086
	InfluxDB     = "traininglog"
)

var (
	InfluxUser = "jonfk"
	InfluxPwd  = "password"
)

func init() {
	if os.Getenv("INFLUX_USER") == "" || os.Getenv("INFLUX_PWD") == "" {
		log.Println("Env INFLUX_USER and INFLUX_PWD are unset, using default user and password")
	} else {
		InfluxUser = os.Getenv("INFLUX_USER")
		InfluxPwd = os.Getenv("INFLUX_PWD")
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
		log.Fatal("Error parsing yaml: %s\n", err)
	}

	connection := initInfluxdB(InfluxDBHost, InfluxDBPort)

	ProjectExerciseIntensity(connection, "low bar squats", traininglogs)
	ProjectExerciseTonnage(connection, "low bar squats", traininglogs)
	ProjectExerciseBarLifts(connection, "low bar squats", traininglogs)
	ProjectExerciseIntensity(connection, "bench press", traininglogs)
	ProjectExerciseTonnage(connection, "bench press", traininglogs)
	ProjectExerciseBarLifts(connection, "bench press", traininglogs)
}

// Projections
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
	writePoints(conn, metricsToBeInserted)
}

func ProjectExerciseTonnage(conn *client.Client, name string, logs []common.TrainingLog) {
	var metricsToBeInserted []ExerciseMetric
	var unit string
	for _, trainLog := range logs {
		var projectDay bool = false
		var tonnagePerDay float32 = 0
		for _, exercise := range trainLog.Workout {
			if exercise.Name == name {
				projectDay = true
				tonnagePerDay += (exercise.Weight.Value * float32(exercise.Reps) * float32(exercise.Sets))
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
	writePoints(conn, metricsToBeInserted)
}

func ProjectExerciseBarLifts(conn *client.Client, name string, logs []common.TrainingLog) {
	var metricsToBeInserted []ExerciseMetric
	var unit string
	for _, trainLog := range logs {
		var projectDay bool = false
		var barliftsPerDay int = 0
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
				Value:     float32(barliftsPerDay),
				Unit:      unit,
				Timestamp: trainLog.Timestamp,
			}
			log.Printf("[Barlifts]Projecting %v for %s", name, trainLog.Timestamp)
			metricsToBeInserted = append(metricsToBeInserted, metric)
		}
	}
	writePoints(conn, metricsToBeInserted)
}
