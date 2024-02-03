package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

type streamOutput struct {
	ID                string `json:"id,omitempty"`
	Object            string `json:"object"`
	Created           int    `json:"created,omitempty"`
	Model             string `json:"model"`
	SystemFingerprint any    `json:"system_fingerprint,omitempty"`
	Choices           []struct {
		Index int `json:"index"`
		Delta struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"delta"`
		Logprobs     any `json:"logprobs,omitempty"`
		FinishReason any `json:"finish_reason"`
	} `json:"choices"`
}

// streamOutput inplement json mashal
// func (s streamOutput) MarshalJSON() ([]byte, error) {
// 	type Alias streamOutput
// 	return json.Marshal(&struct {
// 		*Alias
// 		Object string `json:"object"`
// 		Model  string `json:"model"`
// 	}{
// 		Alias:  (*Alias)(&s),
// 		Object: "chat.completion.chunk",
// 	})
// }

type Proxy struct {
	model  string
	stream bool
	*httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.stream {
		// handle stream
		// ...
		p.ModifyResponse = func(r *http.Response) error {
			p.scanStream(r.Body, w)
			return nil
		}
	}
	p.ReverseProxy.ServeHTTP(w, r)
}

func (p *Proxy) scanStream(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		v := scanner.Bytes()
		if len(v) == 0 {
			continue
		}
		if bytes.Contains(v, []byte("data: [DONE]")) {
			w.Write([]byte("data: [DONE]\n"))
			return
		}
		if len(v) > 5 && bytes.Contains(v, []byte("data")) {
			var output streamOutput
			if err := json.Unmarshal([]byte(v[5:]), &output); err != nil {
				fmt.Printf("json unmarshal error: %v\n", err)
			}
			if len(output.Choices) == 0 {
				continue
			}
			for i, c := range output.Choices {
				if c.Delta.Content == "" {
					continue
				}
				if c.Delta.Role == "" {
					output.Choices[i].Delta.Role = "assistant"
				}
			}
			if output.Model == "" {
				output.Model = p.model
			}
			output.Object = "chat.completion"
			if p.stream {
				output.Object = "chat.completion.chunk"
			}
			jsonOutput, _ := json.Marshal(output)
			w.Write([]byte("data: "))
			w.Write(jsonOutput)
			w.Write([]byte("\n"))
			continue
		}
	}
}

// 所有请求转发到https://api.githubcopilot.com
func handleProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "post method not allowed")
		return
	}
	p := &Proxy{}
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close() //  must close
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyJson := gjson.ParseBytes(bodyBytes)
	if bodyJson.Get("stream").Bool() {
		p.stream = true
	}
	p.model = bodyJson.Get("model").String()

	target, _ := url.Parse("https://api.githubcopilot.com")
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = target.Host
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			// remove /v1/ from req.URL.Path
			req.URL.Path = strings.ReplaceAll(req.URL.Path, "/v1/", "/")

			accToken, err := getAccToken()
			if accToken == "" {
				fmt.Fprintf(w, "get acc token error: %v", err)
				return
			}
			accHeaders := getAccHeaders(accToken)
			for k, v := range accHeaders {
				req.Header.Set(k, v)
			}
		},
	}
	p.ReverseProxy = &proxy
	p.ServeHTTP(w, r)
}
