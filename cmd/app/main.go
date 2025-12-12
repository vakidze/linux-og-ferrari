package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/vakidze/Linux-og-ferari/internal/serial"
	"github.com/vakidze/Linux-og-ferari/internal/logger"
)

func main() {
	a := app.New()
	w := a.NewWindow("Linux Serial Logger - GUI")
	w.Resize(fyne.NewSize(900, 600))

	ports, _ := serial.ListPorts()
	portSelect := widget.NewSelect(ports, func(string) {})

	startBtn := widget.NewButton("Start", nil)
	stopBtn := widget.NewButton("Stop", nil)
	stopBtn.Disable()

	filterCheck := widget.NewCheck("Filter login prompts (login:~$)", func(bool) {})

	saveBtn := widget.NewButton("Save CSV", nil)
	saveBtn.Disable()

	logBox := widget.NewMultiLineEntry()
	logBox.SetReadOnly(true)

	status := widget.NewLabel("Ready")

	// logger
	csvLogger, err := logger.NewCSVLogger()
	if err != nil {
		log.Printf("failed to create logger: %v", err)
	} else {
		saveBtn.Enable()
	}

	// channels and control
	lineCh := make(chan string, 200)
	stopCh := make(chan struct{})
	running := false

	appendLine := func(s string) {
		// prepend newest on top
		if logBox.Text == "" {
			logBox.SetText(s)
		} else {
			logBox.SetText(s + "\n" + logBox.Text)
		}
	}

	// consumer: write to UI and CSV
	go func() {
		for l := range lineCh {
			ts := time.Now().Format("2006-01-02 15:04:05.000")
			appendLine(fmt.Sprintf("%s  %s", ts, l))
			if csvLogger != nil {
				csvLogger.Write(ts, l)
			}
		}
	}()

	startBtn.OnTapped = func() {
		if running {
			return
		}
		port := portSelect.Selected
		if port == "" {
			status.SetText("Select a port first")
			return
		}

		// start reader
		running = true
		startBtn.Disable()
		stopBtn.Enable()
		status.SetText("Logging...")

		go func() {
			filterFn := func(s string) bool {
				if filterCheck.Checked && s == "login:~$" {
					return false
				}
				return true
			}
			ch, err := serial.ReadLinesToChan(port, 115200, filterFn)
			if err != nil {
				status.SetText("Failed to open port: " + err.Error())
				running = false
				startBtn.Enable()
				stopBtn.Disable()
				return
			}
			for l := range ch {
				select {
				default:
					lineCh <- l
				}
			}
		}()
	}

	stopBtn.OnTapped = func() {
		if !running {
			return
		}
		// stopping by closing port handled inside serial package; here we just toggle UI
		running = false
		startBtn.Enable()
		stopBtn.Disable()
		status.SetText("Stopped")
	}

	saveBtn.OnTapped = func() {
		if csvLogger == nil {
			status.SetText("No logger")
			return
		}
		csvLogger.CloseAndSave()
		// create new logger file after saving
		newLogger, err := logger.NewCSVLogger()
		if err == nil {
			csvLogger = newLogger
			status.SetText("Saved and started new CSV")
		} else {
			status.SetText("Saved. Failed to start new CSV")
		}
	}

	// layout
	buttons := container.NewHBox(startBtn, stopBtn, saveBtn, filterCheck)
	top := container.NewVBox(portSelect, buttons, status)
	content := container.NewBorder(top, nil, nil, nil, container.NewVScroll(logBox))

	w.SetContent(content)
	w.ShowAndRun()

	// on exit, ensure logger closed
	if csvLogger != nil {
		csvLogger.CloseAndSave()
	}
}
