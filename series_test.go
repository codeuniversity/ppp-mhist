package mhist_test

import (
	"testing"

	"github.com/codeuniversity/ppp-mhist"
	"github.com/codeuniversity/ppp-mhist/testhelpers"
	. "github.com/smartystreets/goconvey/convey"
)

const maxSize = 100 * 1024 * 1024

func TestSeries(t *testing.T) {
	emptyFilterDefinition := mhist.FilterDefinition{}
	Convey("Series", t, func() {
		Convey("Add()", func() {
			Convey("It only adds measurements it was created with", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				s.Add(&mhist.Numerical{Ts: 1000})
				s.Add(&mhist.Categorical{Ts: 2000})
				returnedMeasurements, _ := s.GetMeasurementsInTimeRange(0, 3000, emptyFilterDefinition)
				So(len(returnedMeasurements), ShouldEqual, 1)
			})
		})
		Convey("GetMeasurementsInTimeRange()", func() {
			Convey("returns no measurements if empty", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				returnedMeasurements, _ := s.GetMeasurementsInTimeRange(1005, 1035, emptyFilterDefinition)
				s.Shutdown()

				So(len(returnedMeasurements), ShouldEqual, 0)
			})
			Convey("returns correct measurements if given range is inside", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				testhelpers.AddMeasurementsToSeries(s)
				returnedMeasurements, _ := s.GetMeasurementsInTimeRange(1005, 1035, emptyFilterDefinition)

				s.Shutdown()
				So(len(returnedMeasurements), ShouldEqual, 3)
			})
			Convey("returns all measurements if it is completly inside given range", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				testhelpers.AddMeasurementsToSeries(s)
				returnedMeasurements, _ := s.GetMeasurementsInTimeRange(500, 4000, emptyFilterDefinition)

				s.Shutdown()
				So(len(returnedMeasurements), ShouldEqual, 5)
			})

			Convey("returns no measurements if given range has no overlap", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				testhelpers.AddMeasurementsToSeries(s)
				returnedMeasurements, _ := s.GetMeasurementsInTimeRange(3000, 4000, emptyFilterDefinition)

				s.Shutdown()
				So(len(returnedMeasurements), ShouldEqual, 0)
			})

			Convey("returns correct if given range has partialy overlaps", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				testhelpers.AddMeasurementsToSeries(s)
				returnedMeasurements, _ := s.GetMeasurementsInTimeRange(1025, 4000, emptyFilterDefinition)

				s.Shutdown()
				So(len(returnedMeasurements), ShouldEqual, 2)
			})
			Convey("returns incomplete = true if start Ts is below lowest measurement in series", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				testhelpers.AddMeasurementsToSeries(s)
				_, incomplete := s.GetMeasurementsInTimeRange(0, 4000, emptyFilterDefinition)

				s.Shutdown()
				So(incomplete, ShouldEqual, true)
			})
		})

		Convey("CutoffBelow()", func() {
			Convey("returns correct measurements", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				testhelpers.AddMeasurementsToSeries(s)

				So(s.Size(), ShouldEqual, 80)
				returnedMeasurements := s.CutoffBelow(1025)
				So(len(returnedMeasurements), ShouldEqual, 3)
				So(s.Size(), ShouldEqual, 32)
				s.Shutdown()
			})

			Convey("returns no measurements if timestamp is below all of series", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				testhelpers.AddMeasurementsToSeries(s)

				So(s.Size(), ShouldEqual, 80)
				returnedMeasurements := s.CutoffBelow(900)
				So(len(returnedMeasurements), ShouldEqual, 0)
				So(s.Size(), ShouldEqual, 80)
				s.Shutdown()
			})

			Convey("returns all measurements if timestamp is above all of series", func() {
				s := mhist.NewSeries(mhist.MeasurementNumerical)
				testhelpers.AddMeasurementsToSeries(s)

				So(s.Size(), ShouldEqual, 80)
				returnedMeasurements := s.CutoffBelow(2000)
				So(len(returnedMeasurements), ShouldEqual, 5)
				So(s.Size(), ShouldEqual, 0)
				s.Shutdown()
			})
		})
	})
}
