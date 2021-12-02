// taken from github.com/kellydunn/golang-geo
package geo

import (
	"math"
)

type Point struct {
	lat float64
	lon float64
}

const (
	EARTH_RADIUS = 6378100 * Meter
)

func NewPoint(lat float64, lon float64) *Point {
	return &Point{lat: lat, lon: lon}
}

func (p *Point) Latitude() float64 {
	return p.lat
}

// Lng returns Point p's longitude.
func (p *Point) Longitude() float64 {
	return p.lon
}

// GreatCircleDistance: Calculates the Haversine distance between two points in Distance unit.
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) GreatCircleDistance(p2 *Point) Distance {
	dLat := Angle(p2.lat-p.lat) * Degree
	dLon := Angle(p2.lon-p.lon) * Degree

	lat1 := Angle(p.lat) * Degree
	lat2 := Angle(p2.lat) * Degree

	a1 := math.Sin(dLat.Radians()/2) * math.Sin(dLat.Radians()/2)
	a2 := math.Sin(dLon.Radians()/2) * math.Sin(dLon.Radians()/2) * math.Cos(lat1.Radians()) * math.Cos(lat2.Radians())

	a := a1 + a2

	c := Distance(2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a)))

	return EARTH_RADIUS * c
}

func (p *Point) PointWithinRadius(p2 *Point, radius Distance) bool {
	return p.GreatCircleDistance(p2) <= radius
}

// BearingTo: Calculates the initial bearing (sometimes referred to as forward azimuth)
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) BearingTo(p2 *Point) Angle {

	dLon := Angle(p2.lon-p.lon) * Degree

	lat1 := Angle(p.lat) * Degree
	lat2 := Angle(p2.lat) * Degree

	y := math.Sin(dLon.Radians()) * math.Cos(lat2.Radians())
	x := math.Cos(lat1.Radians())*math.Sin(lat2.Radians()) -
		math.Sin(lat1.Radians())*math.Cos(lat2.Radians())*math.Cos(dLon.Radians())
	brng := math.Atan2(y, x)

	return Angle(brng)
}

// PointAtDistanceAndBearing returns a Point populated with the lat and lon coordinates
// by transposing the origin point the passed in distance
// by the passed in compass bearing.
// Original Implementation from: http://www.movable-type.co.uk/scripts/latlong.html
func (p *Point) PointAtDistanceAndBearing(dist Distance, bearing Angle) *Point {

	dr := dist / EARTH_RADIUS

	lat1 := Angle(p.lat) * Degree
	lon1 := Angle(p.lon) * Degree

	lat2_part1 := math.Sin(lat1.Radians()) * math.Cos(dr.Meters())
	lat2_part2 := math.Cos(lat1.Radians()) * math.Sin(dr.Meters()) * math.Cos(bearing.Radians())

	lat2 := math.Asin(lat2_part1 + lat2_part2)

	lon2_part1 := math.Sin(bearing.Radians()) * math.Sin(dr.Meters()) * math.Cos(lat1.Radians())
	lon2_part2 := math.Cos(dr.Meters()) - (math.Sin(lat1.Radians()) * math.Sin(lat2))

	lon2 := lon1.Radians() + math.Atan2(lon2_part1, lon2_part2)
	lon2 = math.Mod((lon2+3*math.Pi), (2*math.Pi)) - math.Pi

	lat2 = lat2 * (180.0 / math.Pi)
	lon2 = lon2 * (180.0 / math.Pi)

	return &Point{lat: lat2, lon: lon2}
}
