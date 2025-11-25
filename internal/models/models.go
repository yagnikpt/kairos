package models

import (
	"database/sql"
	"time"
)

type Goal struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"` // ACTIVE, ARCHIVED, COMPLETED
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID                    int64          `json:"id"`
	GoalID                int64          `json:"goal_id"`
	ParentTaskID          sql.NullInt64  `json:"parent_task_id"`
	Description           string         `json:"description"`
	Status                string         `json:"status"` // PENDING, IN_PROGRESS, DONE, SKIPPED
	EstimatedDurationMins sql.NullInt64  `json:"estimated_duration_mins"`
	ProofOfWork           sql.NullString `json:"proof_of_work"`
}
