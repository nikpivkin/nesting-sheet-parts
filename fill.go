package main

import (
	"fmt"
	"sort"
)

// BottomLeftFill implements the Bottom-Left-Fill algorithm for
// placing a sequence of parts in a sheet
type BottomLeftFill struct {
	// the height of the sheet
	height float32
	// the maximum length of the sheet
	maxLength int
	// represents a table of strips
	vacancyTable map[int]Strip
}

// TODO: pass maxLength as option
func NewBottomLeftFill(height float32, maxLength int) *BottomLeftFill {
	return &BottomLeftFill{
		height:       height,
		maxLength:    maxLength,
		vacancyTable: make(map[int]Strip),
	}
}

// Run runs the Bottom-Left-Fill algorithm and returns a list of points
// representing the placement of the parts.
func (r *BottomLeftFill) Run(parts []*Part) {
	for _, part := range parts {
		proj := r.place(part)
		part.Offset = proj.offset
		part.BestOrientationNum = proj.orderNum
	}
}

func (r *BottomLeftFill) getVacancyStrip(num int) Strip {
	col, exists := r.vacancyTable[num]
	if !exists {
		// The sheet vacancy is from between 0 and the height of the sheet
		col = Strip{
			{Start: 0, End: float64(r.height)},
		}
		r.vacancyTable[num] = col
	}
	return col
}

// insert inserts the part into the occupancy table
func (r *BottomLeftFill) insert(proj projection) {
	for stripNum, strip := range proj.val {
		for intervalNum, rng := range strip {
			offseted := make([]Range, len(rng))
			for i, r := range rng {
				offseted[i] = r.Add(proj.offset.y)
			}
			r.insertStrip(proj.offset.column+stripNum, intervalNum, offseted...)
		}
	}
}

type Offset struct {
	column int
	y      float64
}

func (r *BottomLeftFill) insertStrip(stripNum int, rngNum int, strip ...Range) {
	vacantRange := r.vacancyTable[stripNum][rngNum]
	r.vacancyTable[stripNum] = insertSlice(r.vacancyTable[stripNum], rngNum, vacantRange.Split(strip)...)
}

// projection represents the projection of the part to sheet
type projection struct {
	offset   Offset
	orderNum int
	// it is a map of sheet strip number to a map of sheet range number to range
	val map[int]map[int][]Range
}

func (p projection) insert(stripNum int, rngNum int, rng Range) {
	strip, exists := p.val[stripNum]
	if !exists {
		strip = make(map[int][]Range)
		p.val[stripNum] = strip
	}
	strip[rngNum] = append(strip[rngNum], rng)
	p.val[stripNum] = strip
}

func (r *BottomLeftFill) place(part *Part) projection {
	projections := make([]projection, 0, len(part.Orientations))
	for i, orientation := range part.Orientations {
		projection := r.placeOrientation(orientation.occupancy, Offset{})
		projection.orderNum = i
		projections = append(projections, projection)
	}

	sort.Slice(projections, func(i, j int) bool {
		if projections[i].offset.column != projections[j].offset.column {
			return projections[i].offset.column < projections[j].offset.column
		}

		// TODO: if eq, then sort by angle
		return projections[i].offset.y < projections[j].offset.y
	})

	r.insert(projections[0])
	return projections[0]
}

func (r *BottomLeftFill) placeOrientation(part OccupancyTable, offset Offset) projection {
	if offset.column >= r.maxLength {
		panic(fmt.Sprintf("all parts cannot be placed, column %d reached", offset.column))
	}

	var (
		cursor     int // the index of the placed strip
		projection = projection{
			offset: offset, // offset is the current position of the part
			val:    make(map[int]map[int][]Range),
		}
	)

	for cursor != len(part) {
		strip := part[cursor]

		// TODO: check it in Run
		if strip.End() > float64(r.height) {
			// TODO return error
			panic("the end point of the part cannot be greater than ymax")
		}

		for _, stripRange := range strip {
			ok, rngNum, vacantRange := r.findVacantRange(projection.offset, cursor, stripRange)
			if !ok {
				// failed to place a segment of the piece, move to the next column
				return r.placeOrientation(part, Offset{
					column: projection.offset.column + 1,
					y:      0,
				})
			}

			newoffset := toFixed(vacantRange.Start-stripRange.Start, 4)
			projection.offset.y = max(projection.offset.y, newoffset)

			for column, projectionStrip := range projection.val {
				for vacantRngNum, projectionRanges := range projectionStrip {
					for _, projectionRange := range projectionRanges {
						if !r.canPlace(projection.offset, column, vacantRngNum, projectionRange) {
							return r.placeOrientation(part, projection.offset)
						}
					}
				}
			}

			projection.insert(cursor, rngNum, stripRange)
		}
		cursor++
	}

	return projection
}

func (f *BottomLeftFill) findVacantRange(offset Offset, colOffset int, rangeToPlace Range) (bool, int, Range) {
	for idx, vacantRange := range f.getVacancyStrip(offset.column + colOffset) {
		if vacantRange.Includes(rangeToPlace.Add(offset.y)) ||
			vacantRange.Length() >= rangeToPlace.Add(offset.y).Length() &&
				vacantRange.Start >= rangeToPlace.Start+offset.y {
			return true, idx, vacantRange
		}
	}
	return false, 0, Range{}
}

func (r *BottomLeftFill) canPlace(offset Offset, columnOffset int, rngNum int, rng Range) bool {
	vacantRng := r.vacancyTable[offset.column+columnOffset][rngNum]
	return vacantRng.Includes(rng.Add(offset.y))
}

func (r *BottomLeftFill) getVacancyTable() OccupancyTable {
	table := make(OccupancyTable, len(r.vacancyTable))
	keys := make([]int, 0, len(r.vacancyTable))
	for k := range r.vacancyTable {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		table[k] = r.vacancyTable[k]
	}
	return table
}
