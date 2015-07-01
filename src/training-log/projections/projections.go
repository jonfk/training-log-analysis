package projections

import (
	"training-log/common"
)

func ProjectBodyWeight(logs []common.TrainingLog) []DataPoint {
	var dataPoints []DataPoint
	for _, trainLog := range logs {
		point := DataPoint{
			Timestamp: trainLog.Timestamp,
			Value:     float64(trainLog.Bodyweight.Value),
			Unit:      trainLog.Bodyweight.Unit,
		}
		//log.Printf("[Bodyweight]Projecting %v on %s", trainLog.Bodyweight.Value, trainLog.Timestamp)
		dataPoints = append(dataPoints, point)
	}
	return dataPoints
}

// ProjectExerciseIntensity takes the name of a valid exercise and a list of []common.TrainingLog.
// It returns a DataPoint with the intensity of the highest working set for that exercise per day.
func ProjectExerciseIntensity(exerciseName string, logs []common.TrainingLog) []DataPoint {
	// TODO: can we used a map[time.Time]DataPoint
	// and sort by time before returning
	var dataPoints []DataPoint
	for _, trainLog := range logs {
	Loop1:
		for _, exercise := range trainLog.Workout {
			if exercise.Name == exerciseName {
				point := DataPoint{
					Timestamp: trainLog.Timestamp,
					Value:     exercise.Weight.Value,
					Unit:      exercise.Weight.Unit,
				}
				// If 2 values with same timestamp occurs, we take the highest value
				for _, pts := range dataPoints {
					if pts.Timestamp.Equal(point.Timestamp) && pts.Value > point.Value {
						continue Loop1
					}
				}
				//log.Printf("[Intensity]Projecting %v for %s", exercise, trainLog.Timestamp)
				dataPoints = append(dataPoints, point)
			}
		}
	}
	return dataPoints
}

// ProjectExerciseTonnage takes the name of a valid exercise and []common.TrainingLog.
// It returns a []DataPoint with the total weight of the exercise per day.
// It assumes the exercise has the same unit on the same day.
func ProjectExerciseTonnage(name string, logs []common.TrainingLog) []DataPoint {
	var dataPoints []DataPoint
	var unit string
	for _, trainLog := range logs {
		var projectDay = false
		var tonnagePerDay float64
		for _, exercise := range trainLog.Workout {
			if exercise.Name == name {
				projectDay = true
				tonnagePerDay += (exercise.Weight.Value * float64(exercise.Reps) * float64(exercise.Sets))
				unit = exercise.Weight.Unit
			}
		}
		if projectDay {
			point := DataPoint{
				Timestamp: trainLog.Timestamp,
				Value:     tonnagePerDay,
				Unit:      unit,
			}
			//log.Printf("[Tonnage]Projecting %v for %s", name, trainLog.Timestamp)
			dataPoints = append(dataPoints, point)
		}
	}
	return dataPoints
}

// ProjectExerciseBarLifts takes a valid exercise name and a []common.TrainingLog.
// It returns a []DataPoint of the total number of reps done on the exercise per day.
func ProjectExerciseBarLifts(name string, logs []common.TrainingLog) []DataPoint {
	var dataPoints []DataPoint
	var unit string
	for _, trainLog := range logs {
		var projectDay = false
		var barliftsPerDay = 0
		for _, exercise := range trainLog.Workout {
			if exercise.Name == name {
				projectDay = true
				barliftsPerDay += (exercise.Reps * exercise.Sets)
				unit = exercise.Weight.Unit
			}
		}
		if projectDay {
			point := DataPoint{
				Timestamp: trainLog.Timestamp,
				Value:     float64(barliftsPerDay),
				Unit:      unit,
			}
			//log.Printf("[Barlifts]Projecting %v for %s", name, trainLog.Timestamp)
			dataPoints = append(dataPoints, point)
		}
	}
	return dataPoints
}

// ProjectTrainingDuration takes a []common.TrainingLog.
// It returns a []DataPoint of the training duration per day.
func ProjectTrainingDuration(logs []common.TrainingLog) []DataPoint {
	var dataPoints []DataPoint
	for _, trainLog := range logs {
		point := DataPoint{
			Timestamp: trainLog.Timestamp,
			Value:     float64(trainLog.Duration.Hours()),
			Unit:      "hours",
		}
		//log.Printf("[Training Duration]Projecting %s", trainLog.Timestamp)
		dataPoints = append(dataPoints, point)
	}
	return dataPoints
}

// ProjectExerciseFrequency takes a valid exercise name and a []common.TrainingLog.
// It return a []DataPoint of the frequency of training the exercise.
func ProjectExerciseFrequency(name string, logs []common.TrainingLog) []DataPoint {
	var dataPoints []DataPoint
	for _, trainLog := range logs {
		for _, exercise := range trainLog.Workout {
			if exercise.Name == name {
				point := DataPoint{
					Timestamp: trainLog.Timestamp,
					Value:     1,
				}
				//log.Printf("[Exercise Frequency]Projecting %v for %s", name, trainLog.Timestamp)
				dataPoints = append(dataPoints, point)
			}
		}
	}
	return dataPoints
}
