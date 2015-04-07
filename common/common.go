package common

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

const Time_ref = "2006-01-02 3:04PM"

// Yaml structs
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

// Parsed Structs
type Weight struct {
	Value float32
	Unit  string // kg or lbs
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

func ParseWeight(input string) (Weight, error) {
	var (
		value float64
		unit  string
		err   error
	)
	if strings.Contains(input, " ") {
		result := strings.Split(input, " ")
		value, err = strconv.ParseFloat(result[0], 32)
		if err != nil {
			return Weight{}, err
		}
		switch strings.ToLower(result[1]) {
		case "lbs":
			unit = "lbs"
		case "kg":
			unit = "kg"
		case "kgs":
			unit = "kg"
		default:
			return Weight{}, errors.New("Unknown unit: " + result[1])
		}
		return Weight{float32(value), unit}, nil
	}
	switch {
	case strings.Index(input, "lbs") != -1:
		result := strings.Split(input, "lbs")
		value, err = strconv.ParseFloat(result[0], 32)
		if err != nil {
			return Weight{}, err
		}
		return Weight{float32(value), "lbs"}, nil
	case strings.Index(input, "kg") != -1:
		result := strings.Split(input, "kg")
		value, err = strconv.ParseFloat(result[0], 32)
		if err != nil {
			return Weight{}, err
		}
		return Weight{float32(value), "kg"}, nil
	default:
		return Weight{}, errors.New("Unknown unit in : " + input)
	}
}
