package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Choice represents a single option in a choice selector.
type Choice struct {
	Label string
	Value any
}

// ChoiceDialog is a Bubble Tea model for selecting from multiple choices.
type ChoiceDialog struct {
	Prompt      string
	Choices     []Choice
	Cursor      int
	Selected    bool
	Cancelled   bool
	StyleConfig StyleConfig
}

var _ tea.Model = (*ChoiceDialog)(nil)

// StyleConfig holds styling configuration for the choice component.
type StyleConfig struct {
	HighlightStyle lipgloss.Style
	NormalStyle    lipgloss.Style
	PromptStyle    lipgloss.Style
	HelpStyle      lipgloss.Style
}

// DefaultStyleConfig returns sensible default styles.
func DefaultStyleConfig() StyleConfig {
	return StyleConfig{
		HighlightStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")),
		NormalStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),
		PromptStyle: lipgloss.NewStyle(),
		HelpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Faint(true),
	}
}

// NewChoiceDialog creates a new choice selector model.
func NewChoiceDialog(prompt string, choices []Choice) ChoiceDialog {
	return ChoiceDialog{
		Prompt:      prompt,
		Choices:     choices,
		Cursor:      0,
		StyleConfig: DefaultStyleConfig(),
	}
}

func (m ChoiceDialog) Init() tea.Cmd {
	return nil
}

func (m ChoiceDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.Cancelled = true
			return m, tea.Quit

		case "enter":
			m.Selected = true
			return m, tea.Quit

		case "left", "h", "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Choices) - 1
			}

		case "right", "l", "down", "j":
			m.Cursor++
			if m.Cursor >= len(m.Choices) {
				m.Cursor = 0
			}

		case "home":
			m.Cursor = 0

		case "end":
			m.Cursor = len(m.Choices) - 1

		default:
			// Check if the key matches any choice label's first character
			key := strings.ToLower(msg.String())
			if len(key) == 1 {
				for i, choice := range m.Choices {
					if len(choice.Label) > 0 && strings.ToLower(string(choice.Label[0])) == key {
						m.Cursor = i
						m.Selected = true
						return m, tea.Quit
					}
				}
			}
		}
	}

	return m, nil
}

func (m ChoiceDialog) View() string {
	if m.Selected || m.Cancelled {
		return ""
	}

	var b strings.Builder

	if m.Prompt != "" {
		b.WriteString(m.StyleConfig.PromptStyle.Render(m.Prompt))
		b.WriteString(" ")
	}

	for i, choice := range m.Choices {
		if i > 0 {
			b.WriteString(" / ")
		}

		if i == m.Cursor {
			b.WriteString(m.StyleConfig.HighlightStyle.Render(choice.Label))
		} else {
			b.WriteString(m.StyleConfig.NormalStyle.Render(choice.Label))
		}
	}

	b.WriteString("\n")
	b.WriteString(m.StyleConfig.HelpStyle.Render("(Use arrow keys to select, Enter to confirm, Esc to cancel)"))
	b.WriteString("\n")

	return b.String()
}

// GetSelectedChoice returns the currently selected choice.
func (m ChoiceDialog) GetSelectedChoice() *Choice {
	if m.Cancelled || !m.Selected || m.Cursor >= len(m.Choices) {
		return nil
	}
	return &m.Choices[m.Cursor]
}

// GetSelectedValue returns the value of the selected choice.
func (m ChoiceDialog) GetSelectedValue() any {
	choice := m.GetSelectedChoice()
	if choice == nil {
		return nil
	}
	return choice.Value
}

// Run is a convenience method to run the choice selector and return the result.
func (m ChoiceDialog) Run() (*Choice, error) {
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(ChoiceDialog)
	return result.GetSelectedChoice(), nil
}
