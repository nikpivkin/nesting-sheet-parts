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

// Area returns the area of the ring
// https://en.wikipedia.org/wiki/Shoelace_formula#Triangle_formula
func (r Ring) Area() float64 {
	var area float64
	for i := 0; i < len(r)-1; i++ {
		area += r[i].X * r[(i+1)].Y
		area -= r[i].Y * r[(i+1)].X
	}
	// the sign of the area changes because the points are numbered clockwise.
	return area / -2
}

// Centroid returns the centroid of the ring
// https://en.wikipedia.org/wiki/Centroid#Of_a_polygon
func (r Ring) Centroid() Point {
	area := -r.Area()
	var x, y float64
	for i := 0; i < len(r)-1; i++ {
		x += (r[i].X + r[i+1].X) * (r[i].X*r[i+1].Y - r[i+1].X*r[i].Y)
		y += (r[i].Y + r[i+1].Y) * (r[i].X*r[i+1].Y - r[i+1].X*r[i].Y)
	}

	x /= (6 * area)
	y /= (6 * area)
	return NewPoint(x, y)
}

// Rotate returns a new ring that is rotated by the given angle in degrees
// around the given point.
// https://en.wikipedia.org/wiki/Rotation_(mathematics)#Two_dimensions
func (r Ring) Rotate(angle float64, center Point) Ring {
	radians := angle * math.Pi / 180.0
	rotated := make(Ring, len(r))

	for i, point := range r {
		rotatedPoint := NewPoint(
			(point.X-center.X)*math.Cos(radians)-(point.Y-center.Y)*math.Sin(radians)+center.X,
			(point.X-center.X)*math.Sin(radians)+(point.Y-center.Y)*math.Cos(radians)+center.Y,
		)
		rotated[i] = rotatedPoint
	}

	return rotated
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

// Area returns the area of the polygon
func (p Polygon) Area() float64 {
	area := p.outerRing.Area()
	for _, r := range p.innerRings {
		area -= r.Area()
	}
	return area
}

// Centroid returns the centroid of the polygon
func (p Polygon) Centroid() Point {
	if len(p.innerRings) == 0 {
		return p.outerRing.Centroid()
	}

	outerArea := p.outerRing.Area()
	totalArea := outerArea

	var x, y float64
	for i := 0; i < len(p.innerRings); i++ {
		innerArea := p.innerRings[i].Area()
		totalArea -= innerArea

		centroid := p.innerRings[i].Centroid()
		x -= centroid.X * innerArea
		y -= centroid.Y * innerArea
	}

	centroid := p.outerRing.Centroid()
	x += centroid.X * outerArea
	y += centroid.Y * outerArea

	return NewPoint(x/totalArea, y/totalArea)
}

// Rotate returns a new polygon rotated by the given angle in degrees
func (p Polygon) Rotate(angle float64) Polygon {
	inners := make([]Ring, len(p.innerRings))
	for i, innerRing := range p.innerRings {
		inners[i] = innerRing.Rotate(angle, p.outerRing.Centroid())
	}
	outer := p.outerRing.Rotate(angle, p.outerRing.Centroid())
	return NewPolygon(outer, inners...)
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
		{X: x, Y: y + height},
		{X: x + width, Y: y + height},
		{X: x + width, Y: y},
		{X: x, Y: y},
	}
}

func NewCircle(x, y, radius float64, numPoints int) Ring {
	points := make(Ring, numPoints+1)
	angleStep := 2 * math.Pi / float64(numPoints)

	for i := 0; i < numPoints; i++ {
		angle := angleStep * float64(i)
		points[i] = Point{
			X: x + radius*math.Cos(angle+math.Pi),
			Y: y + radius*math.Sin(angle+math.Pi),
		}
	}

	points[numPoints] = points[0]
	return points
}
