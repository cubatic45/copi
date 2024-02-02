package main

import "net/http"

func main() {
	http.HandleFunc("/", handleProxy)
	http.ListenAndServe(":8081", nil)
}
