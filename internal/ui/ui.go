package ui

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Zen / Retro Palette
	// Muted pastels, beige/sand background feel (simulated on terminal), soft accents.
	PrimaryColor   = lipgloss.Color("#D8A657") // Muted Gold/Yellow
	SecondaryColor = lipgloss.Color("#A9B665") // Sage Green
	AccentColor    = lipgloss.Color("#EA6962") // Soft Red/Coral
	TextColor      = lipgloss.Color("#D4BE98") // Sand/Beige text
	SubTextColor   = lipgloss.Color("#928374") // Greyish Brown
	FaintColor     = lipgloss.Color("#504945") // Darker Grey

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(SubTextColor)

	StatusStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Bold(true)

	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(TextColor)

	SelectedStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(PrimaryColor).
			Bold(true)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SubTextColor).
			Padding(1, 2)

	// Huh Theme
	HuhTheme = huh.ThemeBase()
)

func init() {
	HuhTheme.Focused.Base = HuhTheme.Focused.Base.BorderForeground(PrimaryColor)
	HuhTheme.Focused.Title = HuhTheme.Focused.Title.Foreground(PrimaryColor)
	HuhTheme.Focused.NoteTitle = HuhTheme.Focused.NoteTitle.Foreground(SecondaryColor)
	HuhTheme.Focused.Directory = HuhTheme.Focused.Directory.Foreground(PrimaryColor)
	HuhTheme.Focused.Description = HuhTheme.Focused.Description.Foreground(SubTextColor)
	HuhTheme.Focused.ErrorIndicator = HuhTheme.Focused.ErrorIndicator.Foreground(AccentColor)
	HuhTheme.Focused.ErrorMessage = HuhTheme.Focused.ErrorMessage.Foreground(AccentColor)
	HuhTheme.Focused.SelectSelector = HuhTheme.Focused.SelectSelector.Foreground(PrimaryColor)
	HuhTheme.Focused.Option = HuhTheme.Focused.Option.Foreground(TextColor)
	HuhTheme.Focused.MultiSelectSelector = HuhTheme.Focused.MultiSelectSelector.Foreground(PrimaryColor)
	HuhTheme.Focused.SelectedOption = HuhTheme.Focused.SelectedOption.Foreground(PrimaryColor)
	HuhTheme.Focused.TextInput.Cursor = HuhTheme.Focused.TextInput.Cursor.Foreground(PrimaryColor)
	HuhTheme.Focused.TextInput.Placeholder = HuhTheme.Focused.TextInput.Placeholder.Foreground(SubTextColor)
	HuhTheme.Focused.TextInput.Prompt = HuhTheme.Focused.TextInput.Prompt.Foreground(PrimaryColor)

	HuhTheme.Blurred = HuhTheme.Focused
	HuhTheme.Blurred.Base.BorderForeground(SubTextColor)
	HuhTheme.Blurred.Title = HuhTheme.Blurred.Title.Foreground(SubTextColor)
	HuhTheme.Blurred.TextInput.Prompt = HuhTheme.Blurred.TextInput.Prompt.Foreground(SubTextColor)
	HuhTheme.Blurred.TextInput.Text = HuhTheme.Blurred.TextInput.Text.Foreground(SubTextColor)
}

func RenderTitle(text string) {
	fmt.Println(TitleStyle.Render(text))
}

func RenderSubtitle(text string) {
	fmt.Println(SubtitleStyle.Render(text))
}

func RenderError(err error) {
	fmt.Println(lipgloss.NewStyle().Foreground(AccentColor).Render(fmt.Sprintf("Error: %v", err)))
}

func RenderSuccess(text string) {
	fmt.Println(lipgloss.NewStyle().Foreground(SecondaryColor).Render(text))
}

func RenderStatus(label, value string) {
	fmt.Printf("%s %s\n", SubtitleStyle.Render(label), StatusStyle.Render(value))
}
