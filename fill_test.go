package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {

	tests := []struct {
		name     string
		ymax     int
		pieces   []OccypancyTable
		expected []Offset
	}{
		{
			name: "case 1",
			ymax: 10,
			/*
				[ ]
				[ ]
			*/
			pieces: []OccypancyTable{
				NewRectanlePart(2, 2),
				NewRectanlePart(2, 2),
			},
			expected: []Offset{{0, 0}, {0, 2}},
		},
		{
			name: "case 2",
			ymax: 2,
			/*
				[ ][ ]
			*/
			pieces: []OccypancyTable{
				NewRectanlePart(2, 2),
				NewRectanlePart(2, 2),
			},
			expected: []Offset{{0, 0}, {2, 0}},
		},
		{
			name: "case 3",
			ymax: 4,
			/*
				[ ]
				[ ][ ]
			*/
			pieces: []OccypancyTable{
				NewRectanlePart(2, 2),
				NewRectanlePart(2, 2),
				NewRectanlePart(2, 2),
			},
			expected: []Offset{{0, 0}, {0, 2}, {2, 0}},
		},
		{
			name: "case 4",
			ymax: 4,
			/*
				[ ][ ]
				[    ]
			*/
			pieces: []OccypancyTable{
				NewRectanlePart(2, 4),
				NewRectanlePart(2, 2),
				NewRectanlePart(2, 2),
			},
			expected: []Offset{{0, 0}, {0, 2}, {2, 2}},
		},
		{
			name: "case 5",
			ymax: 4,
			/*
				|   ]
				| |[ ]
			*/
			pieces: []OccypancyTable{
				{
					{
						{Start: 0, End: 4},
					},
					{
						{Start: 0, End: 4},
					},
					{
						{Start: 2, End: 4},
					},
					{
						{Start: 2, End: 4},
					},
				},
				NewRectanlePart(2, 2),
			},
			expected: []Offset{{0, 0}, {2, 0}},
		},
		{
			name: "case 6",
			ymax: 6,
			/*
				|x|
				|o|
				|x|
			*/
			pieces: []OccypancyTable{
				{
					{
						{Start: 0, End: 2},
						{Start: 4, End: 6},
					},
				},
				{
					{
						{Start: 0, End: 2},
					},
				},
			},
			expected: []Offset{
				{0, 0}, {0, 2},
			},
		},
		{
			name: "case 7",
			ymax: 6,
			/*
				|x  x|
				   |x|
			*/
			pieces: []OccypancyTable{
				{
					{
						{Start: 2, End: 4},
					},
					{
						{Start: 0, End: 4},
					},
				},
			},
			expected: []Offset{
				{0, 0},
			},
		},
		{
			name: "case 8",
			ymax: 6,
			/*
				|x    x|
				|o || x|
			*/
			pieces: []OccypancyTable{
				{
					{
						{Start: 2, End: 4},
					},
					{
						{Start: 0, End: 4},
					},
				},
				{
					{
						{Start: 0, End: 2},
					},
				},
			},
			expected: []Offset{
				{0, 0}, {0, 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			algo := NewBottomLeftFill(tt.ymax, 10)
			got := algo.Run(tt.pieces)
			assert.Equal(t, tt.expected, got)
		})
	}
}
