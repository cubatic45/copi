package main

import (
	"flag"
	"net/http"
)

var (
	copilotToken    string
	copilotTokenURL string
)

func main() {
	copilotToken = *flag.String("token", "", "token to get copilot token")
	copilotTokenURL = *flag.String("tokenurl", "", "url to get copilot token")
	flag.Parse()
	Init()
	http.HandleFunc("/", handleProxy)
	http.ListenAndServe(":8081", nil)
}
