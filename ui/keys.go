package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Start      key.Binding
	Restart    key.Binding
	Remove     key.Binding
	Inspect    key.Binding
	RawInspect key.Binding
	Logs       key.Binding
	Exec       key.Binding
	ExecCmd    key.Binding
	Refresh    key.Binding
	Search     key.Binding
	Filter     key.Binding
	Quit       key.Binding
	Help       key.Binding
}

var keys = keyMap{
	Up:         key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/↑", "up")),
	Down:       key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/↓", "down")),
	Start:      key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "start/stop")),
	Restart:    key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restart")),
	Remove:     key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "remove")),
	Inspect:    key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "inspect")),
	RawInspect: key.NewBinding(key.WithKeys("I"), key.WithHelp("I", "raw inspect")),
	Logs:       key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "logs")),
	Exec:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "exec")),
	ExecCmd:    key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "exec cmd")),
	Refresh:    key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "refresh")),
	Search:     key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	Filter:     key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "filter")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
}
