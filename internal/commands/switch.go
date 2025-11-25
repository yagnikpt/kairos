package commands

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yagnikpt/kairos/internal/app"
	"github.com/yagnikpt/kairos/internal/ui"
)

type goalItem struct {
	id     int64
	name   string
	status string
}

func (i goalItem) Title() string       { return i.name }
func (i goalItem) Description() string { return i.status }
func (i goalItem) FilterValue() string { return i.name }

type listKeyMap struct {
	delete key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
	}
}

type model struct {
	list   list.Model
	keys   *listKeyMap
	choice *goalItem
	app    *app.App
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.delete):
			if len(m.list.Items()) == 0 {
				return m, nil
			}
			selectedItem := m.list.SelectedItem().(goalItem)

			// Delete from DB
			// Note: We are doing this synchronously for simplicity in this CLI tool.
			// For a larger app, we might want to use a Cmd.
			_, err := m.app.DB.Exec("DELETE FROM goals WHERE id = ?", selectedItem.id)
			if err != nil {
				// In a real app we might want to show an error message
				return m, nil
			}

			// Remove from list
			index := m.list.Index()
			m.list.RemoveItem(index)

			// If list is empty after delete, we might want to show a message or just stay empty
			return m, nil

		case msg.String() == "enter":
			if len(m.list.Items()) > 0 {
				i, ok := m.list.SelectedItem().(goalItem)
				if ok {
					m.choice = &i
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return ui.BoxStyle.Render(m.list.View())
}

func newSwitchCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "switch",
		Short: "Switch to a different goal",
		Run: func(cmd *cobra.Command, args []string) {
			rows, err := a.DB.Query("SELECT id, name, status FROM goals")
			if err != nil {
				ui.RenderError(err)
				return
			}
			defer rows.Close()

			var items []list.Item
			for rows.Next() {
				var id int64
				var name, status string
				if err := rows.Scan(&id, &name, &status); err != nil {
					continue
				}

				items = append(items, goalItem{id: id, name: name, status: status})
			}

			if len(items) == 0 {
				ui.RenderSubtitle("No goals found. Use 'kairos add' to create one.")
				return
			}

			// Setup list
			delegate := list.NewDefaultDelegate()
			delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(ui.PrimaryColor).BorderForeground(ui.PrimaryColor)
			delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(ui.SecondaryColor).BorderForeground(ui.PrimaryColor)

			l := list.New(items, delegate, 0, 0)
			l.Title = "Select a Goal"
			l.Styles.Title = ui.TitleStyle
			l.AdditionalFullHelpKeys = func() []key.Binding {
				return []key.Binding{
					key.NewBinding(
						key.WithKeys("d"),
						key.WithHelp("d", "delete"),
					),
				}
			}
			l.AdditionalShortHelpKeys = func() []key.Binding {
				return []key.Binding{
					key.NewBinding(
						key.WithKeys("d"),
						key.WithHelp("d", "delete"),
					),
				}
			}

			keys := newListKeyMap()
			m := model{list: l, keys: keys, app: a}

			p := tea.NewProgram(m, tea.WithAltScreen())
			finalModel, err := p.Run()
			if err != nil {
				ui.RenderError(err)
				return
			}

			// Check if a choice was made
			if m, ok := finalModel.(model); ok && m.choice != nil {
				selectedGoalID := m.choice.id

				// Update DB
				_, err = a.DB.Exec("INSERT OR REPLACE INTO app_state (key, value) VALUES ('current_goal_id', ?)", selectedGoalID)
				if err != nil {
					ui.RenderError(err)
					return
				}

				tx, _ := a.DB.Begin()
				tx.Exec("UPDATE goals SET status = 'IDLE'")
				tx.Exec("UPDATE goals SET status = 'ACTIVE' WHERE id = ?", selectedGoalID)
				tx.Commit()

				ui.RenderSuccess(fmt.Sprintf("Switched to: %s", m.choice.name))
			}
		},
	}
}
