package main

// BottomLeftFill implements the Bottom-Left-Fill algorithm for
// placing a sequence of parts in a sheet
type BottomLeftFill struct {
	// the maximum height of the sheet
	ymax int
	// represents a table of strips
	occupancyTable map[int]Strip
}

func NewBottomLeftFill(ymax int) *BottomLeftFill {
	return &BottomLeftFill{
		ymax:           ymax,
		occupancyTable: make(map[int]Strip),
	}
}

// Run runs the Bottom-Left-Fill algorithm and returns a list of points
// representing the placement of the parts.
func (r *BottomLeftFill) Run(parts []OccypancyTable) []Offset {
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
			{Start: 0, End: float64(r.ymax)},
		}
		r.occupancyTable[num] = col
	}
	return col
}

// insert inserts the part into the occupancy table
func (r *BottomLeftFill) insert(proection proection, offset Offset) {
	for stripNum, strip := range proection {
		for intervalNum, rng := range strip {
			r.insertStrip(offset.column+stripNum, intervalNum, rng.Add(offset.y))
		}
	}
}

// TODO: use instead of point
type Offset struct {
	column int
	y      float64
}

func (r *BottomLeftFill) insertStrip(stripNum int, rngNum int, strip ...Range) {
	vacantRange := r.occupancyTable[stripNum][rngNum]
	r.occupancyTable[stripNum] = insertSlice(r.occupancyTable[stripNum], rngNum, vacantRange.Split(strip)...)
}

// TODO: must be map of Strip
// proection represents the proection of the part
type proection map[int]map[int]Range

func (p proection) insert(stripNum int, rngNum int, rng Range) {
	strip, exists := p[stripNum]
	if !exists {
		strip = make(map[int]Range)
		p[stripNum] = strip
	}
	strip[rngNum] = rng
	p[stripNum] = strip
}

func (r *BottomLeftFill) place(part OccypancyTable, column int) Offset {

	var (
		offset    = Offset{column, 0} // offset is the current position of the part
		cursor    int                 // the index of the placed strip
		proection = make(proection)
	)

	for cursor != len(part) {
		strip := part[cursor]

		if strip.End() > float64(r.ymax) {
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

			// if the offset has changed, it must be checked if the previous parts can be placed
			for prevStripNum, prevStrip := range part[:cursor+1] {
				for _, rng := range prevStrip {
					if ok, _ := r.canPlace(offset, prevStripNum, rng); !ok {
						// the previous segments cannot be placed with the new offset,
						// so we move to the next column
						return r.place(part, offset.column+1)
					}
				}
			}

			proection.insert(cursor, rngNum, stripRange)
		}
		cursor++
	}

	r.insert(proection, offset)
	return offset
}

func (r *BottomLeftFill) findVacantRange(offset Offset, columnOffset int, rng Range) (bool, int, Range) {
	for segmentNum, freeSegment := range r.getVacancyStrip(int(offset.column) + columnOffset) {
		if freeSegment.Includes(rng.Add(offset.y)) ||
			freeSegment.Length() >= rng.Add(offset.y).Length() && freeSegment.Start >= rng.Start {
			return true, segmentNum, freeSegment
		}
	}
	return false, 0, Range{}
}

func (r *BottomLeftFill) canPlace(offset Offset, columnOffset int, rng Range) (bool, Range) {
	ok, _, vacantRange := r.findVacantRange(offset, columnOffset, rng)

	if !ok || vacantRange.End < rng.End+offset.y {
		return false, Range{}
	}
	return true, vacantRange
}

func insertSlice[E any](slice []E, index int, elements ...E) []E {
	if index < 0 || index > len(slice) {
		return slice
	}

	resultSlice := append(slice[:index], elements...)
	resultSlice = append(resultSlice, slice[index+1:]...)

	return resultSlice
}
