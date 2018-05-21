package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"time"
)

func main() {
	gtk.Init(nil)

	usage := `Timeline - a linear countdown timer

Usage:
  timeline <time> [s|m|h] [inc|dec] [options]
  timeline (-h | --help)

Options:
  -h --help      Show this text
  --normal       Normal window
  --fullscreen   Fullscreen mode
  --dock-top     Docked to the top edge
  --dock-bottom  Docked to the bottom edge (default)
  --dock-left    Docked to the left edge
  --dock-right   Docked to the right edge
  --width=<px>   Width in pixels of window in normal mode, 0=100% [default: 0]
  --height=<px>  Height in pixels of window in normal mode, 0=100% [default: 7]
  --update=<ms>  Update delay in milliseconds [default: 1000]`

	args, _ := docopt.ParseDoc(usage)

	argInc, err := args.Bool("inc")
	handleArgErr(err)

	argTime, err := args.Int("<time>")
	handleArgErr(err)

	argSeconds, err := args.Bool("s")
	handleArgErr(err)

	argMinutes, err := args.Bool("m")
	handleArgErr(err)

	argHours, err := args.Bool("h")
	handleArgErr(err)

	argNormal, err := args.Bool("--normal")
	handleArgErr(err)

	argDockTop, err := args.Bool("--dock-top")
	handleArgErr(err)

	argDockBottom, err := args.Bool("--dock-bottom")
	handleArgErr(err)

	argDockLeft, err := args.Bool("--dock-left")
	handleArgErr(err)

	argDockRight, err := args.Bool("--dock-right")
	handleArgErr(err)

	argFullscreen, err := args.Bool("--fullscreen")
	handleArgErr(err)

	argUpdate, err := args.Int("--update")
	handleArgErr(err)

	argWidth, err := args.Int("--width")
	handleArgErr(err)

	argHeight, err := args.Int("--height")
	handleArgErr(err)

	if !argSeconds && !argMinutes && !argHours {
		argMinutes = true
	}

	if argMinutes {
		argTime *= 60
	}

	if argHours {
		argTime *= 60 * 60
	}

	if argWidth == 0 {
		argWidth = gdk.ScreenWidth()
	}

	if argHeight == 0 {
		argHeight = gdk.ScreenHeight()
	}

	if !(argNormal || argFullscreen ||
		argDockTop || argDockBottom || argDockLeft || argDockRight) {
		argDockBottom = true
	}

	step := float64(argUpdate) / 1000.0 / float64(argTime)

	startFraction := 0.0
	if !argInc {
		startFraction = 1.0
		step *= -1
	}

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Timeline")
	window.Connect("destroy", gtk.MainQuit)

	bar := gtk.NewProgressBar()
	bar.SetFraction(startFraction)
	window.Add(bar)

	window.SetSizeRequest(argWidth, argHeight)
	winWidth, winHeight := window.GetSize()

	if argDockTop || argDockBottom {
		window.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
		bar.SetOrientation(gtk.PROGRESS_LEFT_TO_RIGHT)

		if argDockTop {
			window.SetGravity(gdk.GRAVITY_NORTH)
			window.Move(0, 0)
		} else {
			window.SetGravity(gdk.GRAVITY_SOUTH)
			window.Move(0, gdk.ScreenHeight()-winHeight)
		}
	} else if argDockLeft || argDockRight {
		window.SetTypeHint(gdk.WINDOW_TYPE_HINT_DOCK)
		bar.SetOrientation(gtk.PROGRESS_BOTTOM_TO_TOP)

		if argDockLeft {
			window.SetGravity(gdk.GRAVITY_WEST)
			window.Move(0, 0)
		} else {
			window.SetGravity(gdk.GRAVITY_EAST)
			window.Move(gdk.ScreenWidth()-winWidth, 0)
		}
	} else if argFullscreen {
		window.Fullscreen()
	}

	precision := 0.0000001
	ticker := time.NewTicker(time.Duration(argUpdate) * time.Millisecond)
	go func() {
		for range ticker.C {
			nextFraction := bar.GetFraction() + step
			if nextFraction >= 1.0-precision {
				nextFraction = 1.0
				os.Exit(0)
			} else if nextFraction <= 0.0+precision {
				nextFraction = 0.0
				os.Exit(0)
			}
			bar.SetFraction(nextFraction)
		}
	}()

	window.ShowAll()
	gtk.Main()
}

func handleArgErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
