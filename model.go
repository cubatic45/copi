package main

import "net/http"

func handleModel(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"object":"list","data":[{"id":"gpt-3.5-turbo","object":"model","created":1677610602,"owned_by":"openai"},{"id":"gpt-4","object":"model","created":1677610602,"owned_by":"openai"}]}`))
}
