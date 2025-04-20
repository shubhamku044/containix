package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	headerHeight = 5
	footerHeight = 3
)

func Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Header
	if v, err := g.SetView("header", 0, 0, maxX-1, headerHeight-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, "╔════════════════════════════════════════════════════════════════════╗")
		fmt.Fprintln(v, "║                     Containix - Container Management               ║")
		fmt.Fprintln(v, "╚════════════════════════════════════════════════════════════════════╝")
	}

	// Containers Panel
	if v, err := g.SetView("containers", 0, headerHeight, maxX/3-1, maxY-footerHeight-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Containers"
		fmt.Fprintln(v, "Container 1\nContainer 2\nContainer 3")
	}

	// Images Panel
	if v, err := g.SetView("images", maxX/3, headerHeight, 2*maxX/3-1, maxY-footerHeight-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Images"
		fmt.Fprintln(v, "Image 1\nImage 2\nImage 3")
	}

	// Info Panel
	if v, err := g.SetView("info", 2*maxX/3, headerHeight, maxX-1, maxY-footerHeight-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Info"
		fmt.Fprintln(v, "Containers: 3\nImages: 3\nVersion: 23.0")
	}

	// Footer
	if v, err := g.SetView("footer", 0, maxY-footerHeight, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprintln(v, "Tab: Switch panels | r: refresh | q: quit")
	}

	return nil
}
