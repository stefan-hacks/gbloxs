package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Block represents an interactive block in the terminal
type Block struct {
	ID        string
	Title     string
	Content   string
	Type      BlockType
	Expanded  bool
	Selected  bool
	Progress  float64
	IsLoading bool
	Metadata  map[string]string
	Timestamp time.Time
	Command   string
	Output    string
	Error     string
	TableData [][]string
	Viewport  viewport.Model
}

type BlockType string

const (
	BlockTypeCommand  BlockType = "command"
	BlockTypeOutput   BlockType = "output"
	BlockTypeTable    BlockType = "table"
	BlockTypeProgress BlockType = "progress"
	BlockTypeInfo     BlockType = "info"
	BlockTypeError    BlockType = "error"
	BlockTypeSuccess  BlockType = "success"
)

type model struct {
	blocks      []Block
	selectedIdx int
	width       int
	height      int
	spinner     spinner.Model
	progress    progress.Model
	textInput   textinput.Model
	showInput   bool
	inputMode   bool
	styles      Styles
	table       table.Model
	showTable   bool
	helpMode    bool
}

type Styles struct {
	BlockBorder       lipgloss.Style
	BlockTitle        lipgloss.Style
	BlockContent      lipgloss.Style
	SelectedBlock     lipgloss.Style
	CommandBlock      lipgloss.Style
	OutputBlock       lipgloss.Style
	ErrorBlock        lipgloss.Style
	SuccessBlock      lipgloss.Style
	InfoBlock         lipgloss.Style
	ProgressBar       lipgloss.Style
	TableHeader       lipgloss.Style
	TableCell         lipgloss.Style
	TableSelectedCell lipgloss.Style
}

func NewStyles() Styles {
	return Styles{
		BlockBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2),

		BlockTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			MarginBottom(1),

		BlockContent: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			MarginTop(1),

		SelectedBlock: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(1, 2),

		CommandBlock: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("220")).
			Padding(1, 2),

		OutputBlock: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("34")).
			Padding(1, 2),

		ErrorBlock: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("196")).
			Padding(1, 2),

		SuccessBlock: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("46")).
			Padding(1, 2),

		InfoBlock: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(1, 2),

		ProgressBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")),

		TableHeader: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1),

		TableCell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Padding(0, 1),

		TableSelectedCell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true).
			Padding(0, 1),
	}
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	p := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	p.Width = 40

	ti := textinput.New()
	ti.Placeholder = "Enter command or text..."
	ti.CharLimit = 500
	ti.Width = 50
	ti.Focus()

	styles := NewStyles()

	// Create some example blocks
	blocks := []Block{
		{
			ID:        "1",
			Title:     "Command Execution",
			Command:   "ls -la",
			Output:    "total 48\ndrwxr-xr-x  8 user user  4096 Jan 15 10:30 .\ndrwxr-xr-x 18 user user  4096 Jan 10 09:15 ..\n-rw-r--r--  1 user user  1024 Jan 15 10:25 file.txt",
			Type:      BlockTypeCommand,
			Expanded:  true,
			Selected:  false,
			Timestamp: time.Now(),
		},
		{
			ID:        "2",
			Title:     "System Information",
			Content:   "OS: Linux\nKernel: 6.12.57\nArchitecture: amd64\nUptime: 5 days, 3 hours",
			Type:      BlockTypeInfo,
			Expanded:  true,
			Selected:  false,
			Timestamp: time.Now(),
		},
		{
			ID:        "3",
			Title:     "Progress Indicator",
			Type:      BlockTypeProgress,
			Expanded:  true,
			Selected:  false,
			Progress:  0.65,
			IsLoading: true,
			Timestamp: time.Now(),
		},
		{
			ID:       "4",
			Title:    "Data Table",
			Type:     BlockTypeTable,
			Expanded: true,
			Selected: false,
			TableData: [][]string{
				{"Name", "Status", "CPU %", "Memory %"},
				{"nginx", "Running", "2.5", "15.3"},
				{"postgres", "Running", "1.2", "45.8"},
				{"redis", "Running", "0.8", "12.1"},
			},
			Timestamp: time.Now(),
		},
		{
			ID:        "5",
			Title:     "Success Message",
			Content:   "✓ Operation completed successfully!\n✓ All checks passed\n✓ System is healthy",
			Type:      BlockTypeSuccess,
			Expanded:  true,
			Selected:  false,
			Timestamp: time.Now(),
		},
	}

	// Initialize viewports for blocks that need scrolling
	for i := range blocks {
		vp := viewport.New(50, 10)
		vp.SetContent(blocks[i].Output)
		blocks[i].Viewport = vp
	}

	if len(blocks) > 0 {
		blocks[0].Selected = true
	}

	// Initialize table
	columns := []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Status", Width: 12},
		{Title: "CPU %", Width: 10},
		{Title: "Memory %", Width: 12},
	}

	rows := []table.Row{
		{"nginx", "Running", "2.5", "15.3"},
		{"postgres", "Running", "1.2", "45.8"},
		{"redis", "Running", "0.8", "12.1"},
		{"node", "Running", "5.1", "28.4"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(tableStyles)

	return model{
		blocks:      blocks,
		selectedIdx: 0,
		spinner:     s,
		progress:    p,
		textInput:   ti,
		showInput:   false,
		inputMode:   false,
		styles:      styles,
		table:       t,
		showTable:   false,
		helpMode:    false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		animateProgress(m.blocks[2]),
	)
}

type progressMsg struct {
	blockID string
	value   float64
}

func animateProgress(block Block) tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return progressMsg{blockID: block.ID, value: block.Progress}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 20
		m.textInput.Width = msg.Width - 10

		// Update viewports
		for i := range m.blocks {
			vp := m.blocks[i].Viewport
			vp.Width = msg.Width - 10
			vp.Height = 15
			m.blocks[i].Viewport = vp
		}

	case tea.KeyMsg:
		if m.inputMode {
			switch msg.String() {
			case "esc":
				m.inputMode = false
				m.showInput = false
				m.textInput.Blur()
			case "enter":
				// Process input
				input := m.textInput.Value()
				if input != "" {
					m.addBlockFromInput(input)
					m.textInput.SetValue("")
				}
				m.inputMode = false
				m.showInput = false
				m.textInput.Blur()
			default:
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "i", "I":
			// Toggle input mode
			m.inputMode = !m.inputMode
			m.showInput = !m.showInput
			if m.inputMode {
				m.textInput.Focus()
			} else {
				m.textInput.Blur()
			}

		case "j", "down":
			if m.selectedIdx < len(m.blocks)-1 {
				m.blocks[m.selectedIdx].Selected = false
				m.selectedIdx++
				m.blocks[m.selectedIdx].Selected = true
			}

		case "k", "up":
			if m.selectedIdx > 0 {
				m.blocks[m.selectedIdx].Selected = false
				m.selectedIdx--
				m.blocks[m.selectedIdx].Selected = true
			}

		case "e", "E":
			// Expand/collapse block
			m.blocks[m.selectedIdx].Expanded = !m.blocks[m.selectedIdx].Expanded

		case "c", "C":
			// Copy block content to clipboard
			if m.selectedIdx < len(m.blocks) {
				block := m.blocks[m.selectedIdx]
				contentToCopy := block.Output
				if contentToCopy == "" {
					contentToCopy = block.Content
				}
				if contentToCopy == "" && block.Command != "" {
					contentToCopy = block.Command
				}
				if contentToCopy != "" {
					clipboard.WriteAll(contentToCopy)
					m.blocks[m.selectedIdx].Metadata["copied"] = "true"
					// Show feedback message
					m.addInfoBlock("Content copied to clipboard!")
				}
			}

		case "r", "R":
			// Refresh/reload block
			if m.blocks[m.selectedIdx].Type == BlockTypeProgress {
				m.blocks[m.selectedIdx].Progress = 0
				m.blocks[m.selectedIdx].IsLoading = true
				cmds = append(cmds, animateProgress(m.blocks[m.selectedIdx]))
			}

		case "d", "D":
			// Delete block
			if len(m.blocks) > 1 {
				m.blocks = append(m.blocks[:m.selectedIdx], m.blocks[m.selectedIdx+1:]...)
				if m.selectedIdx >= len(m.blocks) {
					m.selectedIdx = len(m.blocks) - 1
				}
				if len(m.blocks) > 0 {
					m.blocks[m.selectedIdx].Selected = true
				}
			}

		case " ", "enter":
			// Toggle selection or execute action
			m.blocks[m.selectedIdx].Expanded = !m.blocks[m.selectedIdx].Expanded

		case "h", "H":
			// Toggle help mode
			m.helpMode = !m.helpMode

		case "t", "T":
			// Toggle table view
			m.showTable = !m.showTable

		case "ctrl+l":
			// Clear all blocks
			m.blocks = []Block{}
			m.selectedIdx = 0

		case "x", "X":
			// Execute command in selected block
			if m.blocks[m.selectedIdx].Command != "" {
				m.executeCommand(m.blocks[m.selectedIdx].Command)
			}
		}

		// Handle table navigation when table is shown
		if m.showTable {
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			cmds = append(cmds, cmd)
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progressMsg:
		for i := range m.blocks {
			if m.blocks[i].ID == msg.blockID && m.blocks[i].IsLoading {
				m.blocks[i].Progress = msg.value
				if msg.value >= 1.0 {
					m.blocks[i].IsLoading = false
				} else {
					m.blocks[i].Progress += 0.01
					if m.blocks[i].Progress > 1.0 {
						m.blocks[i].Progress = 1.0
					}
					cmds = append(cmds, animateProgress(m.blocks[i]))
				}
			}
		}
	}

	// Update viewports for scrolling
	for i := range m.blocks {
		if m.blocks[i].Expanded {
			var cmd tea.Cmd
			m.blocks[i].Viewport, cmd = m.blocks[i].Viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) addBlockFromInput(input string) {
	newBlock := Block{
		ID:        fmt.Sprintf("%d", len(m.blocks)+1),
		Title:     "User Input",
		Content:   input,
		Command:   input,
		Type:      BlockTypeCommand,
		Expanded:  true,
		Selected:  false,
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}

	// Try to execute as command if it looks like one
	if strings.HasPrefix(input, "/") || strings.HasPrefix(input, "!") {
		cmdStr := strings.TrimPrefix(strings.TrimPrefix(input, "/"), "!")
		m.executeCommandInBlock(cmdStr, &newBlock)
	} else {
		// Simulate command output
		if strings.HasPrefix(input, "ls") {
			newBlock.Output = "file1.txt\nfile2.txt\nfile3.txt\ndirectory1\ndirectory2"
			newBlock.Type = BlockTypeSuccess
		} else if strings.HasPrefix(input, "error") {
			newBlock.Error = "Error: Command failed"
			newBlock.Type = BlockTypeError
		} else {
			newBlock.Output = fmt.Sprintf("Executed: %s\nStatus: OK", input)
			newBlock.Type = BlockTypeSuccess
		}
	}

	vp := viewport.New(m.width-10, 10)
	vp.SetContent(newBlock.Output)
	newBlock.Viewport = vp

	// Deselect all and select new block
	for i := range m.blocks {
		m.blocks[i].Selected = false
	}
	m.blocks = append(m.blocks, newBlock)
	m.selectedIdx = len(m.blocks) - 1
	m.blocks[m.selectedIdx].Selected = true
}

func (m model) executeCommandInBlock(cmdStr string, block *Block) {
	block.IsLoading = true
	block.Metadata["executing"] = "true"

	// Execute command
	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.CombinedOutput()

	block.IsLoading = false
	delete(block.Metadata, "executing")

	if err != nil {
		block.Error = err.Error()
		block.Type = BlockTypeError
		block.Output = string(output)
	} else {
		block.Output = string(output)
		block.Type = BlockTypeSuccess
	}

	vp := viewport.New(m.width-10, 10)
	vp.SetContent(block.Output)
	block.Viewport = vp
}

func (m *model) executeCommand(cmdStr string) {
	if m.selectedIdx >= len(m.blocks) {
		return
	}

	block := &m.blocks[m.selectedIdx]
	m.executeCommandInBlock(cmdStr, block)
}

func (m *model) addInfoBlock(message string) {
	infoBlock := Block{
		ID:        fmt.Sprintf("info-%d", time.Now().UnixNano()),
		Title:     "Info",
		Content:   message,
		Type:      BlockTypeInfo,
		Expanded:  true,
		Selected:  false,
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}

	vp := viewport.New(m.width-10, 5)
	vp.SetContent(message)
	infoBlock.Viewport = vp

	m.blocks = append(m.blocks, infoBlock)
	m.selectedIdx = len(m.blocks) - 1
	m.blocks[m.selectedIdx].Selected = true
}

func (m model) addHelpBlock() {
	helpContent := `Keyboard Shortcuts:
  j / ↓     - Navigate down
  k / ↑     - Navigate up
  e         - Expand/collapse block
  c         - Copy block content
  r         - Refresh/reload block
  d         - Delete block
  i         - Toggle input mode
  h         - Show this help
  q / Ctrl+C - Quit
  Space/Enter - Toggle block expansion

Block Types:
  • Command blocks (yellow) - Show command execution
  • Output blocks (green) - Show command output
  • Error blocks (red) - Show errors
  • Success blocks (green) - Show success messages
  • Info blocks (blue) - Show information
  • Progress blocks - Show progress indicators
  • Table blocks - Show tabular data`

	helpBlock := Block{
		ID:        "help",
		Title:     "Help & Keyboard Shortcuts",
		Content:   helpContent,
		Type:      BlockTypeInfo,
		Expanded:  true,
		Selected:  false,
		Timestamp: time.Now(),
		Metadata:  make(map[string]string),
	}

	vp := viewport.New(m.width-10, 20)
	vp.SetContent(helpContent)
	helpBlock.Viewport = vp

	m.blocks = append(m.blocks, helpBlock)
	m.selectedIdx = len(m.blocks) - 1
	m.blocks[m.selectedIdx].Selected = true
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)

	header := headerStyle.Render("╔═══ Gbloxs - Interactive Terminal Blocks ═══╗")
	b.WriteString(header + "\n\n")

	// Show help overlay if help mode is on
	if m.helpMode {
		helpContent := m.renderHelp()
		helpBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(1, 2).
			Width(m.width - 4).
			Background(lipgloss.Color("235")).
			Render(helpContent)
		b.WriteString(helpBox + "\n\n")
	}

	// Show table if table mode is on
	if m.showTable {
		tableBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("34")).
			Padding(1, 2).
			Render(m.table.View())
		b.WriteString(tableBox + "\n\n")
	}

	// Render blocks
	for i, block := range m.blocks {
		b.WriteString(m.renderBlock(block, i == m.selectedIdx))
		b.WriteString("\n")
	}

	// Input area
	if m.showInput {
		b.WriteString("\n")
		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("220")).
			Padding(1, 2).
			Render(
				m.styles.BlockTitle.Render("Input Mode (ESC to cancel, Enter to submit, /cmd or !cmd to execute):") + "\n" +
					m.textInput.View(),
			)
		b.WriteString(inputBox)
		b.WriteString("\n")
	}

	// Footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Align(lipgloss.Center).
		Width(m.width)

	shortcuts := "i: input | h: help | j/k: navigate | e: expand | c: copy | r: refresh | d: delete | x: execute | t: table | q: quit"
	footer := footerStyle.Render(shortcuts)
	b.WriteString("\n" + footer)

	return b.String()
}

func (m model) renderHelp() string {
	helpText := `╔═══════════════════════════════════════════════════════════════╗
║                    KEYBOARD SHORTCUTS                        ║
╠═══════════════════════════════════════════════════════════════╣
║  Navigation:                                                  ║
║    j / ↓     Navigate down to next block                      ║
║    k / ↑     Navigate up to previous block                    ║
║                                                               ║
║  Block Actions:                                               ║
║    e         Expand/collapse selected block                    ║
║    c         Copy block content to clipboard                   ║
║    r         Refresh/reload block content                      ║
║    d         Delete selected block                             ║
║    x         Execute command in selected block                  ║
║    Space     Toggle block expansion                            ║
║    Enter     Toggle block expansion                            ║
║                                                               ║
║  Modes:                                                       ║
║    i         Toggle input mode                                 ║
║    h         Toggle help (this screen)                         ║
║    t         Toggle table view                                 ║
║                                                               ║
║  Input Mode:                                                  ║
║    /cmd      Execute shell command (e.g., /ls -la)            ║
║    !cmd      Execute shell command (alternative)               ║
║    ESC       Cancel input                                      ║
║    Enter     Submit input                                     ║
║                                                               ║
║  General:                                                     ║
║    q         Quit application                                  ║
║    Ctrl+C    Quit application                                  ║
║    Ctrl+L    Clear all blocks                                 ║
╠═══════════════════════════════════════════════════════════════╣
║  BLOCK TYPES:                                                 ║
║    🟡 Command  - Yellow border, shows command execution        ║
║    🟢 Output   - Green border, shows command output            ║
║    🔴 Error    - Red border, shows error messages             ║
║    🔵 Info     - Blue border, shows information                ║
║    📊 Table    - Shows tabular data                            ║
║    ⏳ Progress - Shows progress indicators                     ║
╚═══════════════════════════════════════════════════════════════╝`

	return helpText
}

func (m model) renderBlock(block Block, selected bool) string {
	var style lipgloss.Style

	// Choose style based on block type and selection
	if selected {
		style = m.styles.SelectedBlock
	} else {
		switch block.Type {
		case BlockTypeCommand:
			style = m.styles.CommandBlock
		case BlockTypeOutput, BlockTypeSuccess:
			style = m.styles.OutputBlock
		case BlockTypeError:
			style = m.styles.ErrorBlock
		case BlockTypeInfo:
			style = m.styles.InfoBlock
		default:
			style = m.styles.BlockBorder
		}
	}

	// Build block content
	var content strings.Builder

	// Title with expand/collapse indicator
	expandIcon := "▼"
	if !block.Expanded {
		expandIcon = "▶"
	}
	title := fmt.Sprintf("%s %s", expandIcon, block.Title)
	if block.Selected {
		title = fmt.Sprintf("● %s", title)
	}
	content.WriteString(m.styles.BlockTitle.Render(title))
	content.WriteString("\n")

	if block.Expanded {
		// Timestamp
		timeStr := block.Timestamp.Format("15:04:05")
		content.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Render(fmt.Sprintf("  %s", timeStr)))
		content.WriteString("\n\n")

		// Render based on block type
		switch block.Type {
		case BlockTypeCommand:
			if block.Command != "" {
				content.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("220")).
					Render(fmt.Sprintf("  $ %s", block.Command)))
				content.WriteString("\n\n")
			}
			if block.Output != "" {
				content.WriteString(m.renderOutput(block.Output))
			}

		case BlockTypeProgress:
			content.WriteString(m.progress.ViewAs(block.Progress))
			if block.IsLoading {
				content.WriteString(" " + m.spinner.View())
			}
			content.WriteString("\n")
			content.WriteString(fmt.Sprintf("  %.0f%% complete", block.Progress*100))

		case BlockTypeTable:
			if len(block.TableData) > 0 {
				content.WriteString(m.renderTable(block.TableData))
			} else {
				content.WriteString(m.table.View())
			}

		case BlockTypeError:
			content.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Render("  ✗ " + block.Error))
			if block.Content != "" {
				content.WriteString("\n" + m.renderOutput(block.Content))
			}

		case BlockTypeSuccess:
			content.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("46")).
				Render("  ✓ " + block.Content))

		default:
			if block.Content != "" {
				content.WriteString(m.renderOutput(block.Content))
			} else if block.Output != "" {
				content.WriteString(m.renderOutput(block.Output))
			}
		}

		// Metadata
		if len(block.Metadata) > 0 {
			content.WriteString("\n")
			for k, v := range block.Metadata {
				content.WriteString(lipgloss.NewStyle().
					Foreground(lipgloss.Color("240")).
					Render(fmt.Sprintf("  [%s: %s]", k, v)))
			}
		}
	}

	return style.Render(content.String())
}

func (m model) renderOutput(output string) string {
	// Enhanced syntax highlighting
	lines := strings.Split(output, "\n")
	var highlighted strings.Builder

	// Patterns for syntax highlighting
	dirPattern := regexp.MustCompile(`^d[rwx-]{9}`)
	filePattern := regexp.MustCompile(`^-rw`)
	execPattern := regexp.MustCompile(`^-rwx`)
	errorPattern := regexp.MustCompile(`(?i)(error|failed|fatal|exception)`)
	successPattern := regexp.MustCompile(`(?i)(success|ok|done|complete)`)
	numberPattern := regexp.MustCompile(`\d+`)
	pathPattern := regexp.MustCompile(`(/[^\s]+|\./[^\s]+|~\w+)`)

	for _, line := range lines {
		if line == "" {
			highlighted.WriteString("\n")
			continue
		}

		style := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

		// Directory detection
		if dirPattern.MatchString(line) {
			style = style.Foreground(lipgloss.Color("39")) // Blue for directories
		} else if execPattern.MatchString(line) {
			style = style.Foreground(lipgloss.Color("46")) // Green for executables
		} else if filePattern.MatchString(line) {
			style = style.Foreground(lipgloss.Color("252")) // White for files
		}

		// Error highlighting
		if errorPattern.MatchString(line) {
			style = style.Foreground(lipgloss.Color("196")).Bold(true)
		}

		// Success highlighting
		if successPattern.MatchString(line) {
			style = style.Foreground(lipgloss.Color("46"))
		}

		// Highlight paths
		line = pathPattern.ReplaceAllStringFunc(line, func(match string) string {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("220")).
				Underline(true).
				Render(match)
		})

		// Highlight numbers
		line = numberPattern.ReplaceAllStringFunc(line, func(match string) string {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render(match)
		})

		highlighted.WriteString(style.Render("  " + line))
		highlighted.WriteString("\n")
	}

	return highlighted.String()
}

func (m model) renderTable(data [][]string) string {
	if len(data) == 0 {
		return ""
	}

	var b strings.Builder

	// Header
	header := data[0]
	headerRow := strings.Builder{}
	for i, cell := range header {
		if i > 0 {
			headerRow.WriteString(" │ ")
		}
		headerRow.WriteString(m.styles.TableHeader.Render(cell))
	}
	b.WriteString("  " + headerRow.String() + "\n")
	b.WriteString("  " + strings.Repeat("─", len(headerRow.String())) + "\n")

	// Rows
	for i := 1; i < len(data); i++ {
		row := strings.Builder{}
		for j, cell := range data[i] {
			if j > 0 {
				row.WriteString(" │ ")
			}
			row.WriteString(m.styles.TableCell.Render(cell))
		}
		b.WriteString("  " + row.String() + "\n")
	}

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
