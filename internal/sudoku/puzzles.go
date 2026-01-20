package sudoku

import "strings"

// puzzle stores a Sudoku puzzle grid and its full solution.
type puzzle struct {
	puzzle   []uint8
	solution []uint8
}

// puzzleSet defines a board size and its box dimensions.
type puzzleSet struct {
	size    int
	boxRows int
	boxCols int
}

// difficulty represents a requested puzzle difficulty level.
type difficulty int

const (
	diffEasy difficulty = iota
	diffMedium
	diffHard
)

// puzzleSets is the size catalog used across the app.
var puzzleSets = map[int]puzzleSet{
	4: {size: 4, boxRows: 2, boxCols: 2},
	6: {size: 6, boxRows: 2, boxCols: 3},
	9: {size: 9, boxRows: 3, boxCols: 3},
}

// difficultyLabel returns the display label for a difficulty value.
func difficultyLabel(d difficulty) string {
	switch d {
	case diffEasy:
		return "Easy"
	case diffMedium:
		return "Medium"
	case diffHard:
		return "Hard"
	default:
		return "Easy"
	}
}

// parseDifficulty converts a string into a difficulty value.
func parseDifficulty(value string) difficulty {
	switch strings.ToLower(value) {
	case "easy":
		return diffEasy
	case "medium":
		return diffMedium
	case "hard":
		return diffHard
	default:
		return diffEasy
	}
}

// nextDifficulty cycles difficulty values.
func nextDifficulty(d difficulty) difficulty {
	switch d {
	case diffEasy:
		return diffMedium
	case diffMedium:
		return diffHard
	default:
		return diffEasy
	}
}

// clueCount returns target clue counts per size/difficulty.
func clueCount(size int, diff difficulty) int {
	switch size {
	case 4:
		switch diff {
		case diffEasy:
			return 10
		case diffMedium:
			return 8
		default:
			return 6
		}
	case 6:
		switch diff {
		case diffEasy:
			return 20
		case diffMedium:
			return 16
		default:
			return 12
		}
	default:
		switch diff {
		case diffEasy:
			return 36
		case diffMedium:
			return 30
		default:
			return 24
		}
	}
}
