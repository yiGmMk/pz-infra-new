package timezoneUtil

import (
	"fmt"
	"testing"

	"github.com/gyf841010/pz-infra-new/tests/base"

	. "github.com/smartystreets/goconvey/convey"
)

func prepare() {
	base.InitConfigFile("roav/api/conf/app.conf")
}

func TestGoogleTimezoneClient(t *testing.T) {
	prepare()
	Convey("Search Google Timezone", t, func() {
		Convey("Test CN Timezone", func() {
			cnLat, cnLng := 22.5617329, 113.9527458
			client := NewGoogleTimezoneClient()
			client.SetLocation(cnLat, cnLng)
			timezoneInfo, err := client.Send()
			So(err, ShouldBeNil)
			So(timezoneInfo, ShouldNotBeNil)
			fmt.Println("Got Timezone Info", timezoneInfo)
			So(timezoneInfo.TimeZoneId, ShouldEqual, "Asia/Shanghai")
			So(timezoneInfo.RawOffset, ShouldEqual, 28800)
		})
		Convey("Test US Timezone", func() {
			usLat, usLng := 40.6475468, -122.3386514
			client := NewGoogleTimezoneClient()
			client.SetLocation(usLat, usLng)
			timezoneInfo, err := client.Send()
			So(err, ShouldBeNil)
			So(timezoneInfo, ShouldNotBeNil)
			fmt.Println("Got Timezone Info", timezoneInfo)
			So(timezoneInfo.TimeZoneId, ShouldEqual, "America/Los_Angeles")
			So(timezoneInfo.RawOffset, ShouldEqual, -28800)
		})
		Convey("Test Invalid Timezone", func() {
			invalidLat, invalidLng := 0.0, 0.0
			client := NewGoogleTimezoneClient()
			client.SetLocation(invalidLat, invalidLng)
			timezoneInfo, err := client.Send()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "No Timezone Info Found")
			So(timezoneInfo, ShouldBeNil)
		})
	})
}
