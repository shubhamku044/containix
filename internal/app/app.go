package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shubhamku044/containix/internal/ui"
)

// Run initializes and starts the application
func Run() {
	p := tea.NewProgram(ui.NewMainModel())

	if err := p.Start(); err != nil {
		panic(err)
	}
}

