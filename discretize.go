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

// OccupancyTable represents the part occupancy table
type OccupancyTable []Strip

// Just construct a rectangle
func NewRectanlePart(height, width int) OccupancyTable {
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

// Discretize decomposes the polygon into a number of vertical strips of the same width
// and the part occupancy is designated by the range of part on each vertical strip
func Discretize(poly Polygon, step float64) OccupancyTable {
	var strips []Strip

	minx, _, maxx, _ := poly.Bounds()
	l := poly.Intersections(minx)

	for i := minx + step; i < maxx+step; i = i + step {
		r := poly.Intersections(i)

		outerRange := findOccupancyRange(poly.outerRing, l.Outer, r.Outer, i, step)

		var inners []Range

		for _, inner := range poly.innerRings {
			ri := inner.Intersections(i)
			li := inner.Intersections(i - step)
			if len(ri) == 0 || len(li) == 0 {
				continue
			}
			innerRange := findOccupancyRange(inner, li, ri, i, step)
			inners = append(inners, innerRange)
		}

		strips = append(strips, outerRange.Split(inners))

		l = r
	}

	return OccupancyTable(strips)
}

func findOccupancyRange(ring Ring, l, r []Point, i float64, step float64) Range {
	findVerticesBetween := func(x1, x2 float64) []float64 {
		var vertices []float64
		for _, point := range ring {
			if point.X > x1 && point.X < x2 {
				vertices = append(vertices, point.Y)
			}
		}

		return vertices
	}

	y1 := slices.MinFunc(append(l, r...), func(a, b Point) int {
		return cmp.Compare(a.Y, b.Y)
	})

	y2 := slices.MaxFunc(append(l, r...), func(a, b Point) int {
		return cmp.Compare(a.Y, b.Y)
	})

	ymin, ymax := y1.Y, y2.Y

	// a convex figure may have a vertex between intersections
	vertices := findVerticesBetween(float64(i-step), float64(i))
	if len(vertices) > 0 {
		ymin = min(ymin, slices.Min(vertices))
		ymax = max(ymax, slices.Max(vertices))
	}

	return NewRange(ymin, ymax)
}
