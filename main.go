package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"
	"sync"

	"catendar-wails/core"
	synccli "catendar-wails/core/ai_sync/cli"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var appCtx context.Context

//go:embed all:frontend/dist
var frontendAssets embed.FS

func main() {
	if len(os.Args) > 1 && os.Args[1] == "cli" {
		os.Exit(synccli.Run(context.Background(), os.Args[2:], os.Stdin, os.Stdout, os.Stderr))
	}

	dbPath, err := databasePath()
	if err != nil {
		log.Fatalf("Failed to prepare database directory: %v", err)
	}
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to locate CATendar executable: %v", err)
	}
	backend, err := core.Open(dbPath, core.Options{ExecutablePath: executablePath})
	if err != nil {
		log.Fatalf("Failed to initialize CATendar Core: %v", err)
	}
	calendarApp := NewCalendarApp(backend.Calendar)
	aiSyncApp := NewAISyncApp(backend.Email, backend.AISync)
	var shutdownOnce sync.Once
	shutdown := func() {
		shutdownOnce.Do(func() {
			stopSystray()
			if err := backend.Close(); err != nil {
				log.Printf("[CATendar] Shutdown error: %v", err)
			}
		})
	}
	defer shutdown()

	startSystray()
	log.Println("[CATendar] Core initialized")

	err = wails.Run(&options.App{
		Title:             "CATendar",
		Width:             1000,
		Height:            700,
		MinWidth:          360,
		MinHeight:         280,
		AssetServer:       &assetserver.Options{Assets: frontendAssets},
		HideWindowOnClose: hideWindowOnClose(),
		OnStartup: func(ctx context.Context) {
			appCtx = ctx
			aiSyncApp.startup(ctx)
		},
		OnDomReady: func(ctx context.Context) {
			log.Println("[CATendar] App ready")
		},
		OnShutdown: func(ctx context.Context) {
			shutdown()
		},
		Bind: []interface{}{
			calendarApp,
			aiSyncApp,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func databasePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(configDir, "CATendar")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "calendar.db"), nil
}
