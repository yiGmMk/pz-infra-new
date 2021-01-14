package staticMapUtil

import (
	"bytes"
	"image"
	"image/png"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"

	. "github.com/yiGmMk/pz-infra-new/logging"
	"github.com/yiGmMk/pz-infra-new/tests/base"

	. "github.com/smartystreets/goconvey/convey"
)

func prepare() {
	base.InitConfigFile("roav/api/conf/app.conf")
}

func TestGoogleStaticMapClient(t *testing.T) {
	prepare()
	_, filename, _, _ := runtime.Caller(0)
	Log.Debug("currentDir Caller filename is ", With("filename", filename))
	currentCallerDir := path.Dir(filename)
	Log.Debug("currentDir Caller gopath is ", With("currentCallerDir", currentCallerDir))
	Convey("Generate Google Static Map", t, func() {
		client := NewGoogleStaticMapClient()
		//client.SetScale(2)
		client.SetImageSize("400x310")
		client.AddDefaultMarker(StartPointColor, 37.765542, -122.477998)
		client.AddDefaultMarker(EndPointColor, 37.76572, -122.47322)
		client.AddCustomMarker("http://goo.gl/w5O7t4", 37.76567, -122.477098)
		client.AddCustomMarker("http://goo.gl/w5O7t4", 37.76567, -122.476098)
		client.AddCustomMarker("http://goo.gl/w5O7t4", 37.76567, -122.47504)
		client.AddCustomMarker("http://goo.gl/w5O7t4", 37.76567, -122.47464)
		client.AddCustomMarker("http://goo.gl/TgkC8N", 37.765542, -122.476598)
		client.AddCustomMarker("http://goo.gl/TgkC8N", 37.765542, -122.47564)
		coordinates := [][]float64{
			{37.765542, -122.477998},
			{37.76567, -122.47504},
			{37.76567, -122.47464},
			{37.76572, -122.47322},
		}
		client.AddPath("0x24b400AA", 10, coordinates)
		content, err := client.Send()
		So(err, ShouldBeNil)
		So(content, ShouldNotBeNil)
		So(len(content), ShouldBeGreaterThan, 0)
		// convert []byte to image for saving to file
		img, _, _ := image.Decode(bytes.NewReader(content))
		//save the imgByte to file
		out, err := os.Create(currentCallerDir + "/result.png")
		So(err, ShouldBeNil)
		contentType := http.DetectContentType(content)
		Log.Debug("DetectContentType contentType Image  ", With("contentType", contentType))
		So(contentType, ShouldEqual, "image/png")
		err = png.Encode(out, img)
		So(err, ShouldBeNil)
	})
}
