package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"catendar-wails/database"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var appCtx context.Context

func main() {
	workDir := getWorkDir()
	frontendAssets := os.DirFS(filepath.Join(workDir, "frontend", "dist"))

	// Start systray event loop in a goroutine (locks OS thread internally)
	go systrayLoop()

	err := wails.Run(&options.App{
		Title:     "CATendar",
		Width:     1000,
		Height:    700,
		MinWidth:  360,
		MinHeight: 280,
		Assets:            frontendAssets,
		HideWindowOnClose: true,
		OnStartup: func(ctx context.Context) {
			appCtx = ctx
			dbPath := filepath.Join(workDir, "..", "backend", "calendar.db")
			if err := database.InitDB(dbPath); err != nil {
				log.Fatalf("Failed to initialize database: %v", err)
			}
			log.Println("[CATendar] Database initialized")
		},
		OnDomReady: func(ctx context.Context) {
			log.Println("[CATendar] App ready")
		},
		OnShutdown: func(ctx context.Context) {
			database.Close()
		},
		Bind: []interface{}{
			&CalendarApp{},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func systrayLoop() {
	// getlantern/systray requires LockOSThread via init(), so we need to call it in a dedicated goroutine
	systray.Run(onSystrayReady, onSystrayExit)
}

func onSystrayReady() {
	iconPath := filepath.Join(getWorkDir(), "build", "windows", "icon.ico")
	iconBytes, err := os.ReadFile(iconPath)
	if err != nil {
		log.Printf("[CATendar] Failed to load tray icon: %v", err)
	} else {
		systray.SetIcon(iconBytes)
	}
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
				database.Close()
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func onSystrayExit() {
	// Called when systray quits
}

func getWorkDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	dir := filepath.Dir(exe)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}
