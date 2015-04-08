package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"github.com/jonfk/training-log-analysis/common"
	"log"
	"net/url"
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

type ExerciseMetric struct {
	Name      string
	Username  string
	Unit      string
	Value     float32
	Set       string
	Timestamp time.Time
}

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

	traininglogs, err := common.ParseYaml(args[0], true)
	if err != nil {
		log.Fatal("Error parsing yaml: %s\n", err)
	}

	connection := initInfluxdB(InfluxDBHost, InfluxDBPort)

	ProjectExerciseIntensity(connection, "low bar squats", traininglogs)
}

func ProjectExerciseIntensity(conn *client.Client, name string, logs []common.TrainingLog) {
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

				log.Printf("Projecting %v for %s", exercise, pTime)

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
				metricsToBeInserted = append(metricsToBeInserted, metric)
			}
		}
	}
	writePoints(conn, metricsToBeInserted)
}

func initInfluxdB(host string, port int) *client.Client {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: InfluxUser,
		Password: InfluxPwd,
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}

	dur, ver, err := con.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connection verified %v, %s", dur, ver)
	return con
}

// queryDB convenience function to query the database
func queryDB(con *client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: InfluxDB,
	}
	if response, err := con.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}

func writePoints(con *client.Client, metrics []ExerciseMetric) error {
	var (
		points = make([]client.Point, len(metrics))
	)

	for i := range metrics {
		points[i] = client.Point{
			Name: metrics[i].Name,
			Tags: map[string]string{
				"username": metrics[i].Username,
				"unit":     metrics[i].Unit,
			},
			Fields: map[string]interface{}{
				"value": metrics[i].Value,
			},
			Timestamp: metrics[i].Timestamp,
			Precision: "s",
		}
	}

	bps := client.BatchPoints{
		Points:          points,
		Database:        InfluxDB,
		RetentionPolicy: "default",
	}
	_, err := con.Write(bps)
	if err != nil {
		return err
	}
	return nil
}
