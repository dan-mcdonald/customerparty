package main

import (
	"fmt"
	. "math"
	"strconv"
)

// Coord represents a coordinate on Earth
type Coord struct {
	// φ is the latitude
	// λ is the longitude
	// stored in radians to work natively with trig functions uses Greek letters
	// because lat/lon associated with degrees and it makes the calculations
	// look more like the formula. Users should build Coord objects using
	// the CoordFromLatLon func.
	φ, λ float64
}

// CoordFromLatLon builds a Coord from latitude and longitude specified in
// degrees
func CoordFromLatLon(latitude, longitude string) (*Coord, error) {
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return nil, err
	}
	lon, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return nil, err
	}
	if IsNaN(lat) || Abs(lat) > 90 {
		return nil, fmt.Errorf("Latitude %f is out of range", lat)
	}
	if IsNaN(lon) || Abs(lon) > 180 {
		return nil, fmt.Errorf("Longitude %f is our of range", lon)
	}
	return &Coord{lat / 180.0 * Pi, lon / 180.0 * Pi}, nil
}

// DistanceKm approximates the distance between two points at sea level in
// kilometers. It is accurate to ~0.5%
func DistanceKm(a, b *Coord) float64 {
	// Uses Spherical Law of Cosines formula to approximate distance
	// https://en.wikipedia.org/wiki/Great-circle_distance

	const r = 6371.0088 // i.e. the mean radius of the Earth in kilometers
	Δλ := b.λ - a.λ
	Δσ := Acos(Sin(a.φ)*Sin(b.φ) + Cos(a.φ)*Cos(b.φ)*Cos(Δλ))
	d := r * Δσ

	return d
}
