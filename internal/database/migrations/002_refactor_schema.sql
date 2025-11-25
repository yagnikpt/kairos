-- +goose Up
CREATE TABLE goals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    status TEXT DEFAULT 'ACTIVE',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    goal_id INTEGER,
    parent_task_id INTEGER,
    description TEXT NOT NULL,
    status TEXT DEFAULT 'PENDING',
    estimated_duration_mins INTEGER,
    proof_of_work TEXT,
    FOREIGN KEY(goal_id) REFERENCES goals(id) ON DELETE CASCADE,
    FOREIGN KEY(parent_task_id) REFERENCES tasks(id) ON DELETE CASCADE
);

CREATE TABLE app_state (
    key TEXT PRIMARY KEY,
    value TEXT
);

-- +goose Down
DROP TABLE app_state;
DROP TABLE tasks;
DROP TABLE goals;
