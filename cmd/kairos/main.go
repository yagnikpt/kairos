package main

import (
	"fmt"
	"os"

	"github.com/yagnikpt/kairos/internal/ai"
	"github.com/yagnikpt/kairos/internal/app"
	"github.com/yagnikpt/kairos/internal/commands"
	"github.com/yagnikpt/kairos/internal/config"
	"github.com/yagnikpt/kairos/internal/database"
	"github.com/yagnikpt/kairos/internal/ui"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		ui.RenderError(fmt.Errorf("failed to load config: %w", err))
		os.Exit(1)
	}

	db, err := database.InitDB(cfg.DBPath)
	if err != nil {
		ui.RenderError(fmt.Errorf("failed to init db: %w", err))
		os.Exit(1)
	}
	defer db.Close()

	aiClient, err := ai.NewClient(cfg.GeminiAPIKey)
	if err != nil {
		ui.RenderError(fmt.Errorf("failed to init ai client: %w", err))
		os.Exit(1)
	}

	app := &app.App{
		DB:     db,
		AI:     aiClient,
		Config: cfg,
	}

	rootCmd := commands.NewRootCmd(app)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
