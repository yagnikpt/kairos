# Kairos

**Kairos** is a CLI productivity tool designed to help you focus on what matters. It breaks down your goals into actionable milestones and tasks using AI, keeping you in the flow with a distraction-free TUI.

**NOTE**: This project is vibe coded and built mainly for my personal use.

## Features

- **AI-Powered Planning**: Automatically breaks down goals into high-level milestones and subtasks.
- **Focus Mode**: A Zen/Retro TUI to keep you focused on the current task.
- **Auto-Advance**: Automatically moves to the next milestone when tasks are completed.
- **Context Switching**: Manage multiple goals and switch between them easily.
- **Chill Mode**: AI-curated content suggestions for your breaks.

## Installation

```bash
go install github.com/yagnikpt/kairos/cmd/kairos@latest
```

## Usage

### Start a New Goal
```bash
kairos add Learn Rust -c "Focus on memory safety and concurrency"
```

### Focus Mode
Run the tool to enter the focus view for your active goal:
```bash
kairos
```

### Switch Goals
```bash
kairos switch
```
(Use `d` to delete a goal)

### Take a Break (not implemented yet)
```bash
kairos chill
```

## Configuration

Kairos requires a Google Gemini API key.
Set it in `~/.config/kairos/config.yaml` or via environment variable `GEMINI_API_KEY`.

## License

MIT
