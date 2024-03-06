package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

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
	Director func(*http.Request)
	model    string
	stream   bool
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	write := io.MultiWriter(w, os.Stdout)
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(write, "http request error: %v\n", err)
	}
	if p.Director != nil {
		p.Director(req)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(write, "http client error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(write, "http status error: %d\nbody: %s\n", resp.StatusCode, string(body))
		go getCopilot().refresh()
		return
	}
	if p.stream {
		p.scanStream(w, resp)
		return
	}
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}
	io.Copy(w, resp.Body)
}

func (p *Proxy) scanStream(w http.ResponseWriter, r *http.Response) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(r.Body)
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
			flusher.Flush()
			time.Sleep(10 * time.Millisecond)
		}
	}
}

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

	p.Director = func(req *http.Request) {
		req.Host = "api.githubcopilot.com"
		req.URL.Scheme = "https"
		req.URL.Host = "api.githubcopilot.com"
		// remove /v1/ from req.URL.Path
		req.URL.Path = strings.ReplaceAll(req.URL.Path, "/v1/", "/")

		accToken, err := getCopilot().token()
		if err != nil {
			accToken, _ = getCopilot().refresh()
		}
		if accToken == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "get acc token error: %v", err)
			return
		}
		accHeaders := getAccHeaders(accToken)
		for k, v := range accHeaders {
			req.Header.Set(k, v)
		}
	}

	p.ServeHTTP(w, r)
}
