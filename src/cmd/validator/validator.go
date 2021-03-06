package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"training-log/common"
	// "time"
)

var Verbose bool = false
var Flexible bool = false
var Output bool = false

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
			err := process(args[i])
			if err != nil {
				fmt.Printf("%s", err)
				os.Exit(1)
			}
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

func process(arg string) error {

	// fmt.Println("Processing " + arg)
	var returnErr error

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
			newErr := process(filepath.Join(arg, toProcess[i].Name()))
			if newErr != nil {
				if returnErr != nil {
					returnErr = fmt.Errorf("%s%s", returnErr, newErr)
				} else {
					returnErr = newErr
				}
			}
		}

	} else {
		data, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatalf("Error reading file %s\n", arg)
			return err
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
			return err
		}

		rawT := common.TrainingLogY{}

		err = yaml.Unmarshal(data, &rawT)
		if err != nil {
			log.Fatalf("Error parsing yaml file %s\n\t%s", arg, err)
		}

		t, err := common.ParseYaml(arg)
		if err != nil {
			returnErr = fmt.Errorf("Error parsing %s with\n\t%s\n", arg, err)
		}

		if Verbose {
			fmt.Printf("--- TrainingLog:\n%#v\n\n", t)
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
	return returnErr
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
