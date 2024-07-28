package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const (
	descritizateStep = 5.0
	sheetHeight      = 200
	maxLength        = 1000 / descritizateStep

	populationSize = 20
	elitismRate    = 0.1
	mutationRate   = 0.1
	numGenerations = 50
)

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

	drawParts(figures, rangeSlice(0, len(figures)), "input.svg")
	// drawParts(figures, []int{0, 6, 3, 4, 1, 7, 2, 5}, "input.svg")

	fitnessFn := func(i Individual) float32 {

		parts := make([]OccupancyTable, 0, len(figures))

		for _, num := range i.Order() {
			parts = append(parts, Discretize(figures[num], descritizateStep))
		}

		fill := NewBottomLeftFill(sheetHeight, maxLength)
		offsets := fill.Run(parts)

		length := calculateSheetLength(parts, offsets, descritizateStep)
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

	drawParts(figures, ga.Best().Order(), "output.svg")

	return nil
}

func calculateSheetLength(parts []OccupancyTable, offsets []Offset, step int) float32 {
	length := 0
	for i, offset := range offsets {
		xoffset := offset.column * step
		width := len(parts[i]) * step

		length = max(length, xoffset+width)
	}

	return float32(length)
}

func drawParts(figures []Polygon, order []int, file string) error {
	// TODO: fit svg to full screen and fix scroll bar
	svgDrawer := NewSVGDrawer(
		WithOffset(100, -100),
		WithScale(1),
		WithSize(300, 300),
	)

	var pieces []OccupancyTable

	for _, num := range order {
		pieces = append(pieces, Discretize(figures[num], descritizateStep))
	}

	fill := NewBottomLeftFill(sheetHeight, maxLength)
	offsets := fill.Run(pieces)

	length := calculateSheetLength(pieces, offsets, descritizateStep)
	fmt.Println("Length:", length)

	sheetArea := length * sheetHeight
	fmt.Println("Area:", sheetArea)

	var figuresArea float64
	for _, fig := range figures {
		figuresArea += fig.Area()
	}
	fmt.Println("Free area:", float64(sheetArea)-figuresArea)

	svgDrawer.AddLine(
		float64(length), 0, float64(length), sheetHeight,
		"stroke-width", "2", "stroke-dasharray", "5", "stroke", "blue",
	)
	svgDrawer.AddLine(0, sheetHeight, float64(length), sheetHeight,
		"stroke-width", "2", "stroke-dasharray", "5", "stroke", "blue")

	svgDrawer.DrawCoordSystem(int(length)+25, sheetHeight+25)

	for i, num := range order {
		offset := offsets[i]
		offsetPoint := NewPoint(float64(offset.column*descritizateStep), offset.y)

		piece := pieces[i]
		fig := figures[num]

		// TODO: generate random color
		svgDrawer.AddPart(piece, descritizateStep, offsetPoint,
			"stroke-width", "1", "stroke", "red")
		svgDrawer.AddPolygon(fig.Offset(offsetPoint),
			"stroke-width", "1", "stroke", "black")
		svgDrawer.AddPoint(fig.Centroid().Offset(offsetPoint), "fill", "blue")
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	defer f.Close()
	svgDrawer.Write(f)

	return nil
}
