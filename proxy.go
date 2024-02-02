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
)

type streamOutput struct {
	ID                string      `json:"id"`
	Object            string      `json:"object"`
	Created           int         `json:"created"`
	Model             string      `json:"model"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
	Choices           []struct {
		Index int `json:"index"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason interface{} `json:"finish_reason"`
	} `json:"choices"`
}

// streamOutput inplement json mashal
func (s streamOutput) MarshalJSON() ([]byte, error) {
	type Alias streamOutput
	return json.Marshal(&struct {
		*Alias
		Object string `json:"object"`
	}{
		Alias:  (*Alias)(&s),
		Object: "chat.completion.chunk",
	})
}

type Proxy struct {
	stream bool
	*httputil.ReverseProxy
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.stream {
		// handle stream
		// ...
		p.ModifyResponse = func(r *http.Response) error {
			scanStream(r.Body, w)
			return nil
		}
	}
	p.ReverseProxy.ServeHTTP(w, r)
}

func scanStream(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		v := scanner.Bytes()
		if len(v) == 0 {
			continue
		}
		if len(v) > 5 && bytes.Contains(v, []byte("data")) {
			var output streamOutput
			if err := json.Unmarshal([]byte(v[5:]), &output); err != nil {
				fmt.Printf("json unmarshal error: %v\n", err)
			}
			jsonOutput, _ := json.Marshal(output)
			w.Write([]byte("data: "))
			w.Write(jsonOutput)
			w.Write([]byte("\n"))
			continue
		}
		w.Write(v)
		w.Write([]byte("\n"))
	}
}

// 所有请求转发到https://api.githubcopilot.com
func handleProxy(w http.ResponseWriter, r *http.Request) {
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
	p := &Proxy{stream: true, ReverseProxy: &proxy}
	p.ServeHTTP(w, r)
}
