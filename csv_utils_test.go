package mhist

import (
	"testing"

	"github.com/codeuniversity/ppp-mhist/models"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_constructCsvLine(t *testing.T) {
	Convey("fills the buffer correctly", t, func() {
		m := &models.Numerical{
			Ts:    1000,
			Value: 42,
		}

		byteSlice, err := constructCsvLine(1, m)
		So(err, ShouldBeNil)

		str := string(byteSlice)
		So(str, ShouldEqual, "1,1000,42\n")
	})
}
