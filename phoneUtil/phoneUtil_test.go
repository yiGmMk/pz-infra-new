package phoneUtil

import (
	"testing"

	"fmt"

	"github.com/yiGmMk/pz-infra-new/tests/base"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	base.InitConfigFile("roav/api/conf/app.conf")
}

func TestValidateMobile(t *testing.T) {
	Convey("Test Validate Email numbers", t, func() {
		Convey("Test Valid Email numbers", func() {
			So(IsValidMobile("13682323959"), ShouldBeTrue)
			So(IsValidMobile("+8613682323959"), ShouldBeTrue)
			So(IsValidMobile("+113682323959"), ShouldBeTrue)
			fmt.Println("Pass testing Valid Email numbers")
		})

		Convey("Test Invalid Email numbers", func() {
			So(IsValidMobile("+86 13-682323959"), ShouldBeFalse)
			So(IsValidMobile("+8613 682323959"), ShouldBeFalse)
			So(IsValidMobile("13-682323959"), ShouldBeFalse)
			So(IsValidMobile("13682323959a"), ShouldBeFalse)
			fmt.Println("Pass testing Invalid Email numbers")
		})

	})
}
