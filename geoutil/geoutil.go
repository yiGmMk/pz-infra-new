package geoutil

import (
	"fmt"
	"math"
	"strconv"
)

const (
	// According to Wikipedia, the Earth's radius is about 6,371km
	EarthRadius = 6371.0

	Eastbound  = "Eastbound"
	Southbound = "Southbound"
	Westbound  = "Westbound"
	Northbound = "Northbound"

	THRESHOLD_SAME_POI = 0.1 // The same POI threshold, kilo meters
)

type Point struct {
	Lat float64
	Lng float64
}

func NewPoint(lat, lng float64) (*Point, error) {
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("geoutil: illegal latitude[%f] (latitude ranges from -90  to 90)", lat)
	}
	if lng < -180 || lng > 180 {
		return nil, fmt.Errorf("geoutil: illegal longitude[%f] (longitude ranges from -180 to 180)", lng)

	}
	return &Point{Lat: lat, Lng: lng}, nil
}

func (p *Point) DistanceTo(p2 *Point) float64 {
	lat := toRadians(p.Lat)
	lng := toRadians(p.Lng)
	lat2 := toRadians(p2.Lat)
	lng2 := toRadians(p2.Lng)
	l := lat2 - lat
	g := lng2 - lng
	a := math.Sin(l/2)*math.Sin(l/2) +
		math.Cos(lat)*math.Cos(lat2)*math.Sin(g/2)*math.Sin(g/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return EarthRadius * c
}

func (p *Point) DestinationPoint(direction, distance float64) *Point {
	r := distance / EarthRadius
	d := toRadians(direction)
	lat1 := toRadians(p.Lat)
	lng1 := toRadians(p.Lng)

	lat2 := math.Asin(math.Sin(lat1)*math.Cos(r) +
		math.Cos(lat1)*math.Sin(r)*math.Cos(d))
	lng2 := lng1 + math.Atan2(math.Sin(d)*math.Sin(r)*math.Cos(lat1),
		math.Cos(r)-math.Sin(lat1)*math.Sin(lat2))
	return &Point{Lat: toDegrees(lat2), Lng: math.Mod(toDegrees(lng2)+540, 360) - 180}
}

func toRadians(f float64) float64 {
	return f * math.Pi / 180
}

func toDegrees(f float64) float64 {
	return f * 180 / math.Pi
}

func DirectionNESW(direction float64) string {
	direction = modDirection(direction)
	switch {
	case direction > 45 && direction <= 135:
		return Eastbound
	case direction > 135 && direction <= 225:
		return Southbound
	case direction > 225 && direction <= 315:
		return Westbound
	default:
		return Northbound
	}
}

func DirectionNS(direction float64) string {
	direction = modDirection(direction)
	switch {
	case direction > 90 && direction <= 270:
		return Southbound
	default:
		return Northbound
	}
}

func DirectionWE(direction float64) string {
	direction = modDirection(direction)
	switch {
	case direction > 0 && direction <= 180:
		return Eastbound
	default:
		return Westbound
	}
}

func modDirection(direction float64) float64 {
	direction = math.Mod(direction, 360)
	if direction < 0 {
		direction = direction + 360
	}
	return direction
}

// Check if two GPS Point could treat as the same POI
func IsSamePoi(srcLat, srcLng, destLat, destLng float64) bool {
	p1, err := NewPoint(srcLat, srcLng)
	if err != nil {
		return false
	}
	p2, err := NewPoint(destLat, destLng)
	if err != nil {
		return false
	}
	if p1.DistanceTo(p2) <= THRESHOLD_SAME_POI {
		return true
	}
	return false
}

// Calculate two GPS Point Distance
func CalculateDistance(srcLat, srcLng, destLat, destLng float64) float64 {
	p1, err := NewPoint(srcLat, srcLng)
	if err != nil {
		return EarthRadius
	}
	p2, err := NewPoint(destLat, destLng)
	if err != nil {
		return EarthRadius
	}
	return p1.DistanceTo(p2)
}

func FormatLatLng(lat, lng float64) string {
	return fmt.Sprintf("%s,%s", strconv.FormatFloat(lat, 'f', 8, 64), strconv.FormatFloat(lng, 'f', 8, 64))
}
