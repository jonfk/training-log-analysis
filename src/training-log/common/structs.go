package common

import (
	"fmt"
	"time"
)

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

// Weight represents a weight value and it's unit.
type Weight struct {
	Value float64
	Unit  string // kg or lbs
}

func (w Weight) String() string {
	return fmt.Sprintf("%.1f %s", w.Value, w.Unit)
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

func (t TrainingLog) SimpleTime() string {
	return t.Timestamp.Format(SimpleTimeRef)
}

// Event is used for special days such as a competition day or
// mock meet
type Event struct {
	Name  string
	Wilks float64
	Total Weight
}
