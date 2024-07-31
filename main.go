package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

const (
	defaultScaleOutput float64 = 0.1
	defaultResolution  float64 = 0.1
)

var defaultAllowedRotations = []int{0}

type intListFlag []int

func (i *intListFlag) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *intListFlag) Set(value string) error {
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*i = append(*i, valueInt)
	return nil
}

var (
	sheetHeight float32 = 200
	maxLength           = int(200 / defaultResolution)
)

var (
	dataset          *string
	scaleOutput      *float64
	resolution       *float64
	allowedRotations intListFlag = defaultAllowedRotations
)

func main() {

	dataset = flag.String("dataset", "datasets/shirts_2007-05-15/shirts.xml", "dataset file")
	scaleOutput = flag.Float64("scale-output", defaultScaleOutput, "scale factor")
	resolution = flag.Float64("resolution", defaultResolution, "resolution")
	flag.Var(&allowedRotations, "rotations", "allowed rotations in degrees")
	flag.Parse()

	println("Loading dataset...")

	f, err := os.Open(*dataset)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var nesting Nesting
	if err := xml.NewDecoder(f).Decode(&nesting); err != nil {
		panic(err)
	}

	polygons := nesting.GetParts()

	for i, poly := range polygons {
		minx, miny, _, _ := poly.Bounds()
		polygons[i] = poly.Offset(NewPoint(-minx, -miny)).Scale(*scaleOutput)
	}

	var sheetWidth float32
	sheetWidth, sheetHeight = nesting.GetBoardSizes()
	var maxLength = int(float64(sheetWidth) * *scaleOutput / *resolution)
	sheetHeight *= float32(*scaleOutput)
	*resolution *= *scaleOutput

	fmt.Println("Dataset loaded")
	fmt.Println("Parts:", len(polygons))
	fmt.Println("Board size:", maxLength, sheetHeight)

	if err := run(polygons); err != nil {
		log.Fatal(err)
	}
}

const (
	angularInterval = 15 // degrees
	angularMin      = 180
	angularMax      = 180

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

func run(figures []Polygon) error {
	var parts []*Part

	var angles []int
	if len(allowedRotations) != 0 {
		angles = allowedRotations
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

		length := calculateSheetLength(parts, *resolution)
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
				occupancy: Discretize(fig, *resolution),
			},
		}
	}
	var orientations []orientation
	for _, i := range angles {
		rotated := fig.Rotate(float64(i))
		orientations = append(orientations, orientation{
			shape:     rotated,
			angle:     float64(i),
			occupancy: Discretize(rotated, *resolution),
		})
	}

	// sort by width
	sort.Slice(orientations, func(i, j int) bool {
		return len(orientations[i].occupancy) < len(orientations[j].occupancy)
	})

	return orientations
}

func calculateSheetLength(parts []*Part, step float64) float32 {
	length := 0.0
	for _, part := range parts {
		xoffset := float64(part.Offset.column) * step
		width := float64(len(part.bestOrienation().occupancy)) * step
		length = max(length, xoffset+width)
	}
	return float32(length)
}

func drawParts(parts []*Part, order []int, file string) error {
	// TODO: fit svg to full screen and fix scroll bar
	svgDrawer := NewSVGDrawer(
		WithOffset(100, -300),
		WithScale(1),
		WithSize(300, 300),
	)

	ordered := make([]*Part, len(order))
	for i, num := range order {
		ordered[i] = parts[num]
	}

	fill := NewBottomLeftFill(sheetHeight, maxLength)
	fill.Run(ordered)

	length := calculateSheetLength(ordered, *resolution)
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

	svgDrawer.DrawCoordSystem(int(length)+25, int(sheetHeight)+25)

	for i, part := range ordered {

		offsetPoint := NewPoint(float64(part.Offset.column)**resolution, part.Offset.y)
		color := fmt.Sprintf("#%02x%02x%02x", randRange(100, 255), randRange(100, 255), randRange(100, 255))
		svgDrawer.AddPart(part.bestOrienation().occupancy, *resolution, offsetPoint,
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
