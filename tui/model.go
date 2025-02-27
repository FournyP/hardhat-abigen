package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the interactive CLI state
type Model struct {
	Params map[string]string
	Keys   []string
	Index  int
}

// NewModel initializes a new Model instance
func NewModel(params map[string]string, keys []string) Model {
	return Model{Params: params, Keys: keys, Index: 0}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case string:
		if m.Index < len(m.Keys) {
			m.Params[m.Keys[m.Index]] = msg
			m.Index++
		}
		if m.Index >= len(m.Keys) {
			return m, tea.Quit
		}
	}

	return m, m.askNext()
}

func (m Model) View() string {
	if m.Index < len(m.Keys) {
		return fmt.Sprintf("Enter %s: ", m.Keys[m.Index])
	}
	return "All parameters set!\n"
}

func (m Model) askNext() tea.Cmd {
	if m.Index < len(m.Keys) {
		return func() tea.Msg {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			return strings.TrimSpace(text)
		}
	}
	return nil
}
