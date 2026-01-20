package sudoku

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// statusView renders the status panel or selection prompts.
func (m model) statusView() string {
	width := boardFrameWidth(m.set.size, m.set.boxCols)
	if m.selectingSize {
		title := statusTitleStyle.Render("Select size")
		body := statusTextStyle.Render("Press 4, 6, or 9 (Esc to cancel)")
		return statusBoxStyle.Width(width).Render(title + "\n" + body)
	}
	if m.selectingSlot {
		mode := "save"
		if m.slotMode == slotLoad {
			mode = "load"
		}
		title := statusTitleStyle.Render("Select slot")
		body := statusTextStyle.Render(fmt.Sprintf("Press 1-%d to %s (Esc to cancel)", slotCount, mode))
		return statusBoxStyle.Width(width).Render(title + "\n" + body)
	}
	if m.gameOver {
		title := statusDangerStyle.Render("Game over")
		body := statusTextStyle.Render("n new | r reset | d difficulty | s size | o load | q quit")
		return statusBoxStyle.Width(width).Render(title + "\n" + body)
	}
	if m.solved {
		title := statusSuccessStyle.Render("Solved")
		body := statusTextStyle.Render("n new | q quit")
		return statusBoxStyle.Width(width).Render(title + "\n" + body)
	}
	if m.flashMessage != "" && time.Now().Before(m.flashUntil) {
		msg := statusAccentStyle.Render(m.flashMessage)
		return statusBoxStyle.Width(width).Render(msg)
	}
	best := bestTimeString(m.stats, m.set.size, m.difficulty)
	statsLine := fmt.Sprintf("Mistakes %d/%d  Hints %d  Best %s  Slot %d", m.mistakes, maxMistakes, m.hintsUsed, best, m.activeSlot)
	statsLine = statusTextStyle.Render(statsLine)
	controlsLine := fmt.Sprintf(
		"Keys: 1-%d set  p notes  v validate  u/y undo  H hint  m strict  s size  d diff  w save  o load  ? help  q quit",
		m.set.size,
	)
	controlsLine = statusHintStyle.Render(controlsLine)
	lines := []string{statsLine, controlsLine}
	return statusBoxStyle.Width(width).Render(strings.Join(lines, "\n"))
}

// badge renders a pill-style label.
func badge(text string, style lipgloss.Style) string {
	return style.Render(" " + text + " ")
}

// toggleBadge renders an ON/OFF badge.
func toggleBadge(label string, on bool, onStyle, offStyle lipgloss.Style) string {
	if on {
		return badge(label+" ON", onStyle)
	}
	return badge(label+" OFF", offStyle)
}

// onOff returns a readable ON/OFF string.
func onOff(value bool) string {
	if value {
		return "ON"
	}
	return "OFF"
}

// statsKey returns a key used for best-time storage.
func statsKey(size int, diff difficulty) string {
	return fmt.Sprintf("%dx%d:%s", size, size, strings.ToLower(difficultyLabel(diff)))
}

// bestTimeString returns the formatted best time for size/difficulty.
func bestTimeString(st stats, size int, diff difficulty) string {
	key := statsKey(size, diff)
	best, ok := st.Best[key]
	if !ok || best <= 0 {
		return "--:--"
	}
	return formatSeconds(best)
}

// formatSeconds formats a duration in mm:ss.
func formatSeconds(seconds int64) string {
	mins := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}
