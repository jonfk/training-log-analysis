package main

import (
	"encoding/json"
	"fmt"
	"github.com/influxdb/influxdb/client" // Requires influx to be started
	"github.com/jonfk/training-log-analysis/common"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

var (
	templates       = template.Must(template.ParseFiles("./templates/index.tmpl.html", "./templates/log.tmpl.html"))
	validPath       = regexp.MustCompile("^([0-9]+-[0-9]+-[0-9]+)$")
	indexData       IndexPage
	traininglogs    []common.TrainingLog
	trainingLogsMap map[string]common.TrainingLog = make(map[string]common.TrainingLog)
	influxDB                                      = "traininglog"
	influxConn      *client.Client
)

type IndexPage struct {
	Logs []string
}

type LogPage struct {
	Date       string
	Bodyweight string
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("usage: statistics <directory>\n")
		os.Exit(0)
	}

	var err error
	traininglogs, err = common.ParseYamlDir(args[0])
	if err != nil {
		log.Fatal("error parsing yaml: %s\n", err)
	}

	for _, tLog := range traininglogs {
		indexData.Logs = append(indexData.Logs, tLog.Timestamp.Format(common.SimpleTimeRef))
		// Generate Map
		trainingLogsMap[tLog.Timestamp.Format(common.SimpleTimeRef)] = tLog
	}

	influxConn = initInfluxDB("localhost", 8086, "jonfk", "password")

	QueryExercise(influxConn, "low_bar_squats", "low_bar_squats_tonnage")

	router := httprouter.New()
	router.GET("/", Index)
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/log/:date", Hello)
	router.GET("/graph/", GraphData)

	fmt.Printf("Serving on 8080\n")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := templates.ExecuteTemplate(w, "index.tmpl.html", indexData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	date := ps.ByName("date")
	path := validPath.FindStringSubmatch(date)

	tLog, exist := trainingLogsMap[date]

	logPage := LogPage{
		Date:       date,
		Bodyweight: tLog.Bodyweight.String(),
	}

	err := templates.ExecuteTemplate(w, "log.tmpl.html", logPage)

	if err == nil && path != nil && exist {
		//fmt.Fprintf(w, "hello, %s!\n%v\n", date, data)
	} else if path == nil || !exist {
		http.Error(w, "Page not found\n", http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GraphData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	exercise := r.Form.Get("exercise")
	measurement := r.Form.Get("measurement")
	series := QueryExercise(influxConn, exercise, measurement)
	data, err := json.Marshal(series)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, string(data))
}

type Series struct {
	Key    string      `json:"key"`
	Values [][]float64 `json:"values"`
}

func QueryExercise(influxConn *client.Client, exercise string, measurement string) Series {
	resp, err := QueryDB(influxConn, fmt.Sprintf("select value from %s", measurement))
	if err != nil {
		log.Fatal(err)
	}
	var result Series = Series{Key: measurement}

	if resp[0].Series != nil {
		var timeIndex int
		var valueIndex int
		for i, column := range resp[0].Series[0].Columns {
			switch column {
			case "time":
				timeIndex = i
			case "value":
				valueIndex = i
			}
		}
		for _, value := range resp[0].Series[0].Values {
			timeStr := value[timeIndex].(string)
			pTime, err := time.Parse(time.RFC3339, timeStr)
			timestamp := float64(pTime.Unix())
			valueJsNum := value[valueIndex].(json.Number)
			point, err := valueJsNum.Float64()
			if err != nil {
				fmt.Printf("%#v\n", resp)
				log.Fatal("Unable to parse response from influx")
			}
			result.Values = append(result.Values, []float64{timestamp, point})
		}
	}
	return result
}

func initInfluxDB(host string, port int, influxUser, influxPwd string) *client.Client {
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
