package main

import "math"

// Point represents a 2D point
type Point struct {
	X, Y float64
}

func NewPoint(x, y float64) Point {
	return Point{x, y}
}

// Offset returns a new point that is offset by the given point
func (p Point) Offset(point Point) Point {
	return Point{
		X: p.X + point.X,
		Y: p.Y + point.Y,
	}
}

// Line represents a 2D line
type Line struct {
	Start, End Point
}

// Ring represents a boundary of a 2D polygon. Ring is a closed collection of points
type Ring []Point

// Offset returns a new ring that is offset by the given point
func (r Ring) Offset(point Point) Ring {
	for i := 0; i < len(r); i++ {
		r[i] = r[i].Offset(point)
	}
	return r
}

// Polygon represents a 2D polygon, which is defained by outer and inner rings.
// Polygon has only one outer ring and zero or more inner rings.
type Polygon struct {
	outerRing  Ring
	innerRings []Ring
}

func NewPolygon(outerRing Ring, innerRings ...Ring) Polygon {
	return Polygon{
		outerRing:  outerRing,
		innerRings: innerRings,
	}
}

// Offset returns a new polygon that is offset by the given point
func (p Polygon) Offset(point Point) Polygon {
	var outherRing []Point
	for i := 0; i < len(p.outerRing); i++ {
		outherRing = append(outherRing, p.outerRing[i].Offset(point))
	}

	var innerRings []Ring
	for i := 0; i < len(p.innerRings); i++ {
		var innerRing []Point
		for j := 0; j < len(p.innerRings[i]); j++ {
			innerRing = append(innerRing, p.innerRings[i][j].Offset(point))
		}
		innerRings = append(innerRings, innerRing)
	}

	return Polygon{
		outerRing:  outherRing,
		innerRings: innerRings,
	}
}

type VerticalLine = float64

type Intersection struct {
	Outer []Point
	Inner [][]Point
}

// Intersections returns the intersections of the polygon with the given line
func (p Polygon) Intersections(line VerticalLine) Intersection {

	outer := p.outerRing.Intersections(line)

	var inner [][]Point
	for _, innerRing := range p.innerRings {
		inner = append(inner, innerRing.Intersections(line))
	}

	return Intersection{
		Outer: outer,
		Inner: inner,
	}
}

func (p Polygon) Bounds() (float64, float64, float64, float64) {
	if len(p.outerRing) == 0 {
		return 0, 0, 0, 0
	}

	var (
		minx, miny float64 = p.outerRing[0].X, p.outerRing[0].Y
		maxx, maxy float64 = p.outerRing[0].X, p.outerRing[0].Y
	)
	for _, point := range p.outerRing[1:] {
		minx = min(minx, point.X)
		miny = min(miny, point.Y)

		maxx = max(maxx, point.X)
		maxy = max(maxy, point.Y)
	}

	return minx, miny, maxx, maxy
}

// Intersections returns the intersections of the ring with the given line
func (r Ring) Intersections(line VerticalLine) []Point {
	var intersections []Point

	for i := 0; i < len(r)-1; i++ {
		l := Line{Start: r[i], End: r[i+1]}
		if point, ok := l.Intersect(line); ok {
			intersections = append(intersections, point)
		}
	}

	return intersections
}

// Intersect returns the intersection of the line with the given line
func (l Line) Intersect(line VerticalLine) (Point, bool) {
	if (l.Start.X < line && l.End.X < line) || (l.Start.X > line && l.End.X > line) {
		return Point{}, false
	}

	intersectY := l.y(line)
	if math.IsNaN(intersectY) { // TODO
		return Point{}, false
	}

	return Point{X: line, Y: intersectY}, true
}

// y returns the y coordinate of the line at the given x coordinate
func (l Line) y(x float64) float64 {
	return l.Start.Y + (l.End.Y-l.Start.Y)/(l.End.X-l.Start.X)*(x-l.Start.X)
}

func NewRectangle(x, y, height, width float64) Ring {
	return Ring{{X: x, Y: y},
		{X: x + width, Y: y},
		{X: x + width, Y: y + height},
		{X: x, Y: y + height},
		{X: x, Y: y},
	}
}

func NewCircle(x, y, radius float64, numPoints int) Ring {
	var ring []Point
	for i := 0; i < numPoints; i++ {
		ring = append(ring, Point{
			X: x + radius*math.Cos(2*math.Pi*float64(i)/float64(numPoints)),
			Y: y + radius*math.Sin(2*math.Pi*float64(i)/float64(numPoints)),
		})
	}

	// Close the circle
	ring = append(ring, ring[0])
	return ring
}
