package main

import (
	"context"
	"fmt"
	"log"
)

// App struct
type App struct {
	ctx  context.Context
	deps *Dependencies
}

// NewApp creates a new App application struct
func NewApp() *App {
	a := &App{}

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %s", err)
	}
	deps, err := BuildDependencies(cfg)
	if err != nil {
		log.Fatalf("Failed to build dependencies: %s", err)
	}
	a.deps = deps
	return a
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	log.Println("Shutting down...")
	//deps.RecurringManager.StopAll()
	//deps.CurrencyService.Close()
	// Close DB connections
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
