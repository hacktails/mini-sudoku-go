package sudoku

// copyGrid clones a sudoku grid.
func copyGrid(src []uint8) []uint8 {
	dst := make([]uint8, len(src))
	copy(dst, src)
	return dst
}

// copyNotes clones a notes slice.
func copyNotes(src []uint16) []uint16 {
	dst := make([]uint16, len(src))
	copy(dst, src)
	return dst
}

// idx returns the flat index for a row/col in a size x size grid.
func idx(row, col, size int) int {
	return row*size + col
}
