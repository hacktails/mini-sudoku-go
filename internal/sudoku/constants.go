package sudoku

// Layout and persistence constants.
const (
	cellW     = 9
	cellH     = 3
	boardPadX = 2
	boardPadY = 1
	headerPad = 1
)

const (
	savesFile   = ".sudoku_saves.json"
	statsFile   = ".sudoku_stats.json"
	puzzlesFile = "puzzles.json"
)

const (
	maxMistakes = 3
	slotCount   = 3
)

const (
	cellGap     = ""
	boxGap      = "  "
	cellGapW    = 0
	boxGapW     = 2
	rowGapLines = 0
	boxGapLines = 1
)
