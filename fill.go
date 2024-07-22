package main

import "fmt"

// BottomLeftFill implements the Bottom-Left-Fill algorithm for
// placing a sequence of parts in a sheet
type BottomLeftFill struct {
	// the height of the sheet
	height int
	// the maximum length of the sheet
	maxLength int
	// represents a table of strips
	occupancyTable map[int]Strip
}

// TODO: pass maxLength as option
func NewBottomLeftFill(height int, maxLength int) *BottomLeftFill {
	return &BottomLeftFill{
		height:         height,
		maxLength:      maxLength,
		occupancyTable: make(map[int]Strip),
	}
}

// Run runs the Bottom-Left-Fill algorithm and returns a list of points
// representing the placement of the parts.
func (r *BottomLeftFill) Run(parts []OccupancyTable) []Offset {
	res := make([]Offset, 0, len(parts))
	for _, part := range parts {
		offset := r.place(part, 0)
		res = append(res, offset)
	}
	return res
}

func (r *BottomLeftFill) getVacancyStrip(num int) Strip {
	col, exists := r.occupancyTable[num]
	if !exists {
		// The sheet vacancy is from between 0 and the height of the sheet
		col = Strip{
			{Start: 0, End: float64(r.height)},
		}
		r.occupancyTable[num] = col
	}
	return col
}

// insert inserts the part into the occupancy table
func (r *BottomLeftFill) insert(projection projection, offset Offset) {
	for stripNum, strip := range projection {
		for intervalNum, rng := range strip {
			r.insertStrip(offset.column+stripNum, intervalNum, rng.Add(offset.y))
		}
	}
}

type Offset struct {
	column int
	y      float64
}

func (r *BottomLeftFill) insertStrip(stripNum int, rngNum int, strip ...Range) {
	vacantRange := r.occupancyTable[stripNum][rngNum]
	r.occupancyTable[stripNum] = insertSlice(r.occupancyTable[stripNum], rngNum, vacantRange.Split(strip)...)
}

// projection represents the projection of the part to sheet
// it is a map of sheet strip number to a map of sheet range number to range
type projection map[int]map[int]Range

func (p projection) insert(stripNum int, rngNum int, rng Range) {
	strip, exists := p[stripNum]
	if !exists {
		strip = make(map[int]Range)
		p[stripNum] = strip
	}
	strip[rngNum] = rng
	p[stripNum] = strip
}

func (r *BottomLeftFill) place(part OccupancyTable, column int) Offset {

	if column >= r.maxLength {
		panic(fmt.Sprintf("all parts cannot be placed, column %d reached", column))
	}

	var (
		offset     = Offset{column, 0} // offset is the current position of the part
		cursor     int                 // the index of the placed strip
		projection = make(projection)
	)

	for cursor != len(part) {
		strip := part[cursor]

		// TODO: check it in Run
		if strip.End() > float64(r.height) {
			// TODO return error
			panic("the end point of the part cannot be greater than ymax")
		}

		for _, stripRange := range strip {
			ok, rngNum, vacantRange := r.findVacantRange(offset, cursor, stripRange)
			if !ok {
				// failed to place a segment of the piece, move to the next column
				return r.place(part, offset.column+1)
			}

			newoffset := vacantRange.Start - stripRange.Start
			offset.y = max(offset.y, newoffset)

			for column, projectionStrip := range projection {
				for vacantRngNum, projectionRng := range projectionStrip {
					if !r.canPlace(offset, column, vacantRngNum, projectionRng) {
						return r.place(part, offset.column+1)
					}
				}
			}

			projection.insert(cursor, rngNum, stripRange)
		}
		cursor++
	}

	r.insert(projection, offset)
	return offset
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
	vacantRng := r.occupancyTable[offset.column+columnOffset][rngNum]
	return vacantRng.Includes(rng.Add(offset.y))
}

func insertSlice[E any](slice []E, index int, elements ...E) []E {
	if index < 0 || index > len(slice) {
		return slice
	}

	resultSlice := append(slice[:index], elements...)
	resultSlice = append(resultSlice, slice[index+1:]...)

	return resultSlice
}
