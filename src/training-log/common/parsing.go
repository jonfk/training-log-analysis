package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
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
		path := filepath.Join(dirPath, toProcess[i].Name())

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}
		isLogDir, err := regexp.Match("^(19|20)\\d\\d$", []byte(toProcess[i].Name()))
		if err != nil {
			return nil, err
		}
		isLogFile, err := regexp.Match("^(19|20)\\d\\d([-])(0[1-9]|1[012])*([-])(0[1-9]|[12][0-9]|3[01])([-])*(.)*.yml$", []byte(toProcess[i].Name()))
		if err != nil {
			return nil, err
		}

		if fInfo.Mode().IsDir() && isLogDir {
			logs, err := ParseYamlDir(path)
			if err != nil {
				return nil, err
			}
			if logs != nil {
				result = append(result, logs...)
			}
		} else if isLogFile {
			t, err := ParseYaml(path)
			if err != nil {
				return nil, err
			}
			result = append(result, t)
		} else {
			//log.Println("skipping " + path)
			// skip
		}
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
		path := filepath.Join(directory, toProcess[i].Name())

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		fInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}

		isLogDir, err := regexp.Match("^(19|20)\\d\\d$", []byte(toProcess[i].Name()))
		if err != nil {
			return nil, err
		}
		isLogFile, err := regexp.Match("^(19|20)\\d\\d([-])(0[1-9]|1[012])*([-])(0[1-9]|[12][0-9]|3[01])([-])*(.)*.yml$", []byte(toProcess[i].Name()))
		if err != nil {
			return nil, err
		}

		if fInfo.Mode().IsDir() && isLogDir {
			logs, err := ParseYamlDirRaw(path)
			if err != nil {
				return nil, err
			}
			if logs != nil {
				result = append(result, logs...)
			}
		} else if isLogFile {
			t := TrainingLogY{}
			buf := new(bytes.Buffer)
			_, err := io.Copy(buf, file)
			if err != nil {
				return nil, err
			}

			data := buf.Bytes()
			err = yaml.Unmarshal(data, &t)
			if err != nil {
				return nil, fmt.Errorf("error parsing yaml file %s\n%v", file, err)
			}
			result = append(result, t)
		} else {

			// skip if not yaml file
		}
	}
	return result, nil
}
