package config

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestViper(t *testing.T) {
	const name = "viper"
	Register(name, &ViperProvider{})

	Convey("Test Providers", t, func() {
		names := Providers()
		So(names, ShouldContain, name)
	})

	Convey("Test GetConfiger", t, func() {
		Convey("option is nil", func() {
			configer, err := GetConfiger(name, nil)
			So(err, ShouldNotBeNil)
			So(configer, ShouldBeNil)
		})

		Convey("option is a wrong value", func() {
			configer, err := GetConfiger(name, "")
			So(err, ShouldNotBeNil)
			So(configer, ShouldBeNil)
		})

	})

	Convey("Test Viper Configer", t, func() {
		configer, err := GetConfiger(name, &ViperOption{
			ConfigFile: "test.yaml",
		})
		So(err, ShouldBeNil)
		So(configer, ShouldNotBeNil)

		Convey("Test Get", func() {
			size := configer.Get("clothing.pants.size")
			So(size, ShouldEqual, "large")

			isHacker := configer.Bool("Hacker")
			So(isHacker, ShouldBeTrue)

			eyes := configer.String("eyes")
			So(eyes, ShouldEqual, "brown")

			age := configer.Int("age")
			So(age, ShouldEqual, 35)

			weight := configer.Float64("weight")
			So(weight, ShouldEqual, 99.9)

			t, _ := time.Parse(time.RFC3339, "1979-05-27T07:32:00Z")
			dob := configer.Time("dob")
			So(dob.Unix(), ShouldEqual, t.Unix())
		})

		Convey("Test Set", func() {

			Convey("Add", func() {
				v := "wow"
				configer.Set("_not_exist_before_", v)
				gv := configer.String("_not_exist_before_")
				So(gv, ShouldEqual, v)
			})

			Convey("Update", func() {
				newAge := 17
				configer.Set("age", 17)
				age := configer.Int("age")
				So(age, ShouldEqual, newAge)
			})
		})

		Convey("Test IsSet", func() {
			isSet := configer.IsSet("_not_exist_forever_")
			So(isSet, ShouldBeFalse)

			isSet = configer.IsSet("Hacker")
			So(isSet, ShouldBeTrue)
		})

	})
}
