package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/shubhamku044/containix/internal/ui"
)

func Run() {
	p := tea.NewProgram(ui.NewMainModel())

	_, err := p.Run()
	if err != nil {
		panic(err)
	}
}
