package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolygon_Offset(t *testing.T) {
	tests := []struct {
		name     string
		poly     Polygon
		offset   Point
		expected Polygon
	}{
		{
			name:     "case 1",
			poly:     NewPolygon(NewRectangle(0, 0, 4, 4)),
			offset:   Point{1, 1},
			expected: NewPolygon(NewRectangle(1, 1, 4, 4)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.poly.Offset(tt.offset)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestPolygon_Intersections(t *testing.T) {
	tests := []struct {
		name     string
		poly     Polygon
		line     VerticalLine
		expected []Point
	}{
		{
			name: "case 1",
			poly: NewPolygon(Ring{
				{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0},
			}, nil),
			line:     1,
			expected: []Point{{1, 2}, {1, 0}},
		},
		{
			name: "case 2",
			poly: NewPolygon(Ring{
				{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0},
			}, nil),
			line:     0,
			expected: []Point{{0, 2}, {0, 0}},
		},
		{
			name: "case 3",
			poly: NewPolygon(
				NewRectangle(0, 0, 4, 4),
				NewRectangle(1, 1, 2, 2),
			),
			line:     2,
			expected: []Point{{2, 4}, {2, 0}, {2, 3}, {2, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.poly.Intersections(tt.line)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestPolygon_Area(t *testing.T) {
	tests := []struct {
		name     string
		poly     Polygon
		expected float64
	}{
		{
			name:     "case 1",
			poly:     NewPolygon(NewRectangle(0, 0, 4, 4)),
			expected: 16,
		},
		{
			name:     "case 2",
			poly:     NewPolygon(NewRectangle(0, 0, 4, 4), NewRectangle(1, 1, 2, 2)),
			expected: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.poly.Area()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestPolygon_Centroid(t *testing.T) {
	tests := []struct {
		name     string
		poly     Polygon
		expected Point
	}{
		{
			name:     "square",
			poly:     NewPolygon(NewRectangle(0, 0, 4, 4)),
			expected: NewPoint(2, 2),
		},
		{
			name:     "square with hole",
			poly:     NewPolygon(NewRectangle(0, 0, 4, 4), NewRectangle(1, 1, 2, 2)),
			expected: NewPoint(2, 2),
		},
		{
			name:     "triangle",
			poly:     NewPolygon(Ring{{0, 0}, {50, 100}, {100, 0}, {0, 0}}),
			expected: NewPoint(50, 50),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.poly.Centroid()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestPolygon_Rotate(t *testing.T) {

	tests := []struct {
		name     string
		poly     Polygon
		angle    float64
		expected Polygon
	}{
		{
			name:  "case 1",
			poly:  NewPolygon(NewRectangle(0, 0, 4, 4)),
			angle: 90,
			expected: NewPolygon(
				Ring{{4, 0}, {0, 0}, {0, 4}, {4, 4}, {4, 0}},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.poly.Rotate(tt.angle)
			assert.Equal(t, tt.expected, got)
		})
	}
}
