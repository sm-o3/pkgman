package ui

import (
	"fmt"
	"pkgman/termux"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type packageItem struct {
	pkg          termux.Package
	selectedPkgs map[string]bool
}

func (i packageItem) Title() string {
	prefix := "[ ] "
	if i.selectedPkgs[i.pkg.Name] {
		prefix = "[x] "
	}
	return prefix + i.pkg.Name
}
func (i packageItem) Description() string {
	status := ""
	if i.pkg.Installed {
		status = " [Installed]"
	}
	return fmt.Sprintf("v%s | %s%s", i.pkg.Version, i.pkg.Category, status)
}
func (i packageItem) FilterValue() string { return i.pkg.Name }

type state int

const (
	stateMenu state = iota
	stateInput
	stateOutput
	statePackageBrowser
	statePackageAction
)

type actionType int

const (
	actionNone actionType = iota
	actionUpdate
	actionSearch
	actionInstall
	actionRemove
	actionReinstall
	actionClean
	actionAutoRemove
	actionRepo
)

type Model struct {
	state       state
	menuOptions []string
	cursor      int
	selected    actionType

	textInput textinput.Model
	output    string
	err       error

	pkgList   list.Model
	allPkgs   []termux.Package
	tabs      []string
	activeTab int
	isLoading bool

	selectedPkgs map[string]bool
	actionMenu   []string
	actionCursor int

	width  int
	height int
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Enter package name..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	delegate := list.NewDefaultDelegate()
	pkgList := list.New([]list.Item{}, delegate, 0, 0)
	pkgList.Title = "Browse Packages"

	return Model{
		state: stateMenu,
		menuOptions: []string{
			"Browse All Packages (Categories & Search)",
			"Update System (pkg upgrade)",
			"Search Package",
			"Install Package",
			"Remove Package",
			"Clean Package Cache (pkg clean)",
			"Auto-remove Unused Packages (apt autoremove)",
			"Repository Management (termux-change-repo)",
		},
		textInput:    ti,
		pkgList:      pkgList,
		tabs:         []string{"All", "Installed", "Stable", "X11", "Root"},
		selectedPkgs: make(map[string]bool),
	}
}

type packagesLoadedMsg struct {
	packages []termux.Package
	err      error
}

func loadPackages() tea.Msg {
	pkgs, err := termux.ListAllPackages()
	return packagesLoadedMsg{packages: pkgs, err: err}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, loadPackages)
}
