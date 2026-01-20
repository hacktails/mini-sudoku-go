package sudoku

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// View renders the entire game screen.
func (m model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	elapsed := time.Since(m.start)
	if m.solved && m.elapsedAtSolve > 0 {
		elapsed = time.Duration(m.elapsedAtSolve) * time.Second
	}
	mins := int(elapsed.Minutes())
	secs := int(elapsed.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d", mins, secs)

	header := m.headerView(timeStr)
	board := boardStyle.Render(m.boardView())
	status := m.statusView()

	content := lipgloss.JoinVertical(lipgloss.Center, header, "", board, "", status)

	base := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
	if m.showHelp {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.helpView())
	}
	return base
}

// headerView renders the title, subtitle, and meta badges.
func (m model) headerView(timeStr string) string {
	width := boardFrameWidth(m.set.size, m.set.boxCols)
	line := alignLeftRight(
		"Mini Sudoku",
		timeStr,
		width-(2*headerPad),
	)
	bar := headerBarStyle.Render(line)

	subtitle := fmt.Sprintf(
		"Fill each row, column, and %dx%d box with 1-%d",
		m.set.boxRows,
		m.set.boxCols,
		m.set.size,
	)
	info := subtitleStyle.Width(width).Align(lipgloss.Center).Render(subtitle)
	meta := metaStyle.Width(width).Align(lipgloss.Center).Render(m.metaView())

	return lipgloss.JoinVertical(lipgloss.Center, bar, info, meta)
}

// alignLeftRight pads a string so left and right content fit a fixed width.
func alignLeftRight(left, right string, width int) string {
	space := width - lipgloss.Width(left) - lipgloss.Width(right)
	if space < 1 {
		space = 1
	}
	return left + strings.Repeat(" ", space) + right
}

// metaView renders the header badges.
func (m model) metaView() string {
	sizeBadge := badge(fmt.Sprintf("Size %dx%d", m.set.size, m.set.size), badgeAccentStyle)
	diffBadge := badge("Diff "+strings.ToUpper(difficultyLabel(m.difficulty)), badgePrimaryStyle)
	notesBadge := toggleBadge("Notes", m.noteMode, badgeOnStyle, badgeOffStyle)
	validateBadge := toggleBadge("Validate", m.showConflicts, badgeOnStyle, badgeOffStyle)
	strictBadge := toggleBadge("Strict", m.strictMode, badgeWarnStyle, badgeOffStyle)
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		sizeBadge,
		" ",
		diffBadge,
		" ",
		notesBadge,
		" ",
		validateBadge,
		" ",
		strictBadge,
	)
}

// helpView shows a quick keyboard reference overlay.
func (m model) helpView() string {
	lines := []string{
		"Navigation: arrows or hjkl",
		"Numbers: 1-9 set value (respecting size)",
		"Notes: p toggle notes, c clear notes in cell",
		"Undo/Redo: u / y",
		"Hint: H (smart hint first, then reveal)",
		"Validate: v",
		"Strict mode: m",
		"Size: s then 4/6/9",
		"Difficulty: d",
		"Save/Load slots: w / o, then 1-3",
		"Quit: q",
	}
	body := strings.Join(lines, "\n")
	content := helpTitleStyle.Render("Help") + "\n\n" + helpBodyStyle.Render(body) + "\n\n" + helpFootStyle.Render("Press ? or Esc to close")
	return helpBoxStyle.Render(content)
}
