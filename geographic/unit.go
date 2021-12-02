package geo

import (
	"math"
)

type Unit float64

// A Distance represents the length between two points as an int64 meter
type Distance Unit

type Angle Unit

// Length constants
const (
	Meter        Distance = 1e0
	Kilometer             = Meter * 1e3
	Yard                  = Meter * 0.0254 * 12 * 3
	Mile                  = Meter * 1609.344
	NauticalMile          = Meter * 1852
	LightYear             = Meter * 9460730472580800
)

// Angle constants
const (
	Radian    Angle = 1e0
	Degree          = Radian * math.Pi / 180
	Arcminute       = Degree / 60
	Arcsecond       = Degree / 3600
)

// Meters returns the distance in m
func (d Distance) Meters() float64 {
	return float64(d)
}

// Kilometers returns the length in km
func (d Distance) Kilometers() float64 {
	return float64(d / Kilometer)
}

// Yards returns the length in yd
func (d Distance) Yards() float64 {
	return float64(d / Yard)
}

// Miles returns the length in mi
func (d Distance) Miles() float64 {
	return float64(d / Mile)
}

// NauticalMiles returns the length in nm
func (d Distance) NauticalMiles() float64 {
	return float64(d / NauticalMile)
}

// LightYears returns the length in ly
func (d Distance) LightYears() float64 {
	return float64(d / LightYear)
}

// Radians returns the angle in ㎭
func (a Angle) Radians() float64 {
	return float64(a / Radian)
}

// Degrees returns the angle in °
func (a Angle) Degrees() float64 {
	return float64(a / Degree)
}

// Arcminutes returns the angle in amin
func (a Angle) Arcminutes() float64 {
	return float64(a / Arcminute)
}

// Arcseconds returns the angle in asec
func (a Angle) Arcseconds() float64 {
	return float64(a / Arcsecond)
}
