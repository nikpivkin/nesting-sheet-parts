package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestVerticalSegments(t *testing.T) {
// 	tests := []struct {
// 		polygon  Polygon
// 		step     float64
// 		expected [][]Line
// 	}{
// 		{
// 			Polygon{
// 				OuterRing: []Point{{0, 0}, {0, 2}, {2, 2}, {2, 0}},
// 			},
// 			1,
// 			[][]Line{
// 				{{Point{0, 1}, Point{1, 1}}, {Point{1, 1}, Point{2, 1}}, {Point{0, 2}, Point{1, 2}}, {Point{1, 2}, Point{2, 2}}},
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		got := VerticalSegments(tt.polygon, tt.step)
// 		assert.Equal(t, tt.expected, got)
// 	}
// }

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
