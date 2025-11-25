package commands

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/yagnikpt/kairos/internal/app"
	"github.com/yagnikpt/kairos/internal/ui"
)

func newChillCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "chill",
		Short: "Take a break",
		Run: func(cmd *cobra.Command, args []string) {
			// Header
			fmt.Println(ui.BoxStyle.Render("[ chill mode ]"))
			ui.RenderStatus("STATUS:", "RECHARGE / INGEST")
			fmt.Println()

			// Hardcoded interests for now, could be stored in config/DB
			interests := []string{"Technology", "Science", "Programming", "Hacker News"}

			suggestion, err := a.AI.SuggestContent(interests)
			if err != nil {
				ui.RenderError(err)
				return
			}

			ui.RenderSubtitle("Top Pick from your Queue:")
			fmt.Println(ui.ItemStyle.Render(fmt.Sprintf("\"%s\"", suggestion)))
			ui.RenderStatus("Time:", "15 min read") // Placeholder
			fmt.Println()

			// Options
			var selectedAction int
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[int]().
						Title("").
						Options(
							huh.NewOption("> Open URL", 1),
							huh.NewOption("> Skip (Only allowed 3 skips/day)", 2),
							huh.NewOption("> Switch back to Code", 3),
						).
						Value(&selectedAction),
				),
			)

			if err := form.Run(); err != nil {
				return
			}

			switch selectedAction {
			case 1:
				ui.RenderSuccess("Opening URL... (Simulated)")
				// In real app, use 'open' package
			case 2:
				ui.RenderStatus("Skipped.", "2 skips remaining")
				// Recursive call to show next? Or just exit.
			case 3:
				ui.RenderSuccess("Back to work!")
			}
		},
	}
}
