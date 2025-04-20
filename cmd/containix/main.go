package main

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/shubhamku044/containix/internal/app"
	"github.com/shubhamku044/containix/internal/docker"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// Fetch real data
	containers := docker.FetchContainers()
	images := docker.FetchImages()

	// Initialize the model with real data
	model := &app.Model{
		Containers: containers, // Use exported field
		Images:     images,     // Use exported field
	}

	g.SetManagerFunc(func(g *gocui.Gui) error {
		return model.Layout(g)
	})

	if err := app.Keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
