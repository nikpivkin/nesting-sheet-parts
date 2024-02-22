package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterval_Difference(t *testing.T) {
	tests := []struct {
		name     string
		interval Range
		others   []Range
		expected []Range
	}{
		{
			name:     "case 1",
			interval: NewRange(0, 10),
			others:   []Range{NewRange(0, 2)},
			expected: []Range{NewRange(2, 10)},
		},
		{
			name:     "case 2",
			interval: NewRange(0, 10),
			others:   []Range{NewRange(2, 4), NewRange(6, 8)},
			expected: []Range{NewRange(0, 2), NewRange(4, 6), NewRange(8, 10)},
		},
		{
			name:     "case 3",
			interval: NewRange(0, 6),
			others:   []Range{NewRange(0, 2), NewRange(4, 6)},
			expected: []Range{NewRange(2, 4)},
		},
		{
			name:     "case 4",
			interval: NewRange(0, 6),
			others:   []Range{NewRange(4, 6)},
			expected: []Range{NewRange(0, 4)},
		},
		{
			name:     "case 5",
			interval: NewRange(0, 4),
			others:   []Range{NewRange(0, 4)},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.interval.Split(tt.others)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewRange(0, 10).Split([]Range{NewRange(0, 2), NewRange(4, 6), NewRange(8, 10)})
	}
}
