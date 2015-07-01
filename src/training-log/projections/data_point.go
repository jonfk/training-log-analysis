package projections

import (
	"time"
)

type DataPoint struct {
	Timestamp time.Time
	Value     float64
	Unit      string
}
