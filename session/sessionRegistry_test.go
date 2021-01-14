package session

import (
	"testing"
	"time"

	"github.com/gyf841010/pz-infra-new/commonUtil"
	"github.com/gyf841010/pz-infra-new/tests/base"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	TEST_USER_ID   = "TEST_USER_ID_SESSION_REGISTRY"
	TEST_CLIENT_ID = "TEST_CLIENT_ID"
	testRegistry   = Registry()
)

func initTest() {
	base.InitConfigFile("roav/api/conf/app.conf")
}

func TestGetSessionRegistry(t *testing.T) {
	initTest()
	Convey("test Get SessionRegistry", t, func() {
		Convey("Previous Not Exist User session", func() {
			//SESSION_REGISTRY_KEY_PREFIX+TEST_USER_ID
			sessionId, clientId, err := testRegistry.GetUserSession(TEST_USER_ID)
			So(err, ShouldBeNil)
			So(sessionId, ShouldBeEmpty)
			So(clientId, ShouldBeEmpty)
		})

		Convey("Exist User session", func() {
			mockSessionId := commonUtil.UUID()
			err := testRegistry.SessionRegistry(TEST_USER_ID, TEST_CLIENT_ID, mockSessionId, 5)
			So(err, ShouldBeNil)
			sessionId, clientId, err := testRegistry.GetUserSession(TEST_USER_ID)
			So(err, ShouldBeNil)
			So(sessionId, ShouldEqual, mockSessionId)
			So(clientId, ShouldEqual, TEST_CLIENT_ID)
		})

		Convey("Expired User session", func() {
			mockSessionId := commonUtil.UUID()
			err := testRegistry.SessionRegistry(TEST_USER_ID, TEST_CLIENT_ID, mockSessionId, 3)
			So(err, ShouldBeNil)
			time.Sleep(time.Duration(3+1) * time.Second)
			sessionId, clientId, err := testRegistry.GetUserSession(TEST_USER_ID)
			So(err, ShouldBeNil)
			So(sessionId, ShouldBeEmpty)
			So(clientId, ShouldBeEmpty)
		})
	})
}
