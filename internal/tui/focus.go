package tui

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/yagnikpt/kairos/internal/app"
	"github.com/yagnikpt/kairos/internal/models"
	"github.com/yagnikpt/kairos/internal/ui"
)

func RunFocusMode(a *app.App, goalID int64) error {
	// Clear screen
	fmt.Print("\033[H\033[2J")

	// Get Goal Name
	var goalName string
	var goalStatus string
	err := a.DB.QueryRow("SELECT name, status FROM goals WHERE id = ?", goalID).Scan(&goalName, &goalStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return err
		}
		return err
	}

	// Header Section
	fmt.Println(ui.BoxStyle.Render(fmt.Sprintf("[ %s ]", goalName)))
	// ui.RenderStatus("STATUS:", goalStatus)
	fmt.Println()

	// Find the first PENDING or IN_PROGRESS high-level task
	var hlTask models.Task
	err = a.DB.QueryRow(`
		SELECT id, description, status
		FROM tasks
		WHERE goal_id = ? AND parent_task_id IS NULL AND status IN ('PENDING', 'IN_PROGRESS')
		ORDER BY id ASC LIMIT 1`, goalID).Scan(&hlTask.ID, &hlTask.Description, &hlTask.Status)

	if err == sql.ErrNoRows {
		// Mark goal as COMPLETED
		_, err := a.DB.Exec("UPDATE goals SET status = 'COMPLETED' WHERE id = ?", goalID)
		if err != nil {
			return err
		}
		ui.RenderSuccess("All milestones completed! Goal marked as COMPLETED.")
		return nil
	} else if err != nil {
		return err
	}

	ui.RenderSubtitle("CURRENT TASK: " + hlTask.Description)
	fmt.Println()

	// Get subtasks for this HL task
	rows, err := a.DB.Query(`
		SELECT id, description, status
		FROM tasks
		WHERE parent_task_id = ?
		ORDER BY id ASC`, hlTask.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var subTasks []models.Task
	var options []huh.Option[int64]

	// Add subtasks to options
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Description, &t.Status); err != nil {
			continue
		}
		subTasks = append(subTasks, t)

		label := fmt.Sprintf("[ ] %s", t.Description)
		switch t.Status {
		case "DONE":
			label = fmt.Sprintf("[x] %s", t.Description)
		case "SKIPPED":
			label = fmt.Sprintf("[-] %s", t.Description)
		}

		options = append(options, huh.NewOption(label, t.ID))
	}

	// Footer Options
	options = append(options, huh.NewOption("---", int64(-1)))
	options = append(options, huh.NewOption("> I'm Exhausted (Switch Context)", int64(-99)))

	var selectedAction int64
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int64]().
				Title("").
				Options(options...).
				Value(&selectedAction).
				// Height(16).
				WithTheme(ui.HuhTheme),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	if selectedAction == -1 {
		return nil // Separator selected, do nothing
	}

	if selectedAction == -99 {
		// Switch Context (Chill Mode)
		ui.RenderSubtitle("Take a break. Run 'kairos chill'.")
		return nil
	}

	// Handle Subtask Selection (Toggle)
	var subTask models.Task
	for _, t := range subTasks {
		if t.ID == selectedAction {
			subTask = t
			break
		}
	}

	if subTask.ID != 0 {
		newStatus := "DONE"
		if subTask.Status == "DONE" {
			newStatus = "PENDING"
		}
		_, err = a.DB.Exec("UPDATE tasks SET status = ? WHERE id = ?", newStatus, selectedAction)
		if err != nil {
			return err
		}

		// Check if all subtasks are DONE
		var pendingCount int
		err = a.DB.QueryRow("SELECT COUNT(*) FROM tasks WHERE parent_task_id = ? AND status != 'DONE'", hlTask.ID).Scan(&pendingCount)
		if err != nil {
			return err
		}

		if pendingCount == 0 {
			// Mark HL task as DONE
			_, err = a.DB.Exec("UPDATE tasks SET status = 'DONE' WHERE id = ?", hlTask.ID)
			if err != nil {
				return err
			}
			ui.RenderSuccess("Milestone completed! Moving to next...")
			// Optional: Sleep briefly to let user see the success message?
			// time.Sleep(1 * time.Second)
		} else {
			// Ensure HL task is IN_PROGRESS
			if hlTask.Status == "PENDING" {
				_, err = a.DB.Exec("UPDATE tasks SET status = 'IN_PROGRESS' WHERE id = ?", hlTask.ID)
				if err != nil {
					return err
				}
			}
		}

		return RunFocusMode(a, goalID)
	}
	return nil
}
