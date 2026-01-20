package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hacktails/mini-sudoku-go/internal/sudoku"
)

// main launches the Bubble Tea program with the Sudoku model.
func main() {
	m := sudoku.NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("error:", err)
	}
}
