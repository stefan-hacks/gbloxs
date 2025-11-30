# Gbloxs - Interactive Terminal Blocks

A powerful, interactive terminal application inspired by Warp Terminal, built with Charmbracelet's Bubble Tea, Lip Gloss, and Bubbles frameworks. Gbloxs provides a modern, block-styled terminal interface with full interactivity, colorization, and advanced features.

## Features

### 🎨 Visual Features
- **Block-Styled Output**: Organize terminal output into interactive, visually distinct blocks
- **Colorized Output**: Automatic syntax highlighting and color coding
- **Multiple Block Types**: Command, Output, Error, Success, Info, Progress, and Table blocks
- **Beautiful Borders**: Rounded borders with color-coded types
- **Responsive Layout**: Adapts to terminal size automatically

### 🖱️ Interactive Features
- **Keyboard Navigation**: Navigate between blocks with `j`/`k` or arrow keys
- **Expand/Collapse**: Toggle block expansion to show/hide content
- **Command Execution**: Execute shell commands directly from blocks
- **Input Mode**: Interactive command input with syntax highlighting
- **Table View**: Interactive table component with navigation
- **Progress Indicators**: Animated progress bars and spinners
- **Viewport Scrolling**: Scroll through long content within blocks

### 🚀 Advanced Features
- **Syntax Highlighting**: Automatic highlighting for:
  - Directory listings
  - File permissions
  - Error messages
  - Success indicators
  - File paths
  - Numbers
- **Command History**: Track executed commands with timestamps
- **Block Management**: Copy, refresh, and delete blocks
- **Help System**: Built-in help overlay with all shortcuts
- **Real-time Updates**: Live progress indicators and status updates

## Installation

### Prerequisites
- Go 1.21 or later
- A terminal with ANSI color support

### Build from Source

```bash
git clone https://github.com/stefan-hacks/gbloxs.git
cd gbloxs
go mod download
go build -o gbloxs
./gbloxs
```

Or run directly:

```bash
go run main.go
```

## Usage

### Basic Navigation

```
j / ↓     Navigate down to next block
k / ↑     Navigate up to previous block
e         Expand/collapse selected block
Space     Toggle block expansion
Enter     Toggle block expansion
```

### Block Actions

```
c         Copy block content
r         Refresh/reload block
d         Delete selected block
x         Execute command in selected block
```

### Modes

```
i         Toggle input mode
h         Toggle help overlay
t         Toggle table view
```

### Input Mode

When in input mode:
- Type commands or text
- Use `/cmd` or `!cmd` prefix to execute shell commands
- Press `Enter` to submit
- Press `ESC` to cancel

### General

```
q         Quit application
Ctrl+C    Quit application
Ctrl+L    Clear all blocks
```

## Block Types

### 🟡 Command Blocks
Yellow border - Display command execution with input and output

### 🟢 Output/Success Blocks
Green border - Show successful command output or success messages

### 🔴 Error Blocks
Red border - Display error messages and failed operations

### 🔵 Info Blocks
Blue border - Show informational content and help text

### 📊 Table Blocks
Display tabular data with interactive navigation

### ⏳ Progress Blocks
Show animated progress bars and loading indicators

## Examples

### Creating Blocks

Blocks are automatically created when:
- Entering input mode and submitting commands
- Executing commands with `/cmd` or `!cmd` prefix
- Programmatically adding blocks

### Command Execution

```bash
# In input mode, type:
/ls -la

# Or:
!ps aux | grep nginx
```

### Interactive Tables

Press `t` to toggle the interactive table view. Navigate with arrow keys when the table is focused.

## Architecture

Gbloxs is built using:

- **Bubble Tea**: TUI framework for state management and user interaction
- **Lip Gloss**: Terminal styling library for colors, borders, and layouts
- **Bubbles**: Collection of reusable components (spinner, progress, table, textinput, viewport)

### Key Components

- **Block Model**: Represents individual interactive blocks
- **Main Model**: Manages application state and all blocks
- **Styles**: Centralized styling system using Lip Gloss
- **Viewport**: Handles scrolling for long content
- **Progress**: Animated progress indicators
- **Table**: Interactive table component

## Customization

### Styling

Modify the `Styles` struct in `main.go` to customize:
- Block borders and colors
- Text colors and formatting
- Table appearance
- Progress bar styles

### Block Types

Add new block types by:
1. Adding a new `BlockType` constant
2. Adding rendering logic in `renderBlock()`
3. Adding appropriate styling

## Keyboard Shortcuts Reference

| Key | Action |
|-----|--------|
| `j` / `↓` | Navigate down |
| `k` / `↑` | Navigate up |
| `e` | Expand/collapse block |
| `c` | Copy block content |
| `r` | Refresh block |
| `d` | Delete block |
| `x` | Execute command |
| `i` | Toggle input mode |
| `h` | Toggle help |
| `t` | Toggle table view |
| `q` / `Ctrl+C` | Quit |
| `Ctrl+L` | Clear all blocks |
| `Space` / `Enter` | Toggle expansion |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- [Charmbracelet](https://github.com/charmbracelet) for the amazing TUI frameworks
- [Warp Terminal](https://www.warp.dev/) for inspiration on block-styled terminals
- The Go community for excellent tooling

## Roadmap

- [ ] Clipboard integration for copy functionality
- [ ] Block templates and presets
- [ ] Custom themes and color schemes
- [ ] Plugin system for custom block types
- [ ] Command history and autocomplete
- [ ] Multi-select blocks
- [ ] Block grouping and nesting
- [ ] Export blocks to files
- [ ] Integration with external tools
- [ ] Mouse support for block interaction

## Screenshots

*Note: Run the application to see the beautiful interactive terminal interface!*

---

Made with ❤️ using [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and [Bubbles](https://github.com/charmbracelet/bubbles)

