package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/yagnikpt/kairos/internal/app"
	"github.com/yagnikpt/kairos/internal/ui"
)

func newAddCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new goal",
		Run: func(cmd *cobra.Command, args []string) {
			var goalName string
			contextInfo, _ := cmd.Flags().GetString("context")

			if len(args) > 0 {
				goalName = strings.Join(args, " ")
			} else {
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewInput().
							Title("What is your goal?").
							// Prompt("I want to... ").
							Value(&goalName),
						huh.NewText().
							Title("Additional Context (Optional)").
							Value(&contextInfo),
					),
				).WithTheme(ui.HuhTheme)

				if err := form.Run(); err != nil {
					ui.RenderError(err)
					return
				}
			}

			if goalName == "" {
				return
			}

			ui.RenderTitle("Analyzing your goal...")
			highLevelTasks, err := a.AI.GenerateHighLevelTasks(goalName, contextInfo)
			if err != nil {
				ui.RenderError(err)
				return
			}

			ui.RenderSubtitle("Proposed milestones:")
			for _, t := range highLevelTasks {
				fmt.Printf("- %s\n", t)
			}

			var confirm bool
			confirmForm := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title("Do these look good?").
						Value(&confirm),
				),
			)

			if err := confirmForm.Run(); err != nil {
				ui.RenderError(err)
				return
			}

			if !confirm {
				ui.RenderSubtitle("Cancelled. Try again.")
				return
			}

			// Save Goal
			res, err := a.DB.Exec("INSERT INTO goals (name, status, created_at) VALUES (?, 'ACTIVE', ?)", goalName, time.Now())
			if err != nil {
				ui.RenderError(err)
				return
			}
			goalID, _ := res.LastInsertId()

			ui.RenderTitle("Generating detailed plan... (this might take a moment)")

			// Generate and Save Tasks
			for _, hlTask := range highLevelTasks {
				// Save High Level Task
				res, err := a.DB.Exec("INSERT INTO tasks (goal_id, description, status) VALUES (?, ?, 'PENDING')", goalID, hlTask)
				if err != nil {
					ui.RenderError(err)
					continue
				}
				hlTaskID, _ := res.LastInsertId()

				// Generate Subtasks
				subTasks, err := a.AI.GenerateSubTasks(hlTask)
				if err != nil {
					ui.RenderError(fmt.Errorf("failed to generate subtasks for '%s': %v", hlTask, err))
					continue
				}

				for _, subTask := range subTasks {
					_, err := a.DB.Exec("INSERT INTO tasks (goal_id, parent_task_id, description, status) VALUES (?, ?, ?, 'PENDING')", goalID, hlTaskID, subTask)
					if err != nil {
						ui.RenderError(err)
					}
				}
			}

			// Set as current goal
			_, err = a.DB.Exec("INSERT OR REPLACE INTO app_state (key, value) VALUES ('current_goal_id', ?)", goalID)
			if err != nil {
				ui.RenderError(err)
			}

			ui.RenderSuccess("Goal setup complete! Run 'kairos' to start working.")
		},
	}
	cmd.Flags().StringP("context", "c", "", "Additional context for the goal")
	return cmd
}
