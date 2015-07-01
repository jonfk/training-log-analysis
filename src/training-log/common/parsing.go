package common

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ParseYaml takes a path to a file containing a
// valid training log in yaml format and returns a TrainingLog.
// If an error occured in parsing the file, an empty TrainingLog and
// the error is returned.
func ParseYaml(inputPath string) (TrainingLog, error) {
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return TrainingLog{}, fmt.Errorf("error reading file %s with\n\t%s\n", inputPath, err)
	}

	var rawLog TrainingLogY

	err = yaml.Unmarshal(data, &rawLog)
	if err != nil {
		return TrainingLog{}, fmt.Errorf("error parsing yaml file %s\n\t%s\n", inputPath, err)
	}

	// validation and parsing

	pTime, err := time.Parse(TimeRef, rawLog.Date+" "+rawLog.Time)
	if err != nil {
		return TrainingLog{}, fmt.Errorf("error parsing time in file %s\n\t%v", inputPath, err)
	}

	pDuration, err := time.ParseDuration(rawLog.Length)
	if err != nil {
		return TrainingLog{}, fmt.Errorf("error parsing duration %s in file %s\n\t%v", rawLog.Length, inputPath, err)
	}

	pBodyweight, err := ParseWeight(rawLog.Bodyweight)
	if err != nil {
		return TrainingLog{}, fmt.Errorf("error parsing bodyweight %s in file %s\n\t%v", rawLog.Bodyweight, inputPath, err)
	}

	// Event
	var pEvent Event
	if rawLog.Event.Name != "" {
		pWilks, err := strconv.ParseFloat(rawLog.Event.Wilks, 64)
		if err != nil {
			return TrainingLog{}, fmt.Errorf("error parsing Event Wilks %s in file %s\n\t%v", rawLog.Event.Wilks, inputPath, err)
		}
		pTotal, err := ParseWeight(rawLog.Event.Total)
		if err != nil {
			return TrainingLog{}, fmt.Errorf("error parsing Event Total %s in file %s\n\t%v", rawLog.Event.Total, inputPath, err)
		}
		pEvent.Wilks = float64(pWilks)
		pEvent.Total = pTotal
	}

	var exercises []Exercise
	for _, exercise := range rawLog.Workout {
		if !IsValidExercise(exercise.Name) {
			return TrainingLog{}, fmt.Errorf("error invalid exercise %s in %s\n", exercise.Name, inputPath)
		}
		pWeight, err := ParseWeight(exercise.Weight)
		if err != nil {
			return TrainingLog{}, fmt.Errorf("error parsing weight for %s in %s with error %s\n", exercise.Name, inputPath, err)
		}
		sets, err := strconv.Atoi(exercise.Sets)
		if err != nil {
			return TrainingLog{}, fmt.Errorf("error parsing sets for %s in %s with error %s\n", exercise.Name, inputPath, err)
		}
		reps, err := strconv.Atoi(exercise.Reps)
		if err != nil {
			return TrainingLog{}, fmt.Errorf("error parsing reps for %s in %s with error %s\n", exercise.Name, inputPath, err)
		}
		exercises = append(exercises, Exercise{
			Name:     exercise.Name,
			Weight:   pWeight,
			Sets:     sets,
			Reps:     reps,
			Exertion: exercise.Exertion})
	}

	return TrainingLog{
		Timestamp:  pTime,
		Duration:   pDuration,
		Bodyweight: pBodyweight,
		Event:      pEvent,
		Workout:    exercises,
		Notes:      rawLog.Notes,
	}, nil
}

// ParseYamlDir takes a path to a directory containing training log files
// and returns a []TrainingLog.
// If an error occured, it returns nil and the error.
func ParseYamlDir(dirPath string) ([]TrainingLog, error) {
	var result []TrainingLog
	toProcess, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for i := range toProcess {
		file := filepath.Join(dirPath, toProcess[i].Name())
		t, err := ParseYaml(file)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

// ParseWeight takes a string such as
//  "315lbs" or "68.5kg"
// and returns a Weight.
// If an error occurs an empty Weight and the error
// is returned
func ParseWeight(input string) (Weight, error) {
	var (
		value float64
		unit  string
		err   error
	)
	if strings.Contains(input, " ") {
		result := strings.Split(input, " ")
		value, err = strconv.ParseFloat(result[0], 64)
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
		return Weight{float64(value), unit}, nil
	}
	switch {
	case strings.Index(input, "lbs") != -1:
		result := strings.Split(input, "lbs")
		value, err = strconv.ParseFloat(result[0], 64)
		if err != nil {
			return Weight{}, err
		}
		return Weight{float64(value), "lbs"}, nil
	case strings.Index(input, "kg") != -1:
		result := strings.Split(input, "kg")
		value, err = strconv.ParseFloat(result[0], 64)
		if err != nil {
			return Weight{}, err
		}
		return Weight{float64(value), "kg"}, nil
	default:
		return Weight{}, fmt.Errorf("unknown unit in : %s", input)
	}
}

// ParseYamlDirRaw takes a path to a directory containing training log files and
// returns []TrainingLogY. It is a helper function and not meant to be used directly.
func ParseYamlDirRaw(directory string) ([]TrainingLogY, error) {
	var result []TrainingLogY
	toProcess, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for i := range toProcess {
		file := filepath.Join(directory, toProcess[i].Name())
		t := TrainingLogY{}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &t)
		if err != nil {
			return nil, fmt.Errorf("error parsing yaml file %s\n%v", file, err)
		}
		result = append(result, t)
	}
	return result, nil
}
