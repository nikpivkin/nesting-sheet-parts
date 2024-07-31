package main

func rangeSlice(start, end int, step int) []int {
	var result []int
	for i := start; i < end; i += step {
		result = append(result, i)
	}
	return result
}

func swapSliceParts(slice []int, index int) []int {
	if index < 0 || index >= len(slice) {
		return slice
	}

	part1 := make([]int, len(slice[index+1:]))
	copy(part1, slice[index+1:])

	part2 := make([]int, len(slice[:index]))
	copy(part2, slice[:index])

	middle := slice[index]

	newSlice := append(part1, middle)
	newSlice = append(newSlice, part2...)

	return newSlice
}

func swap(s []int, i, j int) {
	s[i], s[j] = s[j], s[i]
}
