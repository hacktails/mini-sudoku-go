package sudoku

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// boardView renders the full Sudoku board as a string.
func (m model) boardView() string {
	var lines []string
	for row := 0; row < m.set.size; row++ {
		lines = append(lines, m.rowView(row)...)
		if row == m.set.size-1 {
			continue
		}
		gapLines := rowGapLines
		if (row+1)%m.set.boxRows == 0 {
			gapLines = boxGapLines
		}
		lines = append(lines, blankLines(gapLines, boardWidth(m.set.size, m.set.boxCols))...)
	}
	return strings.Join(lines, "\n")
}

// rowView renders a single logical row (cellH lines).
func (m model) rowView(row int) []string {
	lines := make([]string, cellH)
	for col := 0; col < m.set.size; col++ {
		gap := ""
		if col > 0 {
			if col%m.set.boxCols == 0 {
				gap = boxGap
			} else {
				gap = cellGap
			}
		}

		value := m.grid[idx(row, col, m.set.size)]
		cell := m.renderCellLines(row, col, value)
		for i := 0; i < cellH; i++ {
			if gap != "" {
				lines[i] += gap
			}
			lines[i] += cell[i]
		}
	}
	return lines
}

// renderCellLines renders a single cell as cellH lines.
func (m model) renderCellLines(row, col int, value uint8) []string {
	selected := m.row == row && m.col == col
	fixed := m.isFixed(row, col)
	conflict := m.showConflicts && m.hasConflict(row, col)
	inBoxRow := row % m.set.boxRows
	inBoxCol := col % m.set.boxCols
	checker := (inBoxRow+inBoxCol)%2 == 0
	selectedValue := m.grid[idx(m.row, m.col, m.set.size)]
	isPeer := row == m.row || col == m.col
	sameNumber := selectedValue != 0 && value == selectedValue && !selected

	bg := bgBase1
	if !checker {
		bg = bgBase2
	}
	if isPeer && !selected {
		if checker {
			bg = bgPeer1
		} else {
			bg = bgPeer2
		}
	}
	if sameNumber {
		bg = bgSame
	}
	if selected {
		bg = bgSelected
		if m.pulse {
			bg = bgSelectedPulse
		}
	}
	if conflict {
		bg = bgConflict
	}

	fg := fgMuted
	bold := false
	if fixed {
		fg = fgFixed
		bold = true
	} else if value != 0 {
		fg = fgFilled
		bold = true
	}
	if conflict {
		fg = fgContrast
		bold = true
	}

	cellStyle := lipgloss.NewStyle().
		Width(cellW).
		Foreground(fg).
		Background(bg)

	if bold {
		cellStyle = cellStyle.Bold(true)
	}

	lines := make([]string, cellH)
	for i := 0; i < cellH; i++ {
		lines[i] = cellStyle.Render(strings.Repeat(" ", cellW))
	}

	if value == 0 {
		index := idx(row, col, m.set.size)
		if index < len(m.notes) && m.notes[index] != 0 {
			noteStyle := cellStyle.Foreground(fgNote)
			return renderNotes(m.notes[index], m.set.size, noteStyle)
		}
		return lines
	}

	label := fmt.Sprintf("%d", value)
	rowMid := (cellH / 2)
	colMid := (cellW / 2)
	start := colMid - (len(label) / 2)
	if start < 0 {
		start = 0
	}
	if start+len(label) > cellW {
		start = cellW - len(label)
	}

	content := strings.Repeat(" ", start) + label + strings.Repeat(" ", cellW-start-len(label))
	lines[rowMid] = cellStyle.Render(content)
	return lines
}

// renderNotes prints candidate notes inside a cell.
func renderNotes(notes uint16, size int, style lipgloss.Style) []string {
	tokens := make([]string, 0, size)
	for i := 1; i <= size; i++ {
		if notes&(1<<uint(i-1)) != 0 {
			tokens = append(tokens, fmt.Sprintf("%d", i))
		}
	}
	wrapped := wrapTokens(tokens, cellW)
	lines := make([]string, cellH)
	for i := 0; i < cellH; i++ {
		line := ""
		if i < len(wrapped) {
			line = wrapped[i]
		}
		lines[i] = style.Render(padRight(line, cellW))
	}
	return lines
}

// wrapTokens wraps tokens into lines of a fixed width.
func wrapTokens(tokens []string, width int) []string {
	if len(tokens) == 0 {
		return nil
	}
	lines := []string{}
	current := ""
	for _, token := range tokens {
		if current == "" {
			current = token
			continue
		}
		if len(current)+1+len(token) <= width {
			current += " " + token
		} else {
			lines = append(lines, current)
			current = token
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// padRight pads a string with spaces to a fixed width.
func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}

// boardWidth computes the board width in characters.
func boardWidth(size, boxCols int) int {
	gaps := (size - 1) * cellGapW
	boxExtra := (size/boxCols - 1) * (boxGapW - cellGapW)
	return size*cellW + gaps + boxExtra
}

// boardFrameWidth includes the outer border width.
func boardFrameWidth(size, boxCols int) int {
	return boardWidth(size, boxCols) + boardPadX*2 + 2
}

// blankLines returns empty padding lines of a given width.
func blankLines(count, width int) []string {
	if count <= 0 {
		return nil
	}
	line := strings.Repeat(" ", width)
	lines := make([]string, count)
	for i := 0; i < count; i++ {
		lines[i] = line
	}
	return lines
}
