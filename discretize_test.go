package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescritizate(t *testing.T) {
	tests := []struct {
		name     string
		poly     Polygon
		step     float64
		expected OccupancyTable
	}{
		{
			name:     "square",
			poly:     NewPolygon(Ring{{0, 0}, {0, 1}, {2, 1}, {2, 0}, {0, 0}}, nil),
			step:     1,
			expected: OccupancyTable{{{0, 1}}, {{0, 1}}},
		},
		{
			name: "square with inner ring",
			poly: NewPolygon(
				Ring{{0, 0}, {0, 10}, {10, 10}, {10, 0}, {0, 0}},
				Ring{{2.5, 2.5}, {2.5, 7.5}, {7.5, 7.5}, {7.5, 2.5}, {2.5, 2.5}},
			),
			step:     2.5,
			expected: OccupancyTable{{{0, 10}}, {{0, 2.5}, {7.5, 10}}, {{0, 2.5}, {7.5, 10}}, {{0, 10}}},
		},
		{
			name:     "rhombus 1",
			poly:     NewPolygon(Ring{{0, 2}, {2, 4}, {4, 2}, {2, 0}, {0, 2}}, nil),
			step:     2,
			expected: OccupancyTable{{{0, 4}}, {{0, 4}}},
		},
		{
			name:     "triangle 1",
			poly:     NewPolygon(Ring{{0, 0}, {1, 3}, {4, 0}, {0, 0}}, nil),
			step:     2,
			expected: OccupancyTable{{{0, 3}}, {{0, 2}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Discretize(tt.poly, tt.step)
			assert.Equal(t, tt.expected, got)
		})
	}
}
