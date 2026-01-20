# 🧩 Mini Sudoku Go

![Mini Sudoku Go](docs/mini-sudoku-go.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/hacktails/mini-sudoku-go)](https://goreportcard.com/report/github.com/hacktails/mini-sudoku-go)
[![License](https://img.shields.io/github/license/hacktails/mini-sudoku-go)](LICENSE)
[![Release](https://img.shields.io/github/v/release/hacktails/mini-sudoku-go)](https://github.com/hacktails/mini-sudoku-go/releases)
[![Build Status](https://github.com/hacktails/mini-sudoku-go/actions/workflows/release.yml/badge.svg)](https://github.com/hacktails/mini-sudoku-go/actions/workflows/release.yml)

**Your terminal just got a whole lot cozier.**

Mini Sudoku Go is a delightful, polished Sudoku experience right in your CLI. Built with the lovely [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss).

Whether you're killing time while your code compiles or you're a hardcore logic puzzle fan, we've got you covered with 4x4, 6x6, and classic 9x9 boards.

## ✨ Features

- **🎛 Flexible Boards:** Quick 4x4 snacks, 6x6 mid-sized meals, or the full 9x9 feast.
- **📝 Notes Mode:** Pencil in candidates like a pro.
- **💡 Smart Hints:** Stuck? We'll give you a logical nudge before spoiling the fun.
- **↩️ Undo/Redo:** Because everyone deserves a second chance (or third).
- **😈 Strict Mode:** Challenge yourself with a mistake limit. High stakes!
- **💾 Save Slots:** 3 slots to keep your progress safe.
- **🏆 Best Times:** Race against the clock and beat your personal bests.

## 🚀 How to Play

### 🍺 Homebrew (Preferred)

The easiest way to get started on macOS or Linux:

```bash
brew install hacktails/tap/mini-sudoku-go
```

### 🐹 Go Install

Alternatively, if you have Go installed:

```bash
go install github.com/hacktails/mini-sudoku-go/cmd/mini-sudoku-go@latest
```

### Run Locally

Clone and run instantly:

```bash
go run ./cmd/mini-sudoku-go
```

## 🎮 Controls

Navigate the grid and master the numbers:

| Action | Key |
| :--- | :--- |
| **Move** | Arrows or `h` `j` `k` `l` |
| **Enter Number** | `1`–`9` |
| **Toggle Notes** | `p` (Pencil) |
| **Clear Notes** | `c` |
| **Validate** | `v` |
| **Get Hint** | `H` |
| **Undo / Redo** | `u` / `y` |
| **Strict Mode** | `m` |
| **Change Size** | `s` then `4`, `6`, or `9` |
| **Difficulty** | `d` |
| **Save / Load** | `w` / `o` then `1`–`3` |
| **Help** | `?` |
| **Quit** | `q` |

## 🛠️ Development

Want to hack on the game? Awesome!

### Project Structure

- `cmd/mini-sudoku-go/`: The main entry point.
- `internal/sudoku/`: Where the magic happens (Game logic, UI, etc).
- `puzzles.json`: Our stash of curated brain-teasers.

## 📜 License

MIT © Hacktails — Hack away!
