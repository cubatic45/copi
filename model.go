package main

import "net/http"

func handleModel(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    data := `{"object":"list","data":[{"id":"gpt-3.5-turbo","object":"model","created":1626777600,"owned_by":"openai","permission":[{"id":"modelperm-LwHkVFn8AcMItP432fKKDIKJ","object":"model_permission","created":1626777600,"allow_create_engine":true,"allow_sampling":true,"allow_logprobs":true,"allow_search_indices":false,"allow_view":true,"allow_fine_tuning":false,"organization":"*","group":null,"is_blocking":false}],"root":"gpt-3.5-turbo","parent":null},{"id":"gpt-4","object":"model","created":1626777600,"owned_by":"openai","permission":[{"id":"modelperm-LwHkVFn8AcMItP432fKKDIKJ","object":"model_permission","created":1626777600,"allow_create_engine":true,"allow_sampling":true,"allow_logprobs":true,"allow_search_indices":false,"allow_view":true,"allow_fine_tuning":false,"organization":"*","group":null,"is_blocking":false}],"root":"gpt-4","parent":null}]}`
	w.Write([]byte(data))
}
