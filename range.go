package main

import (
	"fmt"
	"sort"
)

// Range represents the vacancy on a strip
type Range struct {
	Start, End float64
}

// NewRange returns a new range
func NewRange(start, end float64) Range {
	return Range{Start: start, End: end}
}

const epsilon = 0.000001

// Includes returns true if the range includes the other
func (s Range) Includes(other Range) bool {
	return s.Start-epsilon <= other.Start && s.End+epsilon >= other.End
}

// Overlaps returns true if the range overlaps with the other
func (s Range) Overlaps(other Range) bool {
	return s.Includes(other) || other.Includes(s)
}

// Length returns the length of the range
func (s Range) Length() float64 {
	return s.End - s.Start
}

// Add adds a value to the range
func (s Range) Add(val float64) Range {
	return Range{
		Start: toFixed(s.Start+val, 4),
		End:   toFixed(s.End+val, 4),
	}
}

// Offset returns a new range that is offset by the given point
func (s Range) Offset(point Point) Range {
	return Range{
		Start: toFixed(s.Start+point.X, 4),
		End:   toFixed(s.End+point.X, 4),
	}
}

// Split splits the range into multiple ranges
// For example, if the range is (0, 10) and the given ranges are
// [(0, 2), (4, 6), (8, 10)], the result will be [(2, 4), (6, 8)]
func (i Range) Split(others []Range) []Range {
	if len(others) == 0 {
		return []Range{i}
	}

	for j := 0; j < len(others)-1; j++ {
		if others[j].Overlaps(others[j+1]) {
			panic("range.Split: overlapping ranges")
		}
	}

	for _, other := range others {
		if !i.Includes(other) {
			panic(fmt.Sprintf("range.Split: range %v must include all other ranges, but %v does not", i, other))
		}
	}

	sort.Slice(others, func(j, k int) bool {
		return others[j].Start < others[k].Start
	})

	var result []Range

	start := i.Start
	for _, other := range others {
		if start != other.Start {
			result = append(result, NewRange(start, other.Start))
		}
		start = other.End
	}

	if start < i.End {
		result = append(result, NewRange(start, i.End))
	}

	return result
}
