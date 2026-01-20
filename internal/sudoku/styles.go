package sudoku

import "github.com/charmbracelet/lipgloss"

// Palette and styles for the TUI.
var (
	bgBase1         = lipgloss.Color("#20262E")
	bgBase2         = lipgloss.Color("#1A1F26")
	bgPeer1         = lipgloss.Color("#24313C")
	bgPeer2         = lipgloss.Color("#1D2933")
	bgSame          = lipgloss.Color("#3B4C5C")
	bgSelected      = lipgloss.Color("#2F5D62")
	bgSelectedPulse = lipgloss.Color("#387175")
	bgConflict      = lipgloss.Color("#6B2F2F")

	fgFixed    = lipgloss.Color("#F2CC8F")
	fgFilled   = lipgloss.Color("#F4F1DE")
	fgMuted    = lipgloss.Color("#6C7A89")
	fgContrast = lipgloss.Color("#FFFFFF")
	fgNote     = lipgloss.Color("#8A9AA6")

	headerBarStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F4F1DE")).
			Background(lipgloss.Color("#2A9D8F")).
			Padding(0, headerPad)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B8C5C5"))

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B8C5C5"))

	badgePrimaryStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#0B1F1E")).
				Background(lipgloss.Color("#E9C46A"))

	badgeAccentStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F4F1DE")).
				Background(lipgloss.Color("#2A9D8F"))

	badgeOnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E7ECEF")).
			Background(lipgloss.Color("#3A5A40"))

	badgeOffStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9AA5B1")).
			Background(lipgloss.Color("#2B3036"))

	badgeWarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1B1B1B")).
			Background(lipgloss.Color("#E9C46A"))

	boardStyle = lipgloss.NewStyle().
			Padding(boardPadY, boardPadX).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3A4048"))

	statusBoxStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3A4048")).
			Foreground(lipgloss.Color("#C7D2DA"))

	statusTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#E9C46A"))

	statusTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C7D2DA"))

	statusHintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A96A3"))

	statusInfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#B8C5C5"))

	statusAccentStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F4F1DE"))

	statusSuccessStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#9AE6B4"))

	statusDangerStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F28482"))

	helpBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3A4048")).
			Padding(1, 2).
			Background(lipgloss.Color("#1A1F26")).
			Foreground(lipgloss.Color("#E7ECEF"))

	helpTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#E9C46A"))

	helpBodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D8DEE9"))

	helpFootStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9AA5B1"))
)
