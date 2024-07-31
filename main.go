package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
)

var sheetHeight = 200
var maxLength = 1000 / descritizateStep
var availableAngles = []int{0, 180}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const (
	descritizateStep = 5
	angularInterval  = 15 // degrees
	angularMin       = 180
	angularMax       = 180

	populationSize = 20
	elitismRate    = 0.1
	mutationRate   = 0.2
	numGenerations = 50
)

type Part struct {
	Orientations       []orientation
	Shape              Polygon
	Offset             Offset
	BestOrientationNum int
}

func (p Part) bestOrienation() orientation {
	return p.Orientations[p.BestOrientationNum]
}

func run() error {
	squareWithHole := NewPolygon(NewRectangle(0, 0, 100, 100), NewRectangle(25, 25, 50, 50))
	squareWithHole2 := NewPolygon(NewRectangle(0, 0, 100, 100), Ring{{75, 75}, {80, 80}, {85, 75}, {75, 75}})
	triangle := NewPolygon(Ring{{0, 0}, {50, 100}, {100, 0}, {0, 0}})

	smallSquare := NewPolygon(NewRectangle(0, 0, 10, 10))
	smallSquare2 := NewPolygon(NewRectangle(0, 0, 30, 30))
	smallSquare3 := NewPolygon(NewRectangle(0, 0, 50, 50))
	smallSquare4 := NewPolygon(NewRectangle(0, 0, 70, 70))
	smallSquare5 := NewPolygon(NewRectangle(0, 0, 90, 90))

	circle := NewPolygon(NewCircle(50, 50, 50, 16))

	figures := []Polygon{
		squareWithHole,
		squareWithHole2,
		triangle, triangle,
		smallSquare, smallSquare2, smallSquare3, smallSquare4, smallSquare5,
		triangle, triangle,
		smallSquare, smallSquare2, smallSquare3, smallSquare4, smallSquare5,
		squareWithHole, squareWithHole,
		circle,
	}

	var parts []*Part

	var angles []int
	if len(availableAngles) != 0 {
		angles = availableAngles
	} else {
		angles = rangeSlice(angularMin, angularMax, angularInterval)
	}

	for _, fig := range figures {
		parts = append(parts, &Part{
			Orientations: createOrientations(fig, angles...),
			Shape:        fig,
		})
	}

	drawParts(parts, rangeSlice(0, len(figures), 1), "input.svg")

	fitnessFn := func(i Individual) float32 {
		order := i.Order()
		ordered := make([]*Part, len(order))
		for i, num := range order {
			ordered[i] = parts[num]
		}

		fill := NewBottomLeftFill(sheetHeight, maxLength)
		fill.Run(ordered)

		length := calculateSheetLength(parts, descritizateStep)
		return -length
	}

	ga := NewGeneticAlgorithm(
		len(figures),
		fitnessFn,
		WithPopulationSize(populationSize),
		WithElitismRate(elitismRate),
		WithMutationRate(mutationRate),
	)

	ga.Run(numGenerations)
	fmt.Printf("Best fitness: %f, Order: %v\n", ga.Best().Fitness(), ga.Best().Order())

	drawParts(parts, ga.Best().Order(), "output.svg")

	return nil
}

type orientation struct {
	shape     Polygon
	occupancy OccupancyTable
	angle     float64
}

func createOrientations(fig Polygon, angles ...int) []orientation {
	if len(angles) == 0 {
		return []orientation{
			{
				shape:     fig,
				angle:     0,
				occupancy: Discretize(fig, descritizateStep),
			},
		}
	}
	var orientations []orientation
	for _, i := range angles {
		rotated := fig.Rotate(float64(i))
		orientations = append(orientations, orientation{
			shape:     rotated,
			angle:     float64(i),
			occupancy: Discretize(rotated, descritizateStep),
		})
	}

	// sort by width
	sort.Slice(orientations, func(i, j int) bool {
		return len(orientations[i].occupancy) < len(orientations[j].occupancy)
	})

	return orientations
}

func calculateSheetLength(parts []*Part, step int) float32 {
	length := 0
	for _, part := range parts {
		xoffset := part.Offset.column * step
		width := len(part.bestOrienation().occupancy) * step
		length = max(length, xoffset+width)
	}
	return float32(length)
}

func drawParts(parts []*Part, order []int, file string) error {
	// TODO: fit svg to full screen and fix scroll bar
	svgDrawer := NewSVGDrawer(
		WithOffset(100, -100),
		WithScale(1),
		WithSize(300, 300),
	)

	ordered := make([]*Part, len(order))
	for i, num := range order {
		ordered[i] = parts[num]
	}

	fill := NewBottomLeftFill(sheetHeight, maxLength)
	fill.Run(ordered)

	length := calculateSheetLength(ordered, descritizateStep)
	fmt.Println("Length:", length)

	sheetArea := length * float32(sheetHeight)
	fmt.Println("Area:", sheetArea)

	var figuresArea float64
	for _, part := range ordered {
		figuresArea += part.Shape.Area()
	}
	fmt.Println("Free area:", float64(sheetArea)-figuresArea)

	svgDrawer.AddLine(
		float64(length), 0, float64(length), float64(sheetHeight),
		"stroke-width", "2", "stroke-dasharray", "5", "stroke", "blue",
	)
	svgDrawer.AddLine(0, float64(sheetHeight), float64(length), float64(sheetHeight),
		"stroke-width", "2", "stroke-dasharray", "5", "stroke", "blue")

	svgDrawer.DrawCoordSystem(int(length)+25, sheetHeight+25)

	for i, part := range ordered {
		offsetPoint := NewPoint(float64(part.Offset.column*descritizateStep), part.Offset.y)

		color := fmt.Sprintf("#%02x%02x%02x", randRange(100, 255), randRange(100, 255), randRange(100, 255))
		svgDrawer.AddPart(part.bestOrienation().occupancy, descritizateStep, offsetPoint,
			"stroke-width", "1", "stroke", color)
		svgDrawer.AddPolygon(part.bestOrienation().shape.Offset(offsetPoint),
			"stroke-width", "1", "stroke", "black")

		center := part.bestOrienation().shape.Centroid().Offset(offsetPoint)
		svgDrawer.AddPoint(center, "fill", "blue", "r", "1.5")
		// TODO: draw text over figures
		svgDrawer.AddText(center.Offset(NewPoint(2, 2)), fmt.Sprintf("%d", i), "font-size", "4")
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer f.Close()
	svgDrawer.Write(f)

	return nil
}

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}
