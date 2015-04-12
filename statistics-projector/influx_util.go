package main

import (
	"fmt"
	"github.com/influxdb/influxdb/client"
	"log"
	"net/url"
	"time"
)

// ExerciseMetric is a convinience struct to be used by
// WritePoints.
type ExerciseMetric struct {
	Name      string
	Username  string
	Unit      string
	Value     float32
	Set       string
	Timestamp time.Time
}

func initInfluxdB(host string, port int) *client.Client {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", host, port))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: influxUser,
		Password: influxPwd,
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

// QueryDB is a convenience function to query the database
func QueryDB(con *client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: influxDB,
	}
	if response, err := con.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}

// WritePoints is a utility function used to write a []ExerciseMetric
// to the database in a batch using the connection supplied.
// Error is nil unless the write fails.
func WritePoints(con *client.Client, metrics []ExerciseMetric) error {
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
		Database:        influxDB,
		RetentionPolicy: "default",
	}
	_, err := con.Write(bps)
	if err != nil {
		return err
	}
	return nil
}
