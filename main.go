// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/jroimartin/gocui"
	"github.com/mbndr/figlet4go"
	"github.com/phux/pomogoro/config"
	"github.com/phux/pomogoro/history"
	"github.com/phux/pomogoro/timing"
	"github.com/phux/pomogoro/todotxt"
)

var conf *config.Config
var confPath *string
var displayingSummary bool
var timer *timing.Timer

func main() {
	confPath = flag.String("config", "", "absolute path to toml config file")
	flag.Parse()
	conf = config.Load(*confPath)

	runCui()
}

func summary(g *gocui.Gui, v *gocui.View) error {
	if displayingSummary {
		_, err := g.SetCurrentView("todos")
		if err != nil {
			return err
		}

		displayingSummary = false
		v.Highlight = false
		todosView, _ := g.View("todos")
		todosView.Highlight = true
		updateRecentLog(g)
		return nil
	}

	g.Update(func(g *gocui.Gui) error {
		todosView, _ := g.View("todos")
		todosView.Highlight = false
		v, err := g.View("recentLog")
		if err != nil {
			log.Printf("could not update recent log view: %s", err)
		}
		v.Clear()

		history, err := history.Aggregate(conf)
		if err != nil {
			fmt.Fprintln(v, err)
			return nil
		}

		dailyStats := history.GetStatsForLastDays(30)
		for _, stat := range dailyStats {
			for _, line := range stat.Summary() {
				fmt.Fprintln(v, fmt.Sprintf("%s :::: %s", stat.GetDate(), line))
			}
			fmt.Fprintln(v, "================")
		}

		displayingSummary = true
		v.Highlight = true
		_, err = g.SetCurrentView("recentLog")
		return err
	})
	return nil
}

func trackTime(g *gocui.Gui, v *gocui.View) error {
	if timer != nil {
		if timer.IsActive {
			updateTimerView("Multitask", g)
			return nil
		}
		if conf.LogIdleTime {
			history.Append("Idle", int(timer.IdlingTime().Seconds()), conf)
		}
	}

	_, y := v.Cursor()

	var err error
	remaining := math.MaxInt32
	if conf.PomodoroEnabled {
		remaining = conf.PomodoroDuration * 60
	}
	timer = timing.NewTimer(remaining)
	timer.CurrentTodo, err = v.Line(y)
	if err != nil {
		log.Printf("Could not fetch current todo from view: %s", err)
	}

	go startTracking(g)
	return nil
}

func notify(action, message string) {
	body := fmt.Sprintf("%s\n%s", action, message)
	err := beeep.Notify("Pomogoro", body, "")
	if err != nil {
		log.Fatal(err)
	}
}

func startTracking(g *gocui.Gui) {
	for timer.Progress() {
		if conf.PomodoroEnabled {
			updateTimerView(timer.RemainingToString(), g)
		} else {
			updateTimerView(timer.ElapsedToString(), g)
		}
		time.Sleep(time.Second)
	}
	history.Append(timer.CurrentTodo, timer.Elapsed, conf)
	updateRecentLog(g)
	notify("Finished", timer.CurrentTodo)
	timer.Idle()

	if conf.PomodoroEnabled {
		go runBreak(g)
	}
}

func runBreak(g *gocui.Gui) {
	timer = timing.NewTimer(conf.BreakDuration * 60)
	timer.CurrentTodo = "Break"
	notify("Started", timer.CurrentTodo)
	for timer.Progress() {
		updateTimerView(timer.RemainingToString(), g)
		time.Sleep(time.Second)
	}
	updateTimerView("IDLE", g)
	notify("Finished", timer.CurrentTodo)
	timer.CurrentTodo = "--------------"
	if conf.LogBreakTime {
		history.Append("Break", timer.Elapsed, conf)
	}
	timer.Reset()
	timer.Idle()
}

func updateTimerView(text string, g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("timer")
		if err != nil {
			log.Printf("could not update timer view: %s", err)
		}
		v.Clear()

		ascii := figlet4go.NewAsciiRender()
		str, err := ascii.Render(text)
		if err != nil {
			log.Printf("could not render text '%s', err: %s", text, err)
			str = text
		}
		fmt.Fprintln(v, str)
		if timer != nil {
			fmt.Fprintln(v, timer.CurrentTodo)
		}

		return nil
	})
}

func updateTodoView(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	list, err := todotxt.ReadTodoTxt(conf)
	if err != nil {
		return err
	}
	for _, todo := range list {
		fmt.Fprintln(v, todo)
	}
	return nil
}

func updateRecentLog(g *gocui.Gui) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("recentLog")
		if err != nil {
			log.Printf("could not update recent log view: %s", err)
		}
		v.Clear()

		history, err := history.Aggregate(conf)
		if err != nil {
			fmt.Fprintln(v, err)
			return nil
		}

		dailyStats := history.GetStatsForLastDays(14)
		for _, stat := range dailyStats {
			for _, line := range stat.ToList() {
				fmt.Fprintln(v, line)
			}
		}

		return nil
	})
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func pauseToggle(g *gocui.Gui, v *gocui.View) error {
	if timer != nil {
		timer.PauseToggle()
		updateTimerView("Paused", g)
	}
	return nil
}

func cancel(g *gocui.Gui, v *gocui.View) error {
	if timer != nil {
		timer.Cancel()
		updateTimerView("Canceled", g)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	if timer != nil {
		if !timer.IsActive {
			if conf.LogIdleTime {
				history.Append("Idle", int(timer.IdlingTime().Seconds()), conf)
			}
		} else {
			err := cancel(g, v)
			if err != nil {
				return err
			}

			time.Sleep(time.Second)
		}
	}
	return gocui.ErrQuit
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, 1, false)
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	v.MoveCursor(0, -1, false)
	return nil
}

func runCui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone, cursorDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'k', gocui.ModNone, cursorUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'c', gocui.ModNone, cancel); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 's', gocui.ModNone, summary); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("todos", gocui.KeyEnter, gocui.ModNone, trackTime); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("todos", 'p', gocui.ModNone, pauseToggle); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("todos", 'r', gocui.ModNone, updateTodoView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", '?', gocui.ModNone, toggleHelpWindow); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

var displayingHelp bool

func toggleHelpWindow(g *gocui.Gui, v *gocui.View) error {
	v, _ = g.View("help")
	if !displayingHelp {
		_, _ = g.SetViewOnTop("help")
		displayingHelp = true
	} else {
		_, _ = g.SetViewOnBottom("help")
		displayingHelp = false
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("timer", maxX/2, 0, maxX, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Title = "Timer"
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		fmt.Fprintln(v, "WAITING (press enter on a todo to start)")
	}
	if v, err := g.SetView("recentLog", maxX/2, maxY/2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "History"
		v.Highlight = false
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Wrap = true
		updateRecentLog(g)
	}
	if v, err := g.SetView("help", 0, maxY-9, maxX/2, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Wrap = true
		v.Title = "Shortcuts"
		text := `
j/k/Up/Down: move cursor up/down
Enter: Start tracking
p: pause timer
c: cancel timer
s: toggle summary/recent history
r: refresh todo list
CTRL-c: quit
?: show this help
		`
		fmt.Fprintln(v, text)
	}
	if v, err := g.SetView("todos", 0, 0, maxX/2, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Wrap = true
		v.Title = "Todos"
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		err := updateTodoView(g, v)
		if err != nil {
			return err
		}

		if _, err = setCurrentViewOnTop(g, "todos"); err != nil {
			return err
		}
	}

	return nil
}
