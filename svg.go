package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type SVGDrawer struct {
	buffer bytes.Buffer
	offset Point   // offset in pixels
	scale  float64 // scale factor
	size   [2]int  // width, height in pixels
}

type SVGDrawerOption func(*SVGDrawer)

func WithOffset(x, y float64) SVGDrawerOption {
	return func(d *SVGDrawer) {
		d.offset = Point{x, y}
	}
}

func WithScale(scale float64) SVGDrawerOption {
	return func(d *SVGDrawer) {
		d.scale = scale
	}
}

func WithSize(width, height int) SVGDrawerOption {
	return func(d *SVGDrawer) {
		d.size = [2]int{width, height}
	}
}

func NewSVGDrawer(opts ...SVGDrawerOption) *SVGDrawer {
	d := &SVGDrawer{
		size:  [2]int{1000, 1000},
		scale: 1,
	}
	for _, opt := range opts {
		opt(d)
	}
	d.buffer.WriteString(`<?xml version="1.0" standalone="no"?>
<svg width="100vw" height="100vh" version="1.1" xmlns="http://www.w3.org/2000/svg">
`)
	d.buffer.WriteString("<g ")

	var transform strings.Builder
	transform.WriteString(`transform="`)
	// we invert the axis Y
	transform.WriteString(fmt.Sprintf(`translate(%f,%f) `, d.offset.X, float64(d.size[1])-d.offset.Y))
	transform.WriteString(fmt.Sprintf(`scale(%f, -%f) `, d.scale, d.scale))

	d.buffer.WriteString(transform.String())
	d.buffer.WriteString("\" >\n")
	return d
}

func (d *SVGDrawer) AddPoint(pt Point, styles ...string) {
	d.buffer.WriteString(`<circle cx="`)
	d.buffer.WriteString(fmt.Sprintf("%f", pt.X))
	d.buffer.WriteString(`" cy="`)
	d.buffer.WriteString(fmt.Sprintf("%f", pt.Y))
	d.buffer.WriteString(`" `)
	for i := 0; i < len(styles); i += 2 {
		d.buffer.WriteString(fmt.Sprintf(`%s="%s" `, styles[i], styles[i+1]))
	}
	d.buffer.WriteString("\n />")
}

func (d *SVGDrawer) AddPolygon(poly Polygon, styles ...string) {
	d.buffer.WriteString(`<path d="`)
	d.buffer.WriteString("M")
	for i, pt := range poly.outerRing {
		d.buffer.WriteString(fmt.Sprintf("%f,%f", pt.X, pt.Y))
		if i < len(poly.outerRing)-1 {
			d.buffer.WriteString(" ")
		}
	}

	for _, innerRing := range poly.innerRings {
		d.buffer.WriteString("zM")
		for i, pt := range innerRing {
			d.buffer.WriteString(fmt.Sprintf("%f,%f", pt.X, pt.Y))
			if i < len(innerRing)-1 {
				d.buffer.WriteString(" ")
			}
		}
	}

	d.buffer.WriteString(`" `)
	for i := 0; i < len(styles); i += 2 {
		d.buffer.WriteString(fmt.Sprintf(`%s="%s" `, styles[i], styles[i+1]))
	}
	d.buffer.WriteString(`fill="none" />`)
	d.buffer.WriteString("\n")
}

func (d *SVGDrawer) AddPart(piece OccupancyTable, step float64, offset Point, styles ...string) {
	for i, segment := range piece {
		for _, interval := range segment {
			height := interval.End - interval.Start
			x := offset.X + float64(i)*step
			y := offset.Y + interval.Start
			width := step
			d.AddSquare(x, y, width, height, styles...)
		}
	}
}

func (d *SVGDrawer) AddSquare(x, y, width, height float64, styles ...string) {
	d.buffer.WriteString(`<rect x="`)
	d.buffer.WriteString(fmt.Sprintf("%f", x))
	d.buffer.WriteString(`" y="`)
	d.buffer.WriteString(fmt.Sprintf("%f", y))
	d.buffer.WriteString(`" width="`)
	d.buffer.WriteString(fmt.Sprintf("%f", width))
	d.buffer.WriteString(`" height="`)
	d.buffer.WriteString(fmt.Sprintf("%f", height))
	d.buffer.WriteString(`" `)
	for i := 0; i < len(styles); i += 2 {
		d.buffer.WriteString(fmt.Sprintf(`%s="%s" `, styles[i], styles[i+1]))
	}
	d.buffer.WriteString(`fill="none" />`)
	d.buffer.WriteString("\n")
}

func (d *SVGDrawer) AddLine(x1, y1, x2, y2 float64, styles ...string) {
	d.buffer.WriteString(`<line x1="`)
	d.buffer.WriteString(fmt.Sprintf("%f", x1))
	d.buffer.WriteString(`" y1="`)
	d.buffer.WriteString(fmt.Sprintf("%f", y1))
	d.buffer.WriteString(`" x2="`)
	d.buffer.WriteString(fmt.Sprintf("%f", x2))
	d.buffer.WriteString(`" y2="`)
	d.buffer.WriteString(fmt.Sprintf("%f", y2))
	d.buffer.WriteString(`" `)
	for i := 0; i < len(styles); i += 2 {
		d.buffer.WriteString(fmt.Sprintf(`%s="%s" `, styles[i], styles[i+1]))
	}
	d.buffer.WriteString(`fill="none" />`)
	d.buffer.WriteString("\n")
}

func (d *SVGDrawer) AddText(pt Point, text string, styles ...string) {
	d.buffer.WriteString(`\n<text x="`)
	d.buffer.WriteString(fmt.Sprintf("%f", pt.X))
	d.buffer.WriteString(`" y="`)
	d.buffer.WriteString(fmt.Sprintf("%f", -pt.Y))

	d.buffer.WriteString(`" transform="scale(1,-1)" `)

	for i := 0; i < len(styles); i += 2 {
		d.buffer.WriteString(fmt.Sprintf(`%s="%s" `, styles[i], styles[i+1]))
	}
	d.buffer.WriteString(`fill="black" >`)
	d.buffer.WriteString(text)
	d.buffer.WriteString(`</text>`)
	d.buffer.WriteString("\n")
}

func (d *SVGDrawer) DrawCoordSystem(maxX, maxY int) {
	d.buffer.WriteString(`<defs>
<marker
      id="arrow"
      viewBox="0 0 10 10"
      refX="5"
      refY="5"
      markerWidth="6"
      markerHeight="6"
      orient="auto-start-reverse">
      <path d="M 0 0 L 10 5 L 0 10 z" />
    </marker>
</defs>
`)

	d.buffer.WriteString(`<line x1="0" y1="0" x2="0" `)
	d.buffer.WriteString(`y2="`)
	d.buffer.WriteString(fmt.Sprintf("%d", maxY))
	d.buffer.WriteString(`" stroke="black" stroke-width="1" marker-end="url(#arrow)" />
`)

	// d.buffer.WriteString(fmt.Sprintf(`<text x="-15" y="%d" text-anchor="middle" font-size="10" scale=(1, -1) >Y</text>
	// `, d.size[1]))

	d.buffer.WriteString(`<line x1="0" y1="0" x2="`)
	d.buffer.WriteString(fmt.Sprintf("%d", maxX))
	d.buffer.WriteString(`" `)
	d.buffer.WriteString(`y2="0" stroke="black" stroke-width="1" marker-end="url(#arrow)" />
`)
}

func (d *SVGDrawer) Write(w io.Writer) {
	d.buffer.WriteString("</g>\n")
	d.buffer.WriteString("</svg>\n")
	w.Write(d.buffer.Bytes())
}
