package sudoku

import (
	"math/rand"
	"time"
)

// rng is the shared random source for puzzle generation and hints.
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))
