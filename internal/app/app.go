package app

import (
	"database/sql"

	"github.com/yagnikpt/kairos/internal/ai"
	"github.com/yagnikpt/kairos/internal/config"
)

type App struct {
	DB     *sql.DB
	AI     *ai.Client
	Config *config.Config
}
