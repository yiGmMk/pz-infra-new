package geoutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"fmt"
)

var distanceTests = []struct {
	p1       *Point
	p2       *Point
	distance float64
}{
	{&Point{Lat: 22.580812, Lng: 113.961864}, &Point{Lat: 22.580812, Lng: 113.961864}, 0.0},
	{&Point{Lat: 22.580812, Lng: 113.961864}, &Point{Lat: 22.5808, Lng: 113.9619}, 0.003930},
	{&Point{Lat: 22.580812, Lng: 113.961864}, &Point{Lat: 22.6258, Lng: 113.9619}, 5.002},
	{&Point{Lat: 22.580812, Lng: 113.961864}, &Point{Lat: 22.5808, Lng: 114.0106}, 5.004},
	{&Point{Lat: 22.580812, Lng: 113.961864}, &Point{Lat: 22.6198, Lng: 113.9862}, 5.004},
}

func TestDistanceTo(t *testing.T) {
	Convey("Test DestinationPoint", t, func() {
		for _, tt := range distanceTests {
			d := tt.p1.DistanceTo(tt.p2)
			So(d, ShouldAlmostEqual, tt.distance, 0.001)
		}
	})
}

var destinationPointTests = []struct {
	p1        *Point
	direction float64
	distance  float64
	p2        *Point
}{
	{&Point{Lat: 22.580812, Lng: 113.961864}, 0.0, 0.0, &Point{Lat: 22.5808, Lng: 113.9619}},
	{&Point{Lat: 22.580812, Lng: 113.961864}, 90.0, 0.0, &Point{Lat: 22.5808, Lng: 113.9619}},
	{&Point{Lat: 22.580812, Lng: 113.961864}, 0.0, 5.0, &Point{Lat: 22.6258, Lng: 113.9619}},
	{&Point{Lat: 22.580812, Lng: 113.961864}, 90.0, 5.0, &Point{Lat: 22.5808, Lng: 114.0106}},
	{&Point{Lat: 22.580812, Lng: 113.961864}, 30.0, 5.0, &Point{Lat: 22.6198, Lng: 113.9862}},
}

func TestDestinationPoint(t *testing.T) {
	Convey("Test DestinationPoint", t, func() {
		for _, tt := range destinationPointTests {
			p := tt.p1.DestinationPoint(tt.direction, tt.distance)
			So(p.Lat, ShouldAlmostEqual, tt.p2.Lat, 0.0001)
			So(p.Lng, ShouldAlmostEqual, tt.p2.Lng, 0.0001)
		}
	})
}

var directionTests = []struct {
	d float64
	s string
}{
	{10, Northbound},
	{100, Eastbound},
	{190, Southbound},
	{280, Westbound},
}

func TestDirectionNESW(t *testing.T) {
	Convey("Test DirectionNESW", t, func() {
		for _, tt := range directionTests {
			s := DirectionNESW(tt.d)
			So(s, ShouldEqual, tt.s)
		}
	})
}

var samePoiTests = []struct {
	srcLat  float64
	srcLng  float64
	destLat float64
	destLng float64
	isSame  bool
}{
	{srcLat: 22.580812, srcLng: 113.961864, destLat: 22.580812, destLng: 113.961864, isSame: true},
	{srcLat: 22.580812, srcLng: 113.961864, destLat: 22.5808, destLng: 113.9619, isSame: true},
	{srcLat: 22.580812, srcLng: 113.961864, destLat: 22.6258, destLng: 113.9619, isSame: false},
}

func TestIsSamePoi(t *testing.T) {
	Convey("Test Is the Same POI", t, func() {
		for _, test := range samePoiTests {
			So(IsSamePoi(test.srcLat, test.srcLng, test.destLat, test.destLng), ShouldEqual, test.isSame)
		}
	})
}

func TestCalculateDistance(t *testing.T) {
	Convey("Test Calculate Distance", t, func() {
		srcPtLat, srcPtLng := 37.774879,-122.495134
		// destPtLat, destPtLng := 22.5935108,114.1241715  // 0.08360110585897648
		destPtLat, destPtLng := 37.7755882,-122.4956601
		distance := CalculateDistance(srcPtLat, srcPtLng, destPtLat, destPtLng)
		fmt.Println(distance)
		So(distance, ShouldBeLessThan, 0.1)
	})
}