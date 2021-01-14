package timeutil

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

var (
	displayTimeCases = []struct {
		Seconds int
		Result  string
	}{
		{Seconds: 25, Result: "25sec"},
		{Seconds: 60, Result: "1min"},
		{Seconds: 600, Result: "10min"},
		{Seconds: 3660, Result: "1h 1min"},
		{Seconds: 7200 + 600, Result: "2h 10min"},
	}
)

func TestDisplayForTime(t *testing.T) {
	Convey("Test Display for Time", t, func() {
		for _, testCase := range displayTimeCases {
			So(DisplayForTime(testCase.Seconds, ""), ShouldEqual, testCase.Result)
		}
	})
}

func TestTimeZone(t *testing.T) {
	Convey("Test Time Zone", t, func() {
		loc, err := time.LoadLocation("America/Los_Angeles")
		So(err, ShouldBeNil)
		fmt.Println(UnixToDateTime(CurrentUnix(), loc))
	})
}
