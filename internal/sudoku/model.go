package sudoku

import "time"

// model is the Bubble Tea state container for the game.
type model struct {
	set            puzzleSet
	puzzle         puzzle
	grid           []uint8
	notes          []uint16
	row            int
	col            int
	start          time.Time
	width          int
	height         int
	difficulty     difficulty
	mistakes       int
	hintsUsed      int
	noteMode       bool
	showConflicts  bool
	selectingSize  bool
	showHelp       bool
	solved         bool
	elapsedAtSolve int64
	strictMode     bool
	gameOver       bool
	selectingSlot  bool
	slotMode       slotMode
	activeSlot     int
	pulse          bool
	undoStack      []snapshot
	redoStack      []snapshot
	stats          stats
	flashMessage   string
	flashUntil     time.Time
}

// slotMode indicates whether the slot prompt is saving or loading.
type slotMode int

const (
	slotSave slotMode = iota
	slotLoad
)

// snapshot stores undo/redo state.
type snapshot struct {
	grid      []uint8
	notes     []uint16
	row       int
	col       int
	mistakes  int
	hintsUsed int
}

// stats tracks best times per size/difficulty.
type stats struct {
	Best map[string]int64 `json:"best"`
}

// saveState serializes a single game state to disk.
type saveState struct {
	Size          int      `json:"size"`
	BoxRows       int      `json:"box_rows"`
	BoxCols       int      `json:"box_cols"`
	Difficulty    string   `json:"difficulty"`
	Puzzle        []uint8  `json:"puzzle"`
	Solution      []uint8  `json:"solution"`
	Grid          []uint8  `json:"grid"`
	Notes         []uint16 `json:"notes"`
	Row           int      `json:"row"`
	Col           int      `json:"col"`
	StartUnix     int64    `json:"start_unix"`
	Mistakes      int      `json:"mistakes"`
	HintsUsed     int      `json:"hints_used"`
	NoteMode      bool     `json:"note_mode"`
	ShowConflicts bool     `json:"show_conflicts"`
	Solved        bool     `json:"solved"`
	Elapsed       int64    `json:"elapsed"`
	StrictMode    bool     `json:"strict_mode"`
	GameOver      bool     `json:"game_over"`
}

// saveSlots stores all slot saves in one file.
type saveSlots struct {
	Active int               `json:"active"`
	Slots  map[int]saveState `json:"slots"`
}

// tickMsg drives timer updates.
type tickMsg time.Time

// NewModel constructs the initial game model (loads save if present).
func NewModel() model {
	st := loadStats()
	if saved, slot, ok := loadActiveSave(); ok {
		return modelFromSave(saved, st, slot)
	}

	set := puzzleSets[6]
	diff := diffEasy
	p := generatePuzzle(set, diff)
	return model{
		set:           set,
		puzzle:        p,
		grid:          copyGrid(p.puzzle),
		notes:         make([]uint16, set.size*set.size),
		row:           0,
		col:           0,
		start:         time.Now(),
		difficulty:    diff,
		showConflicts: true,
		activeSlot:    1,
		stats:         st,
	}
}
