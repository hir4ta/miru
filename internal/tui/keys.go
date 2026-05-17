package tui

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	HalfUp      key.Binding
	HalfDown    key.Binding
	Top         key.Binding
	Bottom      key.Binding
	PrevSection key.Binding
	NextSection key.Binding
	Browser     key.Binding
	Settings    key.Binding
	Help        key.Binding
	Quit        key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		HalfUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("^u", "half page up"),
		),
		HalfDown: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("^d", "half page down"),
		),
		Top: key.NewBinding(
			key.WithKeys("g", "home"),
			key.WithHelp("g", "top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G", "end"),
			key.WithHelp("G", "bottom"),
		),
		PrevSection: key.NewBinding(
			key.WithKeys("{"),
			key.WithHelp("{", "prev section"),
		),
		NextSection: key.NewBinding(
			key.WithKeys("}"),
			key.WithHelp("}", "next section"),
		),
		Browser: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "open in browser"),
		),
		Settings: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "settings"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "esc"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.PrevSection, k.NextSection, k.Browser, k.Settings, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.HalfUp, k.HalfDown},
		{k.Top, k.Bottom},
		{k.PrevSection, k.NextSection},
		{k.Browser, k.Settings},
		{k.Help, k.Quit},
	}
}
