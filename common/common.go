// Package common provides parsing and datastructures
// used by other packages in training-log-analysis.
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

// Time_ref is the reference time format used it training logs
const TimeRef = "2006-01-02 3:04PM"
const SimpleTimeRef = "2006-01-02" // Used for displaying just the date

// ExerciseY is a 'raw' struct.
// 'raw' structs are not meant to be used directly but instead to be used
// by the yaml marshaller and unmarshaller.
type ExerciseY struct {
	Name     string `name`
	Weight   string `weight,omitempty`
	Sets     string `sets`
	Reps     string `reps`
	Exertion string `exertion`
}

// EventY is a 'raw' struct.
type EventY struct {
	Name  string `name`
	Wilks string `wilks`
	Total string `total`
}

// TrainingLogY is a 'raw' struct.
type TrainingLogY struct {
	Date       string      `date`
	Time       string      `time,omitempty`
	Length     string      `length,omitempty`
	Bodyweight string      `bodyweight,omitempty`
	Event      EventY      `event,omitempty`
	Workout    []ExerciseY `workout`
	Notes      []string    `notes,omitempty`
}

// List of valid exercise names to be used in training logs.
var validExerciseNames = []string{
	// Squats
	"high bar squats",
	"low bar squats",
	"front squats",
	"paused high bar squats",
	"paused low bar squats",
	"paused front squats",
	"db lunges",
	// Pressing
	"close grip bench press",
	"bench press",
	"incline bench press",
	"tng bench press",
	"overhead press",
	"behind the neck press",
	"db incline press",
	"db flyes",
	// Pulling
	"sumo deadlift",
	"conventional deadlift",
	"stiff leg deadlift",
	"deficit conventional deadlift",
	"block pulls",
	"sumo block pulls",
	"bent over rows",
	"pendlay rows",
	"chest supported rows",
	// Back
	"pull ups",
	"chin ups",
	"lat pulldowns",
	// Arms
	"alternating db curls",
}

// IsValidExercise verifies whether it's argument is part of the
// list of valid exercises.
func IsValidExercise(name string) bool {
	for i := range validExerciseNames {
		if name == validExerciseNames[i] {
			return true
		}
	}
	return false
}

// Weight represents a weight value and it's unit.
type Weight struct {
	Value float32
	Unit  string // kg or lbs
}

// Exercise represents an exercise with Name, Weight, Sets, Reps and Exertion.
// Exertion is optional and Name should be part of ValidExerciseNames.
type Exercise struct {
	Name     string
	Weight   Weight
	Sets     int
	Reps     int
	Exertion string
}

// TrainingLog is a parsed representation of a training log
// with a valid timestamp, duration, workout, notes and bodyweight.
// If the name of Event is empty, that means the training log is
// for a normal training day.
type TrainingLog struct {
	Timestamp  time.Time
	Duration   time.Duration
	Bodyweight Weight
	Event      Event
	Workout    []Exercise
	Notes      []string
}

// Event is used for special days such as a competition day or
// mock meet
type Event struct {
	Name  string
	Wilks float32
	Total Weight
}

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
		pWilks, err := strconv.ParseFloat(rawLog.Event.Wilks, 32)
		if err != nil {
			return TrainingLog{}, fmt.Errorf("error parsing Event Wilks %s in file %s\n\t%v", rawLog.Event.Wilks, inputPath, err)
		}
		pTotal, err := ParseWeight(rawLog.Event.Total)
		if err != nil {
			return TrainingLog{}, fmt.Errorf("error parsing Event Total %s in file %s\n\t%v", rawLog.Event.Total, inputPath, err)
		}
		pEvent.Wilks = float32(pWilks)
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
