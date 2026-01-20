package sudoku

import (
	"time"
)

// setSize switches the puzzle size and resets the game.
func (m *model) setSize(size int) {
	set, ok := puzzleSets[size]
	if !ok {
		return
	}
	m.set = set
	m.newPuzzle()
}

// setDifficulty switches difficulty and resets the game.
func (m *model) setDifficulty(diff difficulty) {
	m.difficulty = diff
	m.newPuzzle()
}

// newPuzzle generates and loads a new puzzle.
func (m *model) newPuzzle() {
	m.setPuzzle(generatePuzzle(m.set, m.difficulty))
	m.notes = make([]uint16, m.set.size*m.set.size)
	m.mistakes = 0
	m.hintsUsed = 0
	m.solved = false
	m.elapsedAtSolve = 0
	m.gameOver = false
	m.clearHistory()
	m.flash("New puzzle")
	m.autoSave()
}

// setPuzzle replaces the current puzzle data.
func (m *model) setPuzzle(p puzzle) {
	m.puzzle = p
	m.grid = copyGrid(p.puzzle)
	m.row = 0
	m.col = 0
	m.start = time.Now()
	m.elapsedAtSolve = 0
	m.gameOver = false
}

// reset restores the puzzle to its initial state.
func (m *model) reset() {
	m.grid = copyGrid(m.puzzle.puzzle)
	m.row = 0
	m.col = 0
	m.start = time.Now()
	m.notes = make([]uint16, m.set.size*m.set.size)
	m.mistakes = 0
	m.hintsUsed = 0
	m.solved = false
	m.elapsedAtSolve = 0
	m.gameOver = false
	m.clearHistory()
	m.flash("Reset puzzle")
	m.autoSave()
}

// move shifts the selection by the given delta.
func (m *model) move(dr, dc int) {
	m.row = clamp(m.row+dr, 0, m.set.size-1)
	m.col = clamp(m.col+dc, 0, m.set.size-1)
}

// clamp restricts a value to a [min,max] range.
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// isFixed returns true if the cell is part of the original puzzle.
func (m model) isFixed(row, col int) bool {
	return m.puzzle.puzzle[idx(row, col, m.set.size)] != 0
}

// isSolved returns true if the grid matches the solution or contains no conflicts.
func (m model) isSolved() bool {
	if equalGrid(m.grid, m.puzzle.solution) {
		return true
	}
	for _, v := range m.grid {
		if v == 0 {
			return false
		}
	}
	return !m.anyConflicts()
}

// equalGrid compares two grids for equality.
func equalGrid(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// anyConflicts returns true if any cell violates Sudoku rules.
func (m model) anyConflicts() bool {
	for row := 0; row < m.set.size; row++ {
		for col := 0; col < m.set.size; col++ {
			if m.hasConflict(row, col) {
				return true
			}
		}
	}
	return false
}

// hasConflict checks whether a specific cell violates Sudoku rules.
func (m model) hasConflict(row, col int) bool {
	value := m.grid[idx(row, col, m.set.size)]
	if value == 0 {
		return false
	}

	for c := 0; c < m.set.size; c++ {
		if c != col && m.grid[idx(row, c, m.set.size)] == value {
			return true
		}
	}

	for r := 0; r < m.set.size; r++ {
		if r != row && m.grid[idx(r, col, m.set.size)] == value {
			return true
		}
	}

	boxRow := (row / m.set.boxRows) * m.set.boxRows
	boxCol := (col / m.set.boxCols) * m.set.boxCols
	for r := boxRow; r < boxRow+m.set.boxRows; r++ {
		for c := boxCol; c < boxCol+m.set.boxCols; c++ {
			if (r != row || c != col) && m.grid[idx(r, c, m.set.size)] == value {
				return true
			}
		}
	}

	return false
}

// setValue writes a value to the current cell and updates state.
func (m *model) setValue(value uint8) {
	if m.isFixed(m.row, m.col) {
		return
	}
	index := idx(m.row, m.col, m.set.size)
	if m.grid[index] == value {
		return
	}
	m.pushUndo()
	m.grid[index] = value
	if value != 0 {
		m.notes[index] = 0
		m.pruneNotes(m.row, m.col, value)
		if m.puzzle.solution[index] != value {
			m.mistakes++
			m.flash("Mistake")
			if m.strictMode && m.mistakes >= maxMistakes {
				m.gameOver = true
				m.elapsedAtSolve = int64(time.Since(m.start).Seconds())
				m.flash("Game over")
			}
		}
	}
	if !m.gameOver {
		m.checkSolved()
	}
	m.autoSave()
}

// clearValue clears the current cell.
func (m *model) clearValue() {
	if m.isFixed(m.row, m.col) {
		return
	}
	index := idx(m.row, m.col, m.set.size)
	if m.grid[index] == 0 && m.notes[index] == 0 {
		return
	}
	m.pushUndo()
	m.grid[index] = 0
	m.notes[index] = 0
	if !m.gameOver {
		m.checkSolved()
	}
	m.autoSave()
}

// toggleNote toggles a candidate note in the selected cell.
func (m *model) toggleNote(row, col, value int) {
	if m.isFixed(row, col) {
		return
	}
	index := idx(row, col, m.set.size)
	if m.grid[index] != 0 {
		return
	}
	mask := uint16(1 << uint(value-1))
	m.pushUndo()
	if m.notes[index]&mask != 0 {
		m.notes[index] &^= mask
	} else {
		m.notes[index] |= mask
	}
	m.autoSave()
}

// clearNotes removes all notes from a cell.
func (m *model) clearNotes(row, col int) {
	if m.isFixed(row, col) {
		return
	}
	index := idx(row, col, m.set.size)
	if m.notes[index] == 0 {
		return
	}
	m.pushUndo()
	m.notes[index] = 0
	m.autoSave()
}

// pruneNotes removes a value from notes in the affected row/col/box.
func (m *model) pruneNotes(row, col int, value uint8) {
	if value == 0 {
		return
	}
	mask := uint16(1 << uint(value-1))
	for c := 0; c < m.set.size; c++ {
		if c == col {
			continue
		}
		i := idx(row, c, m.set.size)
		m.notes[i] &^= mask
	}
	for r := 0; r < m.set.size; r++ {
		if r == row {
			continue
		}
		i := idx(r, col, m.set.size)
		m.notes[i] &^= mask
	}
	boxRow := (row / m.set.boxRows) * m.set.boxRows
	boxCol := (col / m.set.boxCols) * m.set.boxCols
	for r := boxRow; r < boxRow+m.set.boxRows; r++ {
		for c := boxCol; c < boxCol+m.set.boxCols; c++ {
			if r == row && c == col {
				continue
			}
			i := idx(r, c, m.set.size)
			m.notes[i] &^= mask
		}
	}
}

// pushUndo records the current state for undo.
func (m *model) pushUndo() {
	snap := snapshot{
		grid:      copyGrid(m.grid),
		notes:     copyNotes(m.notes),
		row:       m.row,
		col:       m.col,
		mistakes:  m.mistakes,
		hintsUsed: m.hintsUsed,
	}
	m.undoStack = append(m.undoStack, snap)
	m.redoStack = nil
}

// undo reverts to the previous snapshot.
func (m *model) undo() {
	if len(m.undoStack) == 0 {
		return
	}
	snap := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]
	m.redoStack = append(m.redoStack, m.snapshot())
	m.applySnapshot(snap)
	m.flash("Undo")
	m.autoSave()
}

// redo reapplies a reverted snapshot.
func (m *model) redo() {
	if len(m.redoStack) == 0 {
		return
	}
	snap := m.redoStack[len(m.redoStack)-1]
	m.redoStack = m.redoStack[:len(m.redoStack)-1]
	m.undoStack = append(m.undoStack, m.snapshot())
	m.applySnapshot(snap)
	m.flash("Redo")
	m.autoSave()
}

// snapshot captures the current game state for undo/redo.
func (m *model) snapshot() snapshot {
	return snapshot{
		grid:      copyGrid(m.grid),
		notes:     copyNotes(m.notes),
		row:       m.row,
		col:       m.col,
		mistakes:  m.mistakes,
		hintsUsed: m.hintsUsed,
	}
}

// applySnapshot restores a snapshot.
func (m *model) applySnapshot(snap snapshot) {
	m.grid = copyGrid(snap.grid)
	m.notes = copyNotes(snap.notes)
	m.row = snap.row
	m.col = snap.col
	m.mistakes = snap.mistakes
	m.hintsUsed = snap.hintsUsed
	m.solved = m.isSolved()
}

// clearHistory clears undo/redo stacks.
func (m *model) clearHistory() {
	m.undoStack = nil
	m.redoStack = nil
}

// applyHint applies a hint (logic first, then a reveal).
func (m *model) applyHint() {
	if m.solved || m.gameOver {
		return
	}
	if m.applyLogicalHint() {
		m.autoSave()
		return
	}

	empties := make([]int, 0)
	for i, value := range m.grid {
		if value == 0 {
			empties = append(empties, i)
		}
	}
	if len(empties) == 0 {
		return
	}
	index := empties[rng.Intn(len(empties))]
	row := index / m.set.size
	col := index % m.set.size
	if m.isFixed(row, col) {
		return
	}
	m.pushUndo()
	m.grid[index] = m.puzzle.solution[index]
	m.notes[index] = 0
	m.pruneNotes(row, col, m.puzzle.solution[index])
	m.hintsUsed++
	m.checkSolved()
	m.flash("Hint used")
	m.autoSave()
}

// applyLogicalHint applies a single logical hint if available.
func (m *model) applyLogicalHint() bool {
	candidates := make([]uint16, len(m.grid))
	for i, value := range m.grid {
		if value != 0 {
			continue
		}
		row := i / m.set.size
		col := i % m.set.size
		candidates[i] = candidatesFor(m.grid, row, col, m.set)
	}

	for i, mask := range candidates {
		if mask == 0 {
			continue
		}
		if bitCount(mask) == 1 {
			value := uint8(firstBit(mask))
			row := i / m.set.size
			col := i % m.set.size
			return m.applyHintValue(row, col, value, "Hint: single candidate")
		}
	}

	for row := 0; row < m.set.size; row++ {
		if value, col, ok := hiddenSingleInRow(m.grid, candidates, m.set, row); ok {
			return m.applyHintValue(row, col, value, "Hint: hidden single")
		}
	}
	for col := 0; col < m.set.size; col++ {
		if value, row, ok := hiddenSingleInCol(m.grid, candidates, m.set, col); ok {
			return m.applyHintValue(row, col, value, "Hint: hidden single")
		}
	}
	for boxRow := 0; boxRow < m.set.size; boxRow += m.set.boxRows {
		for boxCol := 0; boxCol < m.set.size; boxCol += m.set.boxCols {
			if value, row, col, ok := hiddenSingleInBox(m.grid, candidates, m.set, boxRow, boxCol); ok {
				return m.applyHintValue(row, col, value, "Hint: hidden single")
			}
		}
	}

	return false
}

// applyHintValue fills a hinted value and tracks usage.
func (m *model) applyHintValue(row, col int, value uint8, message string) bool {
	if m.isFixed(row, col) || m.grid[idx(row, col, m.set.size)] != 0 {
		return false
	}
	m.pushUndo()
	index := idx(row, col, m.set.size)
	m.grid[index] = value
	m.notes[index] = 0
	m.pruneNotes(row, col, value)
	m.hintsUsed++
	m.checkSolved()
	m.flash(message)
	return true
}

// hiddenSingleInRow searches for a hidden single in a row.
func hiddenSingleInRow(grid []uint8, candidates []uint16, set puzzleSet, row int) (uint8, int, bool) {
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
				pos[n] = col
			}
		}
	}
	for n, count := range counts {
		if count == 1 {
			return uint8(n + 1), pos[n], true
		}
	}
	return 0, 0, false
}

// hiddenSingleInCol searches for a hidden single in a column.
func hiddenSingleInCol(grid []uint8, candidates []uint16, set puzzleSet, col int) (uint8, int, bool) {
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
				pos[n] = row
			}
		}
	}
	for n, count := range counts {
		if count == 1 {
			return uint8(n + 1), pos[n], true
		}
	}
	return 0, 0, false
}

// hiddenSingleInBox searches for a hidden single in a box.
func hiddenSingleInBox(grid []uint8, candidates []uint16, set puzzleSet, boxRow, boxCol int) (uint8, int, int, bool) {
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
	for n, count := range counts {
		if count == 1 {
			idxVal := pos[n]
			return uint8(n + 1), idxVal / set.size, idxVal % set.size, true
		}
	}
	return 0, 0, 0, false
}

// checkSolved updates solved state and best time.
func (m *model) checkSolved() {
	if m.solved || m.gameOver {
		return
	}
	if m.isSolved() {
		m.solved = true
		m.elapsedAtSolve = int64(time.Since(m.start).Seconds())
		m.updateBestTime()
		m.flash("Solved")
	}
}

// updateBestTime stores the best solve time for the current size/difficulty.
func (m *model) updateBestTime() {
	elapsed := m.elapsedAtSolve
	if elapsed == 0 {
		elapsed = int64(time.Since(m.start).Seconds())
	}
	key := statsKey(m.set.size, m.difficulty)
	if best, ok := m.stats.Best[key]; !ok || elapsed < best {
		if m.stats.Best == nil {
			m.stats.Best = map[string]int64{}
		}
		m.stats.Best[key] = elapsed
		_ = saveStats(m.stats)
	}
}

// flash shows a short-lived status message.
func (m *model) flash(message string) {
	m.flashMessage = message
	m.flashUntil = time.Now().Add(2 * time.Second)
}

// autoSave persists the current state without surfacing errors.
func (m *model) autoSave() {
	_ = m.save()
}
