package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
	model  string
}

func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Client{
		client: client,
		model:  "gemini-2.0-flash", // Using a faster model if available, or fallback to user pref
	}, nil
}

func (c *Client) GenerateHighLevelTasks(goal string, contextInfo string) ([]string, error) {
	prompt := fmt.Sprintf(`
You are a productivity assistant.
The user has a goal: "%s".
%s
Break this down into 3-5 high-level, actionable milestones or phases.
Return ONLY a JSON array of strings, where each string is a task description.
Example: ["Learn basic syntax", "Build a small project", "Read documentation"]
`, goal, func() string {
		if contextInfo != "" {
			return fmt.Sprintf("Additional context: %s", contextInfo)
		}
		return ""
	}())

	return c.generateList(prompt)
}

func (c *Client) GenerateSubTasks(parentTask string) ([]string, error) {
	prompt := fmt.Sprintf(`
You are a productivity assistant.
The user has a high-level task: "%s".
Break this down into 3-5 small, actionable sub-tasks that can be done in 15-30 minutes.
Return ONLY a JSON array of strings.
`, parentTask)

	return c.generateList(prompt)
}

func (c *Client) SuggestContent(interests []string) (string, error) {
	prompt := fmt.Sprintf(`
The user needs a break. Their interests are: %s.
Suggest a topic or a type of article/paper they should read to relax but stay inspired.
Keep it short and encouraging.
`, strings.Join(interests, ", "))

	resp, err := c.client.Models.GenerateContent(context.Background(), c.model, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no content generated")
	}

	var textBuilder strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		textBuilder.WriteString(part.Text)
	}
	return textBuilder.String(), nil
}

func (c *Client) generateList(prompt string) ([]string, error) {
	resp, err := c.client.Models.GenerateContent(context.Background(), c.model, genai.Text(prompt), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if resp == nil || len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content generated")
	}

	var textBuilder strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		textBuilder.WriteString(part.Text)
	}
	text := textBuilder.String()

	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var tasks []string
	if err := json.Unmarshal([]byte(text), &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse json response: %w, text: %s", err, text)
	}

	return tasks, nil
}
