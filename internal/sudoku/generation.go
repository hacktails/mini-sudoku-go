package sudoku

// generatePuzzle produces a puzzle for the given size/difficulty.
// It prefers curated puzzles and falls back to generated ones with uniqueness checks.
func generatePuzzle(set puzzleSet, diff difficulty) puzzle {
	if p, ok := randomFromLibrary(set, diff); ok {
		return p
	}

	targetClues := clueCount(set.size, diff)
	attempts := 0
	for attempts < 60 {
		attempts++
		solution := generateSolution(set)
		puzzleGrid := carvePuzzle(solution, set, targetClues)
		if countSolutions(puzzleGrid, set, 2) != 1 {
			continue
		}
		rated := rateDifficulty(puzzleGrid, set)
		if rated == diff {
			return puzzle{puzzle: puzzleGrid, solution: solution}
		}
	}

	solution := generateSolution(set)
	puzzleGrid := carvePuzzle(solution, set, targetClues)
	return puzzle{puzzle: puzzleGrid, solution: solution}
}

// carvePuzzle removes values from a solved grid while keeping uniqueness.
func carvePuzzle(solution []uint8, set puzzleSet, targetClues int) []uint8 {
	puzzleGrid := copyGrid(solution)
	if targetClues < 0 {
		targetClues = 0
	}
	if targetClues > len(puzzleGrid) {
		targetClues = len(puzzleGrid)
	}
	removeCount := len(puzzleGrid) - targetClues
	removed := 0
	for _, idx := range rng.Perm(len(puzzleGrid)) {
		if removed >= removeCount {
			break
		}
		keep := puzzleGrid[idx]
		puzzleGrid[idx] = 0
		if countSolutions(puzzleGrid, set, 2) != 1 {
			puzzleGrid[idx] = keep
			continue
		}
		removed++
	}
	return puzzleGrid
}

// generateSolution builds a full valid Sudoku solution by permuting a base grid.
func generateSolution(set puzzleSet) []uint8 {
	size := set.size
	base := make([]uint8, size*size)
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			value := (row*set.boxCols + row/set.boxRows + col) % size
			base[idx(row, col, size)] = uint8(value + 1)
		}
	}

	digitPerm := rng.Perm(size)
	for i := range base {
		base[i] = uint8(digitPerm[base[i]-1] + 1)
	}

	rowOrder := shuffledBandIndices(size, set.boxRows)
	colOrder := shuffledBandIndices(size, set.boxCols)
	grid := make([]uint8, size*size)
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			grid[idx(row, col, size)] = base[idx(rowOrder[row], colOrder[col], size)]
		}
	}
	return grid
}

// shuffledBandIndices randomizes rows/cols by bands to preserve validity.
func shuffledBandIndices(size, band int) []int {
	bands := size / band
	bandOrder := rng.Perm(bands)
	order := make([]int, 0, size)
	for _, b := range bandOrder {
		rows := rng.Perm(band)
		for _, r := range rows {
			order = append(order, b*band+r)
		}
	}
	return order
}

// rateDifficulty estimates difficulty by running a logic-first solver.
func rateDifficulty(puzzleGrid []uint8, set puzzleSet) difficulty {
	solved, guesses := solveWithLogic(puzzleGrid, set)
	if !solved {
		return diffHard
	}
	if guesses == 0 {
		return diffEasy
	}
	if guesses <= 1 {
		return diffMedium
	}
	return diffHard
}

// solveWithLogic runs a solver and returns whether it solved plus guess count.
func solveWithLogic(puzzleGrid []uint8, set puzzleSet) (bool, int) {
	grid := copyGrid(puzzleGrid)
	solved, guesses := solveRecursive(grid, set)
	return solved, guesses
}

// solveRecursive solves the puzzle using logic, falling back to search.
func solveRecursive(grid []uint8, set puzzleSet) (bool, int) {
	for {
		progress, ok := applyLogic(grid, set)
		if !ok {
			return false, 0
		}
		if !progress {
			break
		}
	}

	emptyIndex := -1
	var candidates uint16
	minCount := 10
	for i, value := range grid {
		if value != 0 {
			continue
		}
		row := i / set.size
		col := i % set.size
		mask := candidatesFor(grid, row, col, set)
		count := bitCount(mask)
		if count == 0 {
			return false, 0
		}
		if count < minCount {
			minCount = count
			emptyIndex = i
			candidates = mask
			if count == 1 {
				break
			}
		}
	}

	if emptyIndex == -1 {
		return true, 0
	}

	row := emptyIndex / set.size
	col := emptyIndex % set.size
	for _, value := range maskToValues(candidates, set.size) {
		next := copyGrid(grid)
		next[idx(row, col, set.size)] = uint8(value)
		solved, guesses := solveRecursive(next, set)
		if solved {
			return true, guesses + 1
		}
	}

	return false, 0
}

// applyLogic fills obvious singles until no further progress is possible.
func applyLogic(grid []uint8, set puzzleSet) (bool, bool) {
	progress := false
	for {
		candidates := make([]uint16, len(grid))
		for i, value := range grid {
			if value != 0 {
				continue
			}
			row := i / set.size
			col := i % set.size
			mask := candidatesFor(grid, row, col, set)
			if mask == 0 {
				return false, false
			}
			candidates[i] = mask
		}

		step := false
		for i, mask := range candidates {
			if mask == 0 {
				continue
			}
			if bitCount(mask) == 1 {
				value := firstBit(mask)
				grid[i] = uint8(value)
				step = true
			}
		}

		for row := 0; row < set.size; row++ {
			if applyHiddenSinglesRow(grid, candidates, set, row) {
				step = true
			}
		}
		for col := 0; col < set.size; col++ {
			if applyHiddenSinglesCol(grid, candidates, set, col) {
				step = true
			}
		}
		for boxRow := 0; boxRow < set.size; boxRow += set.boxRows {
			for boxCol := 0; boxCol < set.size; boxCol += set.boxCols {
				if applyHiddenSinglesBox(grid, candidates, set, boxRow, boxCol) {
					step = true
				}
			}
		}

		if !step {
			break
		}
		progress = true
	}
	return true, progress
}

// applyHiddenSinglesRow finds hidden singles in a row.
func applyHiddenSinglesRow(grid []uint8, candidates []uint16, set puzzleSet, row int) bool {
	counts := make([]int, set.size)
	pos := make([]int, set.size)
	for col := 0; col < set.size; col++ {
		i := idx(row, col, set.size)
		if grid[i] != 0 {
			continue
		}
		mask := candidates[i]
		for n := 0; n < set.size; n++ {
			if mask&(1<<uint(n)) != 0 {
				counts[n]++
				pos[n] = i
			}
		}
	}
	progress := false
	for n, count := range counts {
		if count == 1 {
			grid[pos[n]] = uint8(n + 1)
			progress = true
		}
	}
	return progress
}

// applyHiddenSinglesCol finds hidden singles in a column.
func applyHiddenSinglesCol(grid []uint8, candidates []uint16, set puzzleSet, col int) bool {
	counts := make([]int, set.size)
	pos := make([]int, set.size)
	for row := 0; row < set.size; row++ {
		i := idx(row, col, set.size)
		if grid[i] != 0 {
			continue
		}
		mask := candidates[i]
		for n := 0; n < set.size; n++ {
			if mask&(1<<uint(n)) != 0 {
				counts[n]++
				pos[n] = i
			}
		}
	}
	progress := false
	for n, count := range counts {
		if count == 1 {
			grid[pos[n]] = uint8(n + 1)
			progress = true
		}
	}
	return progress
}

// applyHiddenSinglesBox finds hidden singles in a box.
func applyHiddenSinglesBox(grid []uint8, candidates []uint16, set puzzleSet, boxRow, boxCol int) bool {
	counts := make([]int, set.size)
	pos := make([]int, set.size)
	for r := 0; r < set.boxRows; r++ {
		for c := 0; c < set.boxCols; c++ {
			row := boxRow + r
			col := boxCol + c
			i := idx(row, col, set.size)
			if grid[i] != 0 {
				continue
			}
			mask := candidates[i]
			for n := 0; n < set.size; n++ {
				if mask&(1<<uint(n)) != 0 {
					counts[n]++
					pos[n] = i
				}
			}
		}
	}
	progress := false
	for n, count := range counts {
		if count == 1 {
			grid[pos[n]] = uint8(n + 1)
			progress = true
		}
	}
	return progress
}

// countSolutions counts solutions up to a limit.
func countSolutions(puzzleGrid []uint8, set puzzleSet, limit int) int {
	grid := copyGrid(puzzleGrid)
	return countSolutionsRecursive(grid, set, limit)
}

// countSolutionsRecursive explores solutions with backtracking.
func countSolutionsRecursive(grid []uint8, set puzzleSet, limit int) int {
	emptyIndex := -1
	var candidates uint16
	minCount := 10
	for i, value := range grid {
		if value != 0 {
			continue
		}
		row := i / set.size
		col := i % set.size
		mask := candidatesFor(grid, row, col, set)
		count := bitCount(mask)
		if count == 0 {
			return 0
		}
		if count < minCount {
			minCount = count
			emptyIndex = i
			candidates = mask
			if count == 1 {
				break
			}
		}
	}
	if emptyIndex == -1 {
		return 1
	}
	row := emptyIndex / set.size
	col := emptyIndex % set.size
	total := 0
	for _, value := range maskToValues(candidates, set.size) {
		grid[idx(row, col, set.size)] = uint8(value)
		total += countSolutionsRecursive(grid, set, limit-total)
		if total >= limit {
			return total
		}
	}
	grid[idx(row, col, set.size)] = 0
	return total
}

// candidatesFor returns a bitmask of legal values for a cell.
func candidatesFor(grid []uint8, row, col int, set puzzleSet) uint16 {
	if grid[idx(row, col, set.size)] != 0 {
		return 0
	}
	used := uint16(0)
	for c := 0; c < set.size; c++ {
		value := grid[idx(row, c, set.size)]
		if value != 0 {
			used |= 1 << uint(value-1)
		}
	}
	for r := 0; r < set.size; r++ {
		value := grid[idx(r, col, set.size)]
		if value != 0 {
			used |= 1 << uint(value-1)
		}
	}
	boxRow := (row / set.boxRows) * set.boxRows
	boxCol := (col / set.boxCols) * set.boxCols
	for r := boxRow; r < boxRow+set.boxRows; r++ {
		for c := boxCol; c < boxCol+set.boxCols; c++ {
			value := grid[idx(r, c, set.size)]
			if value != 0 {
				used |= 1 << uint(value-1)
			}
		}
	}
	full := uint16(1<<uint(set.size)) - 1
	return full &^ used
}

// bitCount returns the number of set bits in a mask.
func bitCount(mask uint16) int {
	count := 0
	for mask > 0 {
		mask &= mask - 1
		count++
	}
	return count
}

// firstBit returns the 1-based index of the lowest set bit.
func firstBit(mask uint16) int {
	for i := 0; i < 16; i++ {
		if mask&(1<<uint(i)) != 0 {
			return i + 1
		}
	}
	return 0
}

// maskToValues converts a mask to a list of candidate values.
func maskToValues(mask uint16, size int) []int {
	values := make([]int, 0, size)
	for i := 0; i < size; i++ {
		if mask&(1<<uint(i)) != 0 {
			values = append(values, i+1)
		}
	}
	return values
}
