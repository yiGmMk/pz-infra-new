package commonUtil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"fmt"
)

func TestDeleteSliceIndex(t *testing.T) {
	Convey("Test Delete Slice with Index", t, func() {
		Convey("Index is 0", func() {
			testSlice := []string{"a", "b", "c"}
			result := DeleteSliceIndex(testSlice, 0).([]string)
			fmt.Println(result)
			So(len(result), ShouldEqual, 2)
			So(result[0], ShouldEqual, "b")
			So(result[1], ShouldEqual, "c")
		})

		Convey("Index is last one", func() {
			testSlice := []string{"a", "b", "c"}
			result := DeleteSliceIndex(testSlice, 2).([]string)
			fmt.Println(result)
			So(len(result), ShouldEqual, 2)
			So(result[0], ShouldEqual, "a")
			So(result[1], ShouldEqual, "b")
		})

		Convey("Index is intermediate one", func() {
			testSlice := []string{"a", "b", "c"}
			result := DeleteSliceIndex(testSlice, 1).([]string)
			fmt.Println(result)
			So(len(result), ShouldEqual, 2)
			So(result[0], ShouldEqual, "a")
			So(result[1], ShouldEqual, "c")
		})

		Convey("Index out of bould should do nothing", func() {
			testSlice := []string{"a", "b", "c"}
			result := DeleteSliceIndex(testSlice, 4).([]string)
			fmt.Println(result)
			So(len(result), ShouldEqual, 3)
			So(result[0], ShouldEqual, "a")
			So(result[1], ShouldEqual, "b")
			So(result[2], ShouldEqual, "c")
		})
	})
}
