package ui

import (
	"os/exec"
	"strings"

	"pkgman/termux"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type commandFinishedMsg struct {
	output string
	err    error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.pkgList.SetSize(msg.Width, msg.Height-5) // leave room for tabs

	case packagesLoadedMsg:
		m.isLoading = false
		m.allPkgs = msg.packages
		m.updateListItems()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == stateMenu {
				return m, tea.Quit
			}
			m.state = stateMenu
			return m, nil
		case "esc":
			m.state = stateMenu
			return m, nil
		}

		if m.state == stateMenu {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.menuOptions)-1 {
					m.cursor++
				}
			case "enter":
				m.selected = actionType(m.cursor + 1)
				return m.handleMenuSelection()
			}
		} else if m.state == stateInput {
			switch msg.String() {
			case "enter":
				val := strings.TrimSpace(m.textInput.Value())
				if val != "" {
					return m.executeAction(val)
				}
			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		} else if m.state == statePackageBrowser {
			if m.pkgList.FilterState() == list.Filtering {
				m.pkgList, cmd = m.pkgList.Update(msg)
				return m, cmd
			}

			switch msg.String() {
			case "tab", "right", "l":
				m.activeTab = (m.activeTab + 1) % len(m.tabs)
				m.updateListItems()
				return m, nil
			case "shift+tab", "left", "h":
				m.activeTab--
				if m.activeTab < 0 {
					m.activeTab = len(m.tabs) - 1
				}
				m.updateListItems()
				return m, nil
			case " ":
				selectedItem := m.pkgList.SelectedItem()
				if selectedItem != nil {
					pkg := selectedItem.(packageItem).pkg
					if m.selectedPkgs[pkg.Name] {
						delete(m.selectedPkgs, pkg.Name)
					} else {
						m.selectedPkgs[pkg.Name] = true
					}
					m.updateListItems()
				}
				return m, nil
			case "enter":
				if len(m.selectedPkgs) > 0 {
					m.actionMenu = []string{"Bulk Install", "Bulk Update", "Bulk Reinstall", "Bulk Remove"}
				} else {
					selectedItem := m.pkgList.SelectedItem()
					if selectedItem != nil {
						pkg := selectedItem.(packageItem).pkg
						if pkg.Installed {
							m.actionMenu = []string{"Update", "Reinstall", "Remove"}
						} else {
							m.actionMenu = []string{"Install"}
						}
					} else {
						return m, nil
					}
				}
				m.actionCursor = 0
				m.state = statePackageAction
				return m, nil
			}
			m.pkgList, cmd = m.pkgList.Update(msg)
			return m, cmd
		} else if m.state == statePackageAction {
			switch msg.String() {
			case "up", "k":
				if m.actionCursor > 0 {
					m.actionCursor--
				}
			case "down", "j":
				if m.actionCursor < len(m.actionMenu)-1 {
					m.actionCursor++
				}
			case "esc", "q":
				m.state = statePackageBrowser
			case "enter":
				action := m.actionMenu[m.actionCursor]
				var targets []string
				if len(m.selectedPkgs) > 0 {
					for k := range m.selectedPkgs {
						targets = append(targets, k)
					}
				} else {
					selectedItem := m.pkgList.SelectedItem()
					if selectedItem != nil {
						targets = append(targets, selectedItem.(packageItem).pkg.Name)
					}
				}

				if len(targets) > 0 {
					switch action {
					case "Install", "Bulk Install", "Update", "Bulk Update":
						m.selected = actionInstall
					case "Reinstall", "Bulk Reinstall":
						m.selected = actionReinstall
					case "Remove", "Bulk Remove":
						m.selected = actionRemove
					}
					m.selectedPkgs = make(map[string]bool)
					return m.executeMultiAction(targets)
				}
			}
			return m, nil
		}

	case commandFinishedMsg:
		m.state = stateOutput
		m.output = msg.output
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
		}
	}

	return m, cmd
}

func (m *Model) handleMenuSelection() (tea.Model, tea.Cmd) {
	// Offset action by 1 because index 0 is "Browse All Packages"
	if m.cursor == 0 {
		m.state = statePackageBrowser
		return *m, nil
	}

	m.selected = actionType(m.cursor) // Adjust for the new menu item
	
	switch m.selected {
	case actionUpdate:
		return m.executeAction("")
	case actionClean:
		return m.executeAction("")
	case actionAutoRemove:
		return m.executeAction("")
	case actionRepo:
		return m.executeAction("")
	case actionSearch, actionInstall, actionRemove:
		m.state = stateInput
		m.textInput.SetValue("")
		switch m.selected {
		case actionSearch:
			m.textInput.Placeholder = "Enter package name to search..."
		case actionInstall:
			m.textInput.Placeholder = "Enter package name to install..."
		case actionRemove:
			m.textInput.Placeholder = "Enter package name to remove..."
		}
		return *m, textinput.Blink
	}
	return *m, nil
}

func (m *Model) executeAction(arg string) (tea.Model, tea.Cmd) {
	var c *exec.Cmd
	switch m.selected {
	case actionUpdate:
		c = termux.GetCommand("update")
	case actionClean:
		c = termux.GetCommand("clean")
	case actionAutoRemove:
		c = termux.GetCommand("autoremove")
	case actionRepo:
		c = termux.GetCommand("repo")
	case actionSearch:
		// Search can be captured easily
		return *m, func() tea.Msg {
			out, err := termux.SearchPackages(arg)
			return commandFinishedMsg{output: string(out), err: err}
		}
	case actionInstall:
		c = termux.GetCommand("install", arg)
	case actionRemove:
		c = termux.GetCommand("remove", arg)
	}

	if c != nil {
		// For commands that might need terminal (like repo, install), use Exec
		return *m, tea.ExecProcess(c, func(err error) tea.Msg {
			if err != nil {
				return commandFinishedMsg{output: "Command finished with error", err: err}
			}
			return commandFinishedMsg{output: "Command completed successfully. Press Enter to continue.", err: nil}
		})
	}

	return *m, nil
}

func (m *Model) updateListItems() {
	var items []list.Item
	tab := m.tabs[m.activeTab]

	for _, p := range m.allPkgs {
		include := false
		switch tab {
		case "All":
			include = true
		case "Installed":
			include = p.Installed
		case "Stable":
			include = strings.Contains(strings.ToLower(p.Category), "stable")
		case "X11":
			include = strings.Contains(strings.ToLower(p.Category), "x11")
		case "Root":
			include = strings.Contains(strings.ToLower(p.Category), "root")
		default:
			include = true
		}

		if include {
			items = append(items, packageItem{pkg: p, selectedPkgs: m.selectedPkgs})
		}
	}

	m.pkgList.SetItems(items)
}

func (m *Model) executeMultiAction(targets []string) (tea.Model, tea.Cmd) {
	var c *exec.Cmd
	switch m.selected {
	case actionInstall:
		c = termux.GetCommand("install", targets...)
	case actionRemove:
		c = termux.GetCommand("remove", targets...)
	case actionReinstall:
		c = termux.GetCommand("reinstall", targets...)
	}

	if c != nil {
		return *m, tea.ExecProcess(c, func(err error) tea.Msg {
			if err != nil {
				return commandFinishedMsg{output: "Command finished with error", err: err}
			}
			return commandFinishedMsg{output: "Command completed successfully. Press Enter to continue.", err: nil}
		})
	}
	return *m, nil
}
