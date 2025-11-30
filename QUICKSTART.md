# Quick Start Guide

## Installation

```bash
# Clone the repository
git clone https://github.com/stefan-hacks/gbloxs.git
cd gbloxs

# Install dependencies
go mod download

# Build
go build -o gbloxs main.go

# Run
./gbloxs
```

Or use the Makefile:

```bash
make install  # Install dependencies
make build    # Build the application
make run      # Run the application
```

## First Steps

1. **Start the application**: Run `./gbloxs` or `make run`

2. **Navigate blocks**: Use `j` (down) and `k` (up) to move between blocks

3. **Expand/Collapse**: Press `e` or `Space` to expand/collapse the selected block

4. **Try input mode**: Press `i` to enter input mode, then:
   - Type a command like `/ls -la` to execute it
   - Or type regular text to create a text block
   - Press `Enter` to submit, `ESC` to cancel

5. **View help**: Press `h` to see all available keyboard shortcuts

6. **Try table view**: Press `t` to toggle the interactive table view

## Example Workflow

1. Start the app: `./gbloxs`
2. Press `i` to enter input mode
3. Type `/ps aux | head -10` and press Enter
4. A new block will appear with the command output
5. Navigate to it with `j` or `k`
6. Press `c` to copy the content to clipboard
7. Press `e` to expand/collapse and see more details
8. Press `h` to see all available commands

## Tips

- **Command execution**: Prefix commands with `/` or `!` in input mode
- **Block management**: Use `d` to delete, `r` to refresh, `c` to copy
- **Navigation**: Arrow keys work too (`↑`/`↓`)
- **Quick quit**: Press `q` or `Ctrl+C` to exit

## Troubleshooting

### Colors not showing?
Make sure your terminal supports ANSI colors. Try:
```bash
echo $TERM
# Should show something like xterm-256color
```

### Build errors?
Make sure you have Go 1.21+ installed:
```bash
go version
```

### Dependencies not found?
Run:
```bash
go mod tidy
go mod download
```

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Explore all keyboard shortcuts with `h` in the app
- Try different block types and features
- Customize the styling in `main.go`

Enjoy your interactive terminal experience! 🚀

