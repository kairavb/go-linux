package cmd

import (
	"context"
	"fmt"
	"strings"

	genai "google.golang.org/genai"
)

func GetCommandFromAI(userPrompt string) (string, string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  "", // if empty, it will use GEMINI_API_KEY env var
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to create genai client: %v", err)
	}

	fullPrompt := fmt.Sprintf(`You are a helpful Linux assistant.
Translate the following natural language instruction into a safe Linux command.
Then explain what the command does in one sentence.

Instruction: %s

Format your reply exactly like:
COMMAND: <command>
EXPLANATION: <what it does>`, userPrompt)

	// Build contents (user message) using helper
	contents := []*genai.Content{
		genai.NewContentFromText(fullPrompt, genai.RoleUser),
	}

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, nil)
	if err != nil {
		return "", "", fmt.Errorf("error generating content: %v", err)
	}

	// Extract the text from response

	contentText := resp.Text()

	// Debug: optionally log raw contentText
	// fmt.Println("Raw AI reply:\n", contentText)

	lines := strings.Split(contentText, "\n")
	var command, explanation string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "COMMAND:") {
			command = strings.TrimSpace(strings.TrimPrefix(line, "COMMAND:"))
		} else if strings.HasPrefix(line, "EXPLANATION:") {
			explanation = strings.TrimSpace(strings.TrimPrefix(line, "EXPLANATION:"))
		}
	}

	if command == "" {
		return "", "", fmt.Errorf("failed to parse command from AI response:\n%s", contentText)
	}

	return command, explanation, nil
}
