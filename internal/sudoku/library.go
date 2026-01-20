package sudoku

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// puzzleEntry is a JSON entry for a curated puzzle.
type puzzleEntry struct {
	Size       int     `json:"size"`
	Difficulty string  `json:"difficulty"`
	Puzzle     []uint8 `json:"puzzle"`
	Solution   []uint8 `json:"solution"`
}

// puzzleLibrary is the JSON payload for curated puzzles.
type puzzleLibrary struct {
	Puzzles []puzzleEntry `json:"puzzles"`
}

var (
	libraryLoaded bool
	libraryByKey  map[string][]puzzle
)

// randomFromLibrary returns a puzzle from the curated library when available.
func randomFromLibrary(set puzzleSet, diff difficulty) (puzzle, bool) {
	loadLibrary()
	key := libraryKey(set.size, diff)
	list := libraryByKey[key]
	if len(list) == 0 {
		return puzzle{}, false
	}
	return list[rng.Intn(len(list))], true
}

// libraryKey creates the lookup key used for curated puzzles.
func libraryKey(size int, diff difficulty) string {
	return fmt.Sprintf("%dx%d:%s", size, size, strings.ToLower(difficultyLabel(diff)))
}

// loadLibrary reads puzzles.json and indexes valid entries by size/difficulty.
func loadLibrary() {
	if libraryLoaded {
		return
	}
	libraryLoaded = true
	libraryByKey = map[string][]puzzle{}

	data, err := os.ReadFile(puzzlesFile)
	if err != nil {
		return
	}
	var lib puzzleLibrary
	if err := json.Unmarshal(data, &lib); err != nil {
		return
	}
	for _, entry := range lib.Puzzles {
		set, ok := puzzleSets[entry.Size]
		if !ok {
			continue
		}
		diff := parseDifficulty(entry.Difficulty)
		if !validEntry(entry, set) {
			continue
		}
		if countSolutions(entry.Puzzle, set, 2) != 1 {
			continue
		}
		key := libraryKey(entry.Size, diff)
		libraryByKey[key] = append(libraryByKey[key], puzzle{
			puzzle:   entry.Puzzle,
			solution: entry.Solution,
		})
	}
}

// validEntry ensures the puzzle/solution are consistent with size constraints.
func validEntry(entry puzzleEntry, set puzzleSet) bool {
	expected := set.size * set.size
	if len(entry.Puzzle) != expected || len(entry.Solution) != expected {
		return false
	}
	for i, value := range entry.Puzzle {
		if value == 0 {
			continue
		}
		if value > uint8(set.size) || entry.Solution[i] != value {
			return false
		}
	}
	for _, value := range entry.Solution {
		if value == 0 || value > uint8(set.size) {
			return false
		}
	}
	return true
}
