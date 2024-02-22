package main

import (
	"cmp"
	"slices"
)

// Strip represents the collection of ranges in a strip
// For example, a strip of square with height 2 has a segment [(0,2)]
// If piece has holes, then segment will contain multiple ranges
type Strip []Range

// End returns the maximum end of the segment
func (s Strip) End() float64 {
	if len(s) == 0 {
		return 0
	}

	maxEnd := s[0].End

	for i := 1; i < len(s); i++ {
		if maxEnd < s[i].End {
			maxEnd = s[i].End
		}
	}

	return maxEnd
}

// OccypancyTable represents the part occupancy table
type OccypancyTable []Strip

// Just construct a rectangle
func NewRectanlePart(height, width int) OccypancyTable {
	if width <= 0 || height <= 0 {
		panic("width or height cannot be less or equal to 0")
	}

	segments := make([]Strip, 0, width)
	for i := 0; i < width; i++ {
		segments = append(segments, Strip{
			{
				Start: 0,
				End:   float64(height),
			},
		})
	}
	return segments
}

// TODO: fix discretize for circle
// Discretize decomposes the polygon into a number of vertical strips of the same width
// and the part occupancy is designated by the range of part on each vertical strip
func Discretize(poly Polygon, step float64) OccypancyTable {
	var strips []Strip

	// TODO: optimize
	// TODO: need to return a list of vertices
	findVertexBetween := func(x1, x2 float64) (float64, bool) {
		for _, point := range poly.outerRing {
			if point.X > x1 && point.X < x2 {
				return point.Y, true
			}
		}

		return 0, false
	}

	minx, _, maxx, _ := poly.Bounds()
	l := poly.Intersections(minx)

	for i := minx + step; i < maxx+step; i = i + step {
		r := poly.Intersections(i)

		y1 := slices.MinFunc(append(l.Outer, r.Outer...), func(a, b Point) int {
			return cmp.Compare(a.Y, b.Y)
		})

		y2 := slices.MaxFunc(append(l.Outer, r.Outer...), func(a, b Point) int {
			return cmp.Compare(a.Y, b.Y)
		})

		ymin, ymax := y1.Y, y2.Y

		// a convex figure may have a vertex between intersections
		if vertex, ok := findVertexBetween(float64(i)*step, float64(i+1)*step); ok {
			ymin = min(ymin, vertex)
			ymax = max(ymax, vertex)
		}

		// TODO: handle holes

		strips = append(strips, Strip{{Start: ymin, End: ymax}})

		l = r
	}

	return OccypancyTable(strips)
}
