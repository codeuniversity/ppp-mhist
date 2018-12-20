package mhist

import (
	"testing"
	"time"

	"github.com/codeuniversity/ppp-mhist/models"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Passes(t *testing.T) {
	Convey("correct timestamps pass the filter", t, func() {
		definition := models.FilterDefinition{
			Granularity: 2 * time.Millisecond,
			Names:       []string{"bla", "blup"},
		}
		filter := models.NewFilterCollection(definition)
		So(filter.Passes("foo", &models.Numerical{Ts: 1000000}), ShouldBeFalse)
		So(filter.Passes("bla", &models.Numerical{Ts: 1000000}), ShouldBeTrue)
		So(filter.Passes("bla", &models.Numerical{Ts: 2000000}), ShouldBeFalse)
		So(filter.Passes("bla", &models.Numerical{Ts: 3000000}), ShouldBeTrue)
		So(filter.Passes("bla", &models.Numerical{Ts: 4000000}), ShouldBeFalse)
	})
}

func Test_TimestampFilter_Passes(t *testing.T) {
	Convey("correct timestamps pass the filter", t, func() {
		filter := &models.TimestampFilter{Granularity: 2 * time.Millisecond}
		So(filter.Passes(&models.Numerical{Ts: 1000000}), ShouldBeTrue)
		So(filter.Passes(&models.Numerical{Ts: 2000000}), ShouldBeFalse)
		So(filter.Passes(&models.Numerical{Ts: 3000000}), ShouldBeTrue)
		So(filter.Passes(&models.Numerical{Ts: 4000000}), ShouldBeFalse)
	})
}
