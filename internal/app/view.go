package app

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/jroimartin/gocui"
)

type Model struct {
	Width      int
	Height     int
	Err        error
	Containers []types.Container
	Images     []types.ImageSummary
}

func (m *Model) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	m.Width, m.Height = maxX, maxY

	// Header
	if v, err := g.SetView("header", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, "Containix - Container Management")
	}

	// Containers Panel
	if v, err := g.SetView("containers", 0, 3, maxX/3-1, maxY-4); err != nil { // Adjusted height to avoid footer overlap
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Containers"
		fmt.Fprintln(v, renderContainers(m.Containers))
	}

	// Images Panel
	if v, err := g.SetView("images", maxX/3, 3, 2*maxX/3-1, maxY-4); err != nil { // Adjusted height to avoid footer overlap
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Images"
		fmt.Fprintln(v, renderImages(m.Images))
	}

	// Info Panel
	if v, err := g.SetView("info", 2*maxX/3, 3, maxX-1, maxY-4); err != nil { // Adjusted height to avoid footer overlap
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Info"
		fmt.Fprintln(v, renderInfoPanel(len(m.Containers), len(m.Images)))
	}

	// Footer
	if v, err := g.SetView("footer", 0, maxY-3, maxX-1, maxY-1); err != nil { // Corrected footer position
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Clear() // Clear the view to avoid overlapping content
		fmt.Fprintln(v, "Tab: Switch panels | R: Refresh | Q: Quit")
	}

	return nil
}

func renderContainers(containers []types.Container) string {
	var sb strings.Builder
	for _, c := range containers {
		name := strings.TrimPrefix(c.Names[0], "/")
		sb.WriteString(fmt.Sprintf("%s\n%s\n", name, c.Status))
	}
	return sb.String()
}

func renderImages(images []types.ImageSummary) string {
	var sb strings.Builder
	for _, img := range images {
		if len(img.RepoTags) > 0 {
			tag := strings.Split(img.RepoTags[0], ":")[0]
			sb.WriteString(fmt.Sprintf("%s\n%s\n", tag, img.ID[7:19]))
		}
	}
	return sb.String()
}

func renderInfoPanel(containerCount, imageCount int) string {
	return fmt.Sprintf("Containers: %d\nImages: %d\nVersion: 23.0", containerCount, imageCount)
}

func Keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	views := []string{"containers", "images", "info"}
	var next string

	for i, name := range views {
		if v != nil && v.Name() == name {
			next = views[(i+1)%len(views)]
			break
		}
	}

	if next == "" {
		next = views[0]
	}

	_, err := g.SetCurrentView(next)
	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
