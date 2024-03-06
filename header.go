package main

import (
	"github.com/google/uuid"
)

func getAccHeaders(accessToken *vscodeCopilot) map[string]string {
	return map[string]string{
		"Host":                   "api.githubcopilot.com",
		"Authorization":          "Bearer " + accessToken.token,
		"X-Request-Id":           uuid.New().String(),
		"X-Github-Api-Version":   "2023-07-07",
		"Vscode-Sessionid":       accessToken.sessionid,
		"Vscode-machineid":       accessToken.machineid,
		"Editor-Version":         "vscode/1.85.1",
		"Editor-Plugin-Version":  "copilot-chat/0.11.1",
		"Openai-Organization":    "github-copilot",
		"Openai-Intent":          "conversation-panel",
		"Content-Type":           "application/json",
		"User-Agent":             "GitHubCopilotChat/0.11.1",
		"Copilot-Integration-Id": "vscode-chat",
		"Accept":                 "*/*",
		"Accept-Encoding":        "gzip, deflate, br",
	}
}
