package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func getAccHeaders(accessToken string) map[string]string {
	sessionId := fmt.Sprintf("%s%d", uuid.New().String(), time.Now().UnixNano()/int64(time.Millisecond))
	machineID := sha256.Sum256([]byte(uuid.New().String()))
	machineIDStr := hex.EncodeToString(machineID[:])
	return map[string]string{
		"Host":                   "api.githubcopilot.com",
		"Authorization":          "Bearer " + accessToken,
		"X-Request-Id":           uuid.New().String(),
		"X-Github-Api-Version":   "2023-07-07",
		"Vscode-Sessionid":       sessionId,
		"Vscode-machineid":       machineIDStr,
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
