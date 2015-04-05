package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jonfk/training-log-analysis/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var Verbose bool = false
var Flexible bool = false
var Output bool = false

var validExerciseNames []string = []string{
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

func main() {
	// Parse Command Line arguments
	args := os.Args[1:] // Ignore program name

	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}

	for i := range args {
		if isLongFlag(args[i]) {
			switch args[i] {
			case "--help":
				printUsage()
				os.Exit(0)
			case "--verbose":
				Verbose = true
			case "--flexible":
				Flexible = true
			case "--output":
				Output = true
			}
		} else if isFlag(args[i]) {
			for j := range args[i] {
				switch args[i][j] {
				case 'h':
					printUsage()
					os.Exit(0)
				case 'v':
					Verbose = true
				case 'f':
					Flexible = true
				case 'o':
					Output = true
				}
			}
		} else {
			process(args[i])
		}
	}

}

func isLongFlag(arg string) bool {
	idx := strings.Index(arg, "--")
	if idx == 0 {
		return true
	}
	return false
}

func isFlag(arg string) bool {
	idx := strings.Index(arg, "-")
	if idx == 0 {
		return true
	}
	return false
}

func process(arg string) {

	// fmt.Println("Processing " + arg)

	isDir, err := IsDirectory(arg)

	if err != nil {
		log.Fatalf("Invalid path to file or directory\n")
	}
	if isDir {
		// fmt.Println("isdirectory")
		toProcess, err := ioutil.ReadDir(arg)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		for i := range toProcess {
			process(filepath.Join(arg, toProcess[i].Name()))
		}

	} else {
		data, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatalf("Error reading file %s\n", arg)
			return
		}

		if Flexible {
			m := make(map[interface{}]interface{})
			err := yaml.Unmarshal(data, &m)
			if err != nil {
				log.Fatalf("Error parsing yaml file %s\n%v", arg, err)
			}
			if Verbose {
				fmt.Printf("--- Flexible:\n%#v\n\n", m)
			}
			return
		}

		t := common.TrainingLog{}

		err = yaml.Unmarshal(data, &t)
		if err != nil {
			log.Fatalf("Error parsing yaml file %s\n%v", arg, err)
		}

		// validate exercises
		for i := range t.Workout {
			if !isValidExercise(t.Workout[i].Name) {
				fmt.Printf("File %s contains invalid exercise %s\n", arg, t.Workout[i].Name)
			}
			//fmt.Printf("Exercise: %s\n", t.Workout[i].Name)
		}

		// validate time
		const reference = "2006-01-02 3:04PM"
		_, err = time.Parse(reference, t.Date+" "+t.Time)
		if err != nil {
			fmt.Printf("Error parsing time in %s\n", arg)
		}

		if Verbose {
			//fmt.Printf("--- TrainingLog:\n%#v\n\n", t)
			spew.Printf("--- TrainingLog:\n%#v\n\n", t)
		}

		// Output flag set
		if Output {
			d, err := yaml.Marshal(&t)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			fmt.Printf("%s\n", string(d))
		}
	}

}

func isValidExercise(name string) bool {
	for i := range validExerciseNames {
		if name == validExerciseNames[i] {
			return true
		}
	}
	return false
}

func printUsage() {
	usage := `
Usage [-hvfo][--help][--verbose][--flexible][--output] [Arguments...]

Options:
--help, -h       show this message, then exit
--verbose, -v    Print the internal datastructure the yaml mapped to
--flexible, -f   Verify the yaml rather than the template
--output, -o     Output the idealized template
`
	fmt.Printf("%s\n", usage)
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}
