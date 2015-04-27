package main

import (
	"fmt"
	//"github.com/influxdb/influxdb/client" // if influx is hosted
	"github.com/jonfk/training-log-analysis/common"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var (
	templates       = template.Must(template.ParseFiles("./templates/index.tmpl.html", "./templates/log.tmpl.html"))
	validPath       = regexp.MustCompile("^([0-9]+-[0-9]+-[0-9]+)$")
	indexData       IndexPage
	traininglogs    []common.TrainingLog
	trainingLogsMap map[string]common.TrainingLog = make(map[string]common.TrainingLog)
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

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/:date", Hello)

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
