//go:build !darwin

package main

import (
	_ "embed"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/windows/icon.ico
var trayIcon []byte

func startSystray() {
	go systray.Run(onSystrayReady, nil)
}

func stopSystray() {
	systray.Quit()
}

func hideWindowOnClose() bool {
	return true
}

func onSystrayReady() {
	systray.SetIcon(trayIcon)
	systray.SetTitle("CATendar")

	showItem := systray.AddMenuItem("Show Window", "Show CATendar window")
	quitItem := systray.AddMenuItem("Quit", "Exit CATendar")

	go func() {
		for {
			select {
			case <-showItem.ClickedCh:
				if appCtx != nil {
					runtime.WindowShow(appCtx)
				}
			case <-quitItem.ClickedCh:
				if appCtx != nil {
					runtime.Quit(appCtx)
				}
			}
		}
	}()
}
