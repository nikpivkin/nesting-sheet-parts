package main

import "encoding/xml"

type Nesting struct {
	Problem  Problem    `xml:"problem"`
	Polygons []NPolygon `xml:"polygons>polygon"`
}

type Problem struct {
	XMLName xml.Name `xml:"problem"`
	Boards  []Piece  `xml:"boards>piece"`
	Lot     []Piece  `xml:"lot>piece"`
}

type Piece struct {
	XMLName   xml.Name  `xml:"piece"`
	ID        string    `xml:"id,attr"`
	Quantity  int       `xml:"quantity,attr"`
	Component Component `xml:"component"`
}

type Component struct {
	IDPolygon string `xml:"idPolygon,attr"`
	Type      int    `xml:"type,attr"`
	XOffset   int    `xml:"xOffset,attr"`
	YOffset   int    `xml:"yOffset,attr"`
}

type NPolygon struct {
	ID    string `xml:"id,attr"`
	N     int    `xml:"nVertices,attr"`
	Lines struct {
		Segment []Segment `xml:"segment"`
	} `xml:"lines"`
	XMin float32 `xml:"xMin"`
	XMax float32 `xml:"xMax"`
	YMin float32 `xml:"yMin"`
	YMax float32 `xml:"yMax"`
}

func (n *NPolygon) ToGeomPolygon() Polygon {
	var outer Ring
	for i := 0; i < n.N; i++ {
		outer = append(outer,
			NewPoint(n.Lines.Segment[i].X0, n.Lines.Segment[i].Y0),
			NewPoint(n.Lines.Segment[i].X1, n.Lines.Segment[i].Y1),
		)
	}
	return Polygon{outer, nil}
}

type Segment struct {
	N  int     `xml:"n,attr"`
	X0 float64 `xml:"x0,attr"`
	X1 float64 `xml:"x1,attr"`
	Y0 float64 `xml:"y0,attr"`
	Y1 float64 `xml:"y1,attr"`
}

func (n *Nesting) GetBoardSizes() (float32, float32) {
	if len(n.Problem.Boards) == 0 {
		panic("no board in nesting")
	}

	polyid := n.Problem.Boards[0].Component.IDPolygon
	for _, polygon := range n.Polygons {
		if polygon.ID == polyid {
			return polygon.XMax - polygon.XMin, polygon.YMax - polygon.YMin
		}
	}
	panic("board not found")
}

func (n *Nesting) GetParts() []Polygon {
	var parts []Polygon

	polygons := make(map[string]NPolygon)

	for _, polygon := range n.Polygons {
		polygons[polygon.ID] = polygon
	}

	for _, lot := range n.Problem.Lot {
		pieces := make([]Polygon, lot.Quantity)

		npoly := polygons[lot.Component.IDPolygon]

		for i := 0; i < lot.Quantity; i++ {
			pieces[i] = npoly.ToGeomPolygon()
		}
		parts = append(parts, pieces...)
	}

	return parts
}
