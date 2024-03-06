package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	copilotToken    string
	copilotTokenURL string
)

func main() {
	flag.StringVar(&copilotToken, "token", "", "token to get copilot token")
	flag.StringVar(&copilotTokenURL, "token_url", "", "url to get copilot token")
	flag.Parse()
	fmt.Printf("copilotToken: %s\n", copilotToken)
	fmt.Printf("copilotTokenURL: %s\n", copilotTokenURL)
	Init()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/models", handleModel)
	mux.HandleFunc("POST /v1/chat/completions", handleProxy)
	http.ListenAndServe(":8081", mux)
}
