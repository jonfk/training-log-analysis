package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"github.com/jonfk/training-log-analysis/common"
	"log"
	"os"
	"strings"
	"time"
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

	traininglogs, err := common.ParseYamlDirRaw(args[0])
	if err != nil {
		log.Fatal("Error parsing yaml: %s\n", err)
	}

	connection := initInfluxdB(InfluxDBHost, InfluxDBPort)

	ProjectExerciseIntensity(connection, "low bar squats", traininglogs)
}

// Projections
func ProjectExerciseIntensity(conn *client.Client, name string, logs []common.TrainingLogY) {
	var metricsToBeInserted []ExerciseMetric
	for _, trainLog := range logs {
	Loop1:
		for _, exercise := range trainLog.Workout {
			if exercise.Name == name {
				pTime, err := time.Parse(common.Time_ref, trainLog.Date+" "+trainLog.Time)
				if err != nil {
					log.Printf("Error parsing time in %s\n", trainLog.Date)
					continue
				}

				weight, err := common.ParseWeight(exercise.Weight)
				if err != nil {
					log.Printf("Error parsing weight for %s in %s with error %s\n", exercise.Name, trainLog.Date, err)
					continue
				}
				// Intensity metric
				metric := ExerciseMetric{
					Name:      strings.Replace(exercise.Name, " ", "_", -1) + "_intensity",
					Username:  User,
					Value:     weight.Value,
					Unit:      weight.Unit,
					Timestamp: pTime,
				}
				// If 2 values with same timestamp are inserted the first one will
				// be overwritten by influx
				for _, mtrs := range metricsToBeInserted {
					if mtrs.Timestamp.Equal(metric.Timestamp) && mtrs.Value > metric.Value {
						continue Loop1
					}
				}
				log.Printf("Projecting %v for %s", exercise, pTime)
				metricsToBeInserted = append(metricsToBeInserted, metric)
			}
		}
	}
	writePoints(conn, metricsToBeInserted)
}

// func ProjectExerciseTonnage(conn *client.Client, name string, logs []common.TrainingLogY) {
// 	var metricsToBeInserted []ExerciseMetric
// 	for _, trainLog := range logs {
// 		var tonnagePerDay float32
// 		for _, exercise := range trainLog.Workout {
// 			if exercise.Name == name {
// 			}
// 		}
// 	}
// }
