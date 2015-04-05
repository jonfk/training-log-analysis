package common

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Exercise struct {
	Name     string `name`
	Weight   string `weight,omitempty`
	Sets     string `sets`
	Reps     string `reps`
	Exertion string `exertion`
}

type TrainingLog struct {
	Date       string     `date`
	Time       string     `time,omitempty`
	Length     string     `length,omitempty`
	Bodyweight string     `bodyweight,omitempty`
	Event      string     `event,omitempty`
	Wilks      string     `wilks,omitempty`
	Total      string     `total,omitempty`
	Workout    []Exercise `workout`
	Notes      []string   `notes,omitempty`
}

func ParseYaml(directory string, print bool) ([]TrainingLog, error) {
	var result []TrainingLog
	toProcess, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for i := range toProcess {
		file := filepath.Join(directory, toProcess[i].Name())
		t := TrainingLog{}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &t)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error parsing yaml file %s\n%v", file, err))
		}
		result = append(result, t)
	}

	if print {
		fmt.Printf("[ParseYaml] TrainingLogs:\n")
		for i := range result {
			//fmt.Printf("- %#v\n", result[i])
			spew.Printf("- %v\n", result[i])
		}
	}
	return result, nil
}
