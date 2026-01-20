package sudoku

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Init starts the periodic tick for UI updates.
func (m model) Init() tea.Cmd {
	return tickCmd()
}

// tickCmd emits a tick every second.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles Bubble Tea messages and keyboard input.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		m.pulse = !m.pulse
		return m, tickCmd()
	case tea.KeyMsg:
		if m.selectingSlot {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc":
				m.selectingSlot = false
				return m, nil
			case "1", "2", "3":
				slot := int(msg.Runes[0] - '0')
				if slot >= 1 && slot <= slotCount {
					m.activeSlot = slot
					if m.slotMode == slotSave {
						if err := m.save(); err == nil {
							m.flash(fmt.Sprintf("Saved to slot %d", slot))
						} else {
							m.flash("Save failed")
						}
					} else {
						if saved, ok := loadSlot(slot); ok {
							m = modelFromSave(saved, m.stats, slot)
							m.flash(fmt.Sprintf("Loaded slot %d", slot))
						} else {
							m.flash("Empty slot")
						}
					}
				}
				m.selectingSlot = false
				return m, nil
			default:
				return m, nil
			}
		}

		if m.showHelp {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc", "?":
				m.showHelp = false
				return m, nil
			default:
				return m, nil
			}
		}

		if m.selectingSize {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc", "s":
				m.selectingSize = false
				return m, nil
			case "4", "6", "9":
				m.setSize(int(msg.Runes[0] - '0'))
				m.selectingSize = false
				return m, nil
			default:
				return m, nil
			}
		}

		if m.gameOver {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "n":
				m.newPuzzle()
				return m, nil
			case "r":
				m.reset()
				return m, nil
			case "s":
				m.selectingSize = true
				return m, nil
			case "d":
				m.setDifficulty(nextDifficulty(m.difficulty))
				return m, nil
			case "o":
				m.slotMode = slotLoad
				m.selectingSlot = true
				return m, nil
			case "?":
				m.showHelp = true
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "n":
			m.newPuzzle()
			return m, nil
		case "r":
			m.reset()
			return m, nil
		case "s":
			m.selectingSize = true
			return m, nil
		case "m":
			m.strictMode = !m.strictMode
			if m.strictMode {
				m.flash(fmt.Sprintf("Strict mode (max %d mistakes)", maxMistakes))
			} else {
				m.flash("Strict mode off")
			}
			return m, nil
		case "d":
			m.setDifficulty(nextDifficulty(m.difficulty))
			return m, nil
		case "p":
			m.noteMode = !m.noteMode
			if m.noteMode {
				m.flash("Notes mode")
			} else {
				m.flash("Entry mode")
			}
			return m, nil
		case "v":
			m.showConflicts = !m.showConflicts
			if m.showConflicts {
				m.flash("Validation on")
			} else {
				m.flash("Validation off")
			}
			return m, nil
		case "u":
			m.undo()
			return m, nil
		case "y":
			m.redo()
			return m, nil
		case "c":
			m.clearNotes(m.row, m.col)
			return m, nil
		case "w":
			m.slotMode = slotSave
			m.selectingSlot = true
			return m, nil
		case "o":
			m.slotMode = slotLoad
			m.selectingSlot = true
			return m, nil
		case "?":
			m.showHelp = true
			return m, nil
		case "left", "h":
			m.move(0, -1)
			return m, nil
		case "right", "l":
			m.move(0, 1)
			return m, nil
		case "up", "k":
			m.move(-1, 0)
			return m, nil
		case "down", "j":
			m.move(1, 0)
			return m, nil
		case "backspace", "delete", " ", "space":
			m.clearValue()
			return m, nil
		}

		if len(msg.Runes) == 1 {
			r := msg.Runes[0]
			if r == '?' {
				m.showHelp = true
				return m, nil
			}
			if r == 'H' {
				m.applyHint()
				return m, nil
			}
			if r >= '1' && r <= '9' {
				value := int(r - '0')
				if value <= m.set.size {
					if m.noteMode {
						m.toggleNote(m.row, m.col, value)
					} else {
						m.setValue(uint8(value))
					}
					return m, nil
				}
			}
		}
	}

	return m, nil
}
