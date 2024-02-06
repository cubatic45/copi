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
	http.HandleFunc("/", handleProxy)
	http.ListenAndServe(":8081", nil)
}
