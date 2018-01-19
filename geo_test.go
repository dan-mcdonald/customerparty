package main

import (
	"math"
	"testing"
)

type latlon struct {
	lat, lon string
}

func TestCoordFromLatLonValid(t *testing.T) {
	var valids = []struct {
		input    latlon
		expected Coord
	}{
		{latlon{"0", "0"}, Coord{0, 0}},
		{latlon{"12.53675423423", "-43"}, Coord{0.21880763890065386, -0.7504915783575618}},
		{latlon{"-90", "-180"}, Coord{math.Pi / -2, math.Pi * -1}},
		{latlon{"90", "180"}, Coord{math.Pi / 2, math.Pi}},
	}
	for _, testCase := range valids {
		output, err := CoordFromLatLon(testCase.input.lat, testCase.input.lon)
		if err != nil {
			t.Errorf("input %#v generated unexpected err %#v", testCase.input, err)
		}
		if output.φ != testCase.expected.φ {
			t.Errorf("input %#v got phi %#v but expected %#v", testCase.input, output.φ, testCase.expected.φ)
		}
		if output.λ != testCase.expected.λ {
			t.Errorf("input %#v got lambda %#v but expected %#v", testCase.input, output.λ, testCase.expected.λ)
		}
	}
}

func TestCoordFromLatLonInvalid(t *testing.T) {
	var invalids = []latlon{
		// latitude out of range
		{"90.1", "-2"},
		{"-105", "23"},
		{"NaN", "34"},
		{"Inf", "34"},
		{"-Inf", "-13"},
		// longitude out of range
		{"-12", "180.5"},
		{"44", "-200"},
		{"0", "NaN"},
		{"1", "Inf"},
		{"-1", "-Inf"},

		{"123", "-21"},              // lat/lon swapped
		{"1,3", "-21"},              // comma instead of period
		{"NaN", "13"},               // NaN
		{"23N", "32W"},              // cardinal directions instead of sign
		{"sda", "42"},               // stray text
		{"52", "fda"},               // stray text
		{"10°31'42\"", "1°45'20\""}, // DMS format
		{"  10", "54"},              // leading space
		{"0x343", "54"},             // hexadecimal
		{"3", ""},                   // empty string
	}
	for _, testCase := range invalids {
		coord, err := CoordFromLatLon(testCase.lat, testCase.lon)
		if err == nil {
			t.Errorf("Invalid getCoord on %#v should return an err", testCase)
		}
		if coord != nil {
			t.Errorf("Invalid argument %#v to CoordFromLatLon returned %#v for coord instead of nil", testCase, coord)
		}
	}
}

func isClose(a, b float64) bool {
	delta := math.Abs(b - a)
	return delta < 1e-9 || math.Abs(delta/b) < 0.0001
}

func makeCoord(lat, lon string) *Coord {
	coord, err := CoordFromLatLon(lat, lon)
	if err != nil {
		panic(err)
	}
	return coord
}

func TestDistanceMeters(t *testing.T) {
	places := map[string]*Coord{
		"intercom": makeCoord("53.339428", "-6.257664"),
		"pub":      makeCoord("53.3394075", "-6.2584701"),
		"phoenix":  makeCoord("33.3728236", "-112.1059911"),
		"cork":     makeCoord("51.8960528", "-8.4980692"),
	}
	testCases := []struct {
		from, to string
		expected float64
	}{
		// Generated reference values from https://www.movable-type.co.uk/scripts/latlong.html
		{"intercom", "phoenix", 8032},
		{"intercom", "cork", 220.5},
		{"intercom", "pub", 0.05357},
		{"phoenix", "phoenix", 0},
	}
	for _, testCase := range testCases {
		distance := DistanceKm(places[testCase.from], places[testCase.to])
		if !isClose(distance, testCase.expected) {
			t.Errorf("Expected distance between %s and %s is %fm but instead got %fm", testCase.from, testCase.to, testCase.expected, distance)
		}
		distanceBack := DistanceKm(places[testCase.to], places[testCase.from])
		if distance != distanceBack {
			t.Errorf("Expected distance to be commutative")
		}
	}
}
