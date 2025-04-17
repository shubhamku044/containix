package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shubhamku044/containix/ui"
)

func Execute() {
	p := tea.NewProgram(ui.NewModel())

	if err := p.Start(); err != nil {
		panic(err)
	}
}
