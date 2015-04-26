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
)

var templates = template.Must(template.ParseFiles("./templates/index.html"))
var traininglogs []common.TrainingLog
var indexData IndexData

type IndexData struct {
	Logs []string
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
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("%v", traininglogs)
	err := templates.ExecuteTemplate(w, "index.html", indexData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
