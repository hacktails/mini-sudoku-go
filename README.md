# Mini Sudoku Go

![Mini Sudoku Go](docs/mini-sudoku-go.png)

A beautiful terminal Sudoku in Go, built with Bubble Tea + Lip Gloss. Supports 4x4, 6x6, and 9x9 boards with notes, hints, undo/redo, strict mode, and save slots.

## Features

- 4x4 / 6x6 / 9x9 boards
- Notes (candidate marks)
- Smart hints (logic first, then reveal)
- Undo/redo
- Strict mode (mistake limit)
- Save slots (1–3)
- Best time tracking
- Curated puzzle library + generator fallback

## Install

### Homebrew (planned)

Once the tap is published, you will be able to install with:

```bash
brew install hacktails/tap/mini-sudoku-go
```

### Go

```bash
go install github.com/hacktails/mini-sudoku-go/cmd/mini-sudoku-go@latest
```

## Run locally

```bash
go run ./cmd/mini-sudoku-go
```

## Controls

- Move: Arrow keys or `h` `j` `k` `l`
- Enter: `1–9` (respects board size)
- Notes: `p` toggle, `c` clear notes in cell
- Validate: `v`
- Hint: `H`
- Undo/Redo: `u` / `y`
- Strict mode: `m`
- Size: `s` then `4` / `6` / `9`
- Difficulty: `d`
- Save/Load: `w` / `o` then slot `1–3`
- Help: `?`
- Quit: `q`

## Development

### Hot reload

```bash
go install github.com/air-verse/air@latest
air
```

### Formatting

```bash
gofmt -w .
```

## Project structure

```text
cmd/mini-sudoku-go/    Entrypoint
internal/sudoku/       Game logic, UI, persistence
puzzles.json           Curated puzzles
```

## Releases

Tag releases like `v0.1.0`. The GitHub Actions workflow builds binaries and attaches them to the release.

## License

MIT © Hacktails
