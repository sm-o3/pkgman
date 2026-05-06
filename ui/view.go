package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(lipgloss.Color("#EE6FF8")).
				Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E88388")).
			Bold(true).
			MarginBottom(1)
)

func (m Model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Termux Package Manager"))
	s.WriteString("\n")

	switch m.state {
	case stateMenu:
		for i, choice := range m.menuOptions {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
				s.WriteString(selectedItemStyle.Render(fmt.Sprintf("%s %s", cursor, choice)))
			} else {
				s.WriteString(itemStyle.Render(fmt.Sprintf("%s %s", cursor, choice)))
			}
			s.WriteString("\n")
		}
		s.WriteString("\nPress q to quit.\n")

	case statePackageBrowser:
		if m.isLoading {
			s.WriteString("Loading packages...\n")
		} else {
			var tabs []string
			for i, t := range m.tabs {
				if i == m.activeTab {
					tabs = append(tabs, lipgloss.NewStyle().Foreground(lipgloss.Color("#EE6FF8")).Bold(true).Render(t))
				} else {
					tabs = append(tabs, t)
				}
			}
			s.WriteString(strings.Join(tabs, " | "))
			s.WriteString("\n\n")
			s.WriteString(m.pkgList.View())
		}

	case statePackageAction:
		s.WriteString("Select an Action:\n\n")
		for i, action := range m.actionMenu {
			cursor := " "
			if m.actionCursor == i {
				cursor = ">"
				s.WriteString(selectedItemStyle.Render(fmt.Sprintf("%s %s", cursor, action)))
			} else {
				s.WriteString(itemStyle.Render(fmt.Sprintf("%s %s", cursor, action)))
			}
			s.WriteString("\n")
		}
		s.WriteString("\nPress Enter to execute, or Esc to cancel.\n")

	case stateInput:
		s.WriteString(m.textInput.View())
		s.WriteString("\n\nPress Esc to return to menu.\n")

	case stateOutput:
		if m.err != nil {
			s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
			s.WriteString("\n")
		}
		
		// Optional: Truncate output if it's too long, or use a viewport.
		// For simplicity, we just print it.
		s.WriteString(m.output)
		s.WriteString("\n\nPress Enter or Esc to return to menu.\n")
	}

	return s.String()
}
