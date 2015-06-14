package main

import (
	"fmt"
	//"github.com/influxdb/influxdb/client" // if influx is hosted
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
	"training-log/common"
)

type IndexPage struct {
	Logs []string
}

type LogPage struct {
	Date       string
	Bodyweight string
	Time       string
	Length     string
	Exercises  []common.Exercise
	Notes      []string
}

const (
	outputDir = "target"
)

var (
	templates       = template.Must(template.ParseFiles("./templates/index.tmpl.html", "./templates/log.tmpl.html"))
	validPath       = regexp.MustCompile("^([0-9]+-[0-9]+-[0-9]+)$")
	indexData       IndexPage
	traininglogs    []common.TrainingLog
	trainingLogsMap map[string]common.TrainingLog = make(map[string]common.TrainingLog)
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("usage: generate-website <directory>\n")
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

	// write index file
	path := outputDir + string(os.PathSeparator) + "index.html"
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	err = templates.ExecuteTemplate(file, "index.tmpl.html", indexData)
	if err != nil {
		log.Fatal(err)
	}

	// write log files
	for date, tlog := range trainingLogsMap {
		logPath := outputDir + string(os.PathSeparator) + date + ".html"
		file, err := os.Create(logPath)
		if err != nil {
			log.Fatal(err)
		}

		logPage := LogPage{
			Date:       date,
			Bodyweight: tlog.Bodyweight.String(),
			Time:       tlog.Timestamp.Format(time.Kitchen),
			Length:     tlog.Duration.String(),
			Exercises:  tlog.Workout,
			Notes:      tlog.Notes,
		}

		err = templates.ExecuteTemplate(file, "log.tmpl.html", logPage)
		if err != nil {
			log.Fatal(err)
		}
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/date/:date", Hello)
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	log.Printf("serving on :8080")
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
		Time:       tLog.Timestamp.Format(time.Kitchen),
		Length:     tLog.Duration.String(),
		Exercises:  tLog.Workout,
		Notes:      tLog.Notes,
	}

	if path != nil && exist {
		// correct path and log exists
		err := templates.ExecuteTemplate(w, "log.tmpl.html", logPage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if path == nil || !exist {
		http.Error(w, "Page not found\n", http.StatusNotFound)
	} else {
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}
}
