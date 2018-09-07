package test_helpers

import (
	"github.com/codeuniversity/ppp-mhist"
)

//AddMeasurementsToSeries sample to series
func AddMeasurementsToSeries(s *mhist.Series) {
	measurements := GetSampleMeasurements(5, 1000, 10)
	for _, m := range measurements {
		s.Add(m)
	}
}

//GetSampleMeasurements ...
func GetSampleMeasurements(amount, startTs, increment int64) []*mhist.Measurement {
	measurements := []*mhist.Measurement{}
	for i := int64(0); i < amount; i++ {
		measurements = append(measurements, &mhist.Measurement{Ts: startTs + increment*i, Value: float64(10 + i)})
	}
	return measurements
}
