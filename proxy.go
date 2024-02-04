package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
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
		} `json:"delta"`
	} `json:"choices"`
}

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
			if r.StatusCode != http.StatusOK {
				return nil
			}
			p.scanStream(r)
			return nil
		}
	}
	p.ReverseProxy.ServeHTTP(w, r)
}

func (p *Proxy) scanStream(r *http.Response) {
	w := &bytes.Buffer{}
	scanner := bufio.NewScanner(r.Body)
	for scanner.Scan() {
		v := scanner.Bytes()
		if len(v) == 0 {
			continue
		}
		if bytes.Contains(v, []byte("data: [DONE]")) {
			w.Write([]byte("data: [DONE]\n"))
			r.Body = io.NopCloser(w)
			return
		}
		if len(v) > 5 && bytes.Contains(v, []byte("data")) {
			var output streamOutput
			if err := json.Unmarshal([]byte(v[5:]), &output); err != nil {
				fmt.Printf("json unmarshal error: %v\n", err)
				continue
			}
			if len(output.Choices) == 0 {
				continue
			}
			for _, c := range output.Choices {
				if c.Delta.Content == "" {
					continue
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
			w.Write([]byte("\n\n"))
		}
	}
}

// 所有请求转发到https://api.githubcopilot.com
func handleProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	p := &Proxy{}

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close() //  must close
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	bodyJson := gjson.ParseBytes(bodyBytes)
	p.stream = bodyJson.Get("stream").Bool()
	p.model = bodyJson.Get("model").String()

	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Host = "api.githubcopilot.com"
			req.URL.Scheme = "https"
			req.URL.Host = "api.githubcopilot.com"
			// remove /v1/ from req.URL.Path
			req.URL.Path = strings.ReplaceAll(req.URL.Path, "/v1/", "/")

			accToken, err := getAccToken()
			if accToken == "" {
				w.WriteHeader(http.StatusInternalServerError)
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
