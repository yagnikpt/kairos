package commands

import (
	"database/sql"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yagnikpt/kairos/internal/app"
	"github.com/yagnikpt/kairos/internal/tui"
	"github.com/yagnikpt/kairos/internal/ui"
)

func NewRootCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kairos",
		Short: "Kairos: Focus on what matters",
		Long:  `Kairos is a CLI tool to help you manage your goals and stay focused.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check for current goal
			var currentGoalIDStr string
			err := a.DB.QueryRow("SELECT value FROM app_state WHERE key = ?", "current_goal_id").Scan(&currentGoalIDStr)
			if err == sql.ErrNoRows {
				ui.RenderSubtitle("No active goal selected. Use 'kairos add' to start or 'kairos switch' to pick one.")
				return
			} else if err != nil {
				ui.RenderError(err)
				return
			}

			currentGoalID, _ := strconv.ParseInt(currentGoalIDStr, 10, 64)
			if err := tui.RunFocusMode(a, currentGoalID); err != nil {
				if err == sql.ErrNoRows {
					ui.RenderSubtitle("No active goal selected. Use 'kairos add' to start or 'kairos switch' to pick one.")
					// Clean up invalid state
					a.DB.Exec("DELETE FROM app_state WHERE key = 'current_goal_id'")
				} else {
					ui.RenderError(err)
				}
				return
			}
			// fmt.Print("\033[H\033[2J")
			// fmt.Println("Keep grinding 💪")
		},
	}

	cmd.AddCommand(newAddCmd(a))
	cmd.AddCommand(newSwitchCmd(a))
	cmd.AddCommand(newChillCmd(a))

	return cmd
}
