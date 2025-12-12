package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/vakidze/Linux-og-ferari/internal/logger"
	"github.com/vakidze/Linux-og-ferari/internal/serial"
)

func main() {
	a := app.New()
	w := a.NewWindow("Linux Serial Logger GUI")
	w.Resize(fyne.NewSize(600, 400))

	// Serial port input
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("/dev/ttyUSB0")

	baudEntry := widget.NewEntry()
	baudEntry.SetPlaceHolder("115200")

	logBox := widget.NewMultiLineEntry()
	logBox.Wrapping = fyne.TextWrapWord
	logBox.SetMinRowsVisible(15)
	logBox.Disable() // instead of SetReadOnly

	status := widget.NewLabel("Idle")

	startBtn := widget.NewButton("Start Logging", nil)
	stopBtn := widget.NewButton("Stop", nil)
	stopBtn.Disable()

	var stopCh chan struct{}
	var l *logger.CSVLogger

	startBtn.OnTapped = func() {
		if portEntry.Text == "" {
			status.SetText("Port required")
			return
		}

		baud := 115200
		if baudEntry.Text != "" {
			fmt.Sscan(baudEntry.Text, &baud)
		}

		p, err := serial.Open(portEntry.Text, baud)
		if err != nil {
			status.SetText("Open error: " + err.Error())
			return
		}

		l = logger.New("logs")
		stopCh = make(chan struct{})

		startBtn.Disable()
		stopBtn.Enable()
		status.SetText("Running...")

		go func() {
			buf := make([]byte, 256)

			for {
				select {
				case <-stopCh:
					p.Close()
					return
				default:
					n, err := p.Read(buf)
					if err != nil {
						time.Sleep(100 * time.Millisecond)
						continue
					}

					data := string(buf[:n])

					// update GUI thread-safe
					logBox.SetText(logBox.Text + data)

					l.WriteLine(data)
				}
			}
		}()
	}

	stopBtn.OnTapped = func() {
		if stopCh != nil {
			close(stopCh)
		}

		startBtn.Enable()
		stopBtn.Disable()
		status.SetText("Stopped")
		if l != nil {
			l.Close()
		}
	}

	form := container.NewVBox(
		widget.NewLabel("Serial Port:"),
		portEntry,
		widget.NewLabel("Baud:"),
		baudEntry,
		startBtn,
		stopBtn,
		status,
		logBox,
	)

	w.SetContent(form)
	w.ShowAndRun()
}
