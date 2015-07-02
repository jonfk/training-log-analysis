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
	var bodyweightDataPoints []DataPoint
	for _, bodyweight := range bodyweights {
		bodyweightDataPoints = append(bodyweightDataPoints, DataPoint{
			Date:  bodyweight.Timestamp.Format(common.SimpleTimeRef),
			Value: bodyweight.Value,
		})
	}

	bodyweightjson, err := json.Marshal(bodyweightDataPoints)

	err = ioutil.WriteFile(c.String("output")+string(os.PathSeparator)+"bodyweight.json", bodyweightjson, 493)
	if err != nil {
		log.Fatal(err)
	}

}
