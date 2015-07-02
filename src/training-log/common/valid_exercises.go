package common

// List of valid exercise names to be used in training logs.
var validExerciseNames = []string{
	// Squats
	"high bar squats",
	"low bar squats",
	"belted low bar squats",
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
	"paused conventional deadlift",
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

/*
 * Exercise Variations and filtering functions
 */

type ExerciseVariation int

const (
	SquatVariation = iota
	BenchVariation
	DeadliftVariation
)

// Important Exercise variations for exercises we care about
// such as squats, bench and deadlift
var SquatVariations = []string{
	"high bar squats",
	"low bar squats",
	"belted low bar squats",
}
var BenchVariations = []string{
	"close grip bench press",
	"bench press",
	"tng bench press",
}
var DeadliftVariations = []string{
	"sumo deadlift",
	"conventional deadlift",
}

func IsExerciseVariation(variation ExerciseVariation, exerciseName string) bool {
	var variations []string
	switch variation {
	case SquatVariation:
		variations = SquatVariations
	case BenchVariation:
		variations = BenchVariations
	case DeadliftVariation:
		variations = DeadliftVariations
	}
	for i := range variations {
		if variations[i] == exerciseName {
			return true
		}
	}
	return false
}

func Filter(log TrainingLog, filters ...string) []Exercise {
	if len(filters) == 0 {
		return log.Workout
	}
	var result []Exercise
	for i := range log.Workout {
		for j := range filters {

			if log.Workout[i].Name == filters[j] {
				result = append(result, log.Workout[i])
				break
			}
		}
	}
	return result
}

func FilterVariation(variation ExerciseVariation, t TrainingLog) []Exercise {
	var result []Exercise
	for i := range t.Workout {
		if IsExerciseVariation(variation, t.Workout[i].Name) {
			result = append(result, t.Workout[i])
		}
	}
	return result
}
