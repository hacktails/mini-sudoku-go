package sudoku

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// save writes the current game state to the active slot.
func (m *model) save() error {
	state := m.makeSaveState()
	return saveToSlot(m.activeSlot, state)
}

// makeSaveState serializes the current model into a saveState.
func (m *model) makeSaveState() saveState {
	return saveState{
		Size:          m.set.size,
		BoxRows:       m.set.boxRows,
		BoxCols:       m.set.boxCols,
		Difficulty:    strings.ToLower(difficultyLabel(m.difficulty)),
		Puzzle:        copyGrid(m.puzzle.puzzle),
		Solution:      copyGrid(m.puzzle.solution),
		Grid:          copyGrid(m.grid),
		Notes:         copyNotes(m.notes),
		Row:           m.row,
		Col:           m.col,
		StartUnix:     m.start.Unix(),
		Mistakes:      m.mistakes,
		HintsUsed:     m.hintsUsed,
		NoteMode:      m.noteMode,
		ShowConflicts: m.showConflicts,
		Solved:        m.solved,
		Elapsed:       m.elapsedAtSolve,
		StrictMode:    m.strictMode,
		GameOver:      m.gameOver,
	}
}

// loadActiveSave returns the active slot's save if present.
func loadActiveSave() (saveState, int, bool) {
	slots := loadSlotsFile()
	slot := slots.Active
	if slot == 0 {
		slot = 1
	}
	state, ok := slots.Slots[slot]
	return state, slot, ok
}

// loadSlot loads a specific save slot.
func loadSlot(slot int) (saveState, bool) {
	slots := loadSlotsFile()
	state, ok := slots.Slots[slot]
	return state, ok
}

// saveToSlot persists a saveState into a slot.
func saveToSlot(slot int, state saveState) error {
	slots := loadSlotsFile()
	if slots.Slots == nil {
		slots.Slots = map[int]saveState{}
	}
	slots.Active = slot
	slots.Slots[slot] = state
	return saveSlotsFile(slots)
}

// loadSlotsFile reads the save slots file from disk.
func loadSlotsFile() saveSlots {
	data, err := os.ReadFile(savesFile)
	if err != nil {
		return saveSlots{Active: 1, Slots: map[int]saveState{}}
	}
	var slots saveSlots
	if err := json.Unmarshal(data, &slots); err != nil {
		return saveSlots{Active: 1, Slots: map[int]saveState{}}
	}
	if slots.Active == 0 {
		slots.Active = 1
	}
	if slots.Slots == nil {
		slots.Slots = map[int]saveState{}
	}
	return slots
}

// saveSlotsFile writes save slots to disk.
func saveSlotsFile(slots saveSlots) error {
	data, err := json.MarshalIndent(slots, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(savesFile, data, 0o644)
}

// modelFromSave reconstructs a model from a saved state.
func modelFromSave(state saveState, st stats, slot int) model {
	set, ok := puzzleSets[state.Size]
	if !ok {
		set = puzzleSets[6]
	}
	diff := parseDifficulty(state.Difficulty)
	m := model{
		set:            set,
		puzzle:         puzzle{puzzle: state.Puzzle, solution: state.Solution},
		grid:           state.Grid,
		notes:          state.Notes,
		row:            state.Row,
		col:            state.Col,
		start:          time.Unix(state.StartUnix, 0),
		difficulty:     diff,
		mistakes:       state.Mistakes,
		hintsUsed:      state.HintsUsed,
		noteMode:       state.NoteMode,
		showConflicts:  state.ShowConflicts,
		solved:         state.Solved,
		elapsedAtSolve: state.Elapsed,
		strictMode:     state.StrictMode,
		gameOver:       state.GameOver,
		activeSlot:     slot,
		stats:          st,
	}
	expected := set.size * set.size
	if len(m.grid) != expected {
		m.grid = copyGrid(m.puzzle.puzzle)
	}
	if len(m.notes) != expected {
		m.notes = make([]uint16, expected)
	}
	if (m.solved || m.gameOver) && m.elapsedAtSolve == 0 {
		m.elapsedAtSolve = int64(time.Since(m.start).Seconds())
	}
	m.solved = m.solved || m.isSolved()
	return m
}

// loadStats loads best-time stats.
func loadStats() stats {
	st := stats{Best: map[string]int64{}}
	data, err := os.ReadFile(statsFile)
	if err != nil {
		return st
	}
	if err := json.Unmarshal(data, &st); err != nil {
		return st
	}
	if st.Best == nil {
		st.Best = map[string]int64{}
	}
	return st
}

// saveStats writes best-time stats to disk.
func saveStats(st stats) error {
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statsFile, data, 0o644)
}
