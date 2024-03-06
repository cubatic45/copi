package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
)

var (
	caches     = new(cache.Cache)
	copilotCFG = new(copilot)
)

func getCopilot() copiloter {
	return copilotCFG
}

type copilot struct {
	GithubCom struct {
		User        string `json:"user"`
		OauthToken  string `json:"oauth_token"`
		DevOverride struct {
			CopilotTokenURL string `json:"copilot_token_url"`
		} `json:"dev_override"`
	} `json:"github.com"`
}

func Init() {
	// init cache
	one := &sync.Once{}
	one.Do(func() {
		caches = cache.New(15*time.Minute, 20*time.Minute)
	})

	// ~/.config/github-copilot/hosts.json
	// get user home directory
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	hostsJson := path.Join(dir, ".config", "github-copilot", "hosts.json")
	_, err = os.Stat(hostsJson)

	if !os.IsNotExist(err) {
		bts, err := os.ReadFile(hostsJson)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(bts, copilotCFG); err != nil {
			log.Fatal(err)
		}
		log.Printf("read config file path: %s", hostsJson)
	} else {
		if copilotToken == "" {
			log.Fatal("no token and no hosts.json found, please set token or create hosts.json file in ~/.config/github-copilot/hosts.json")
		}
	}
	if copilotToken != "" {
		copilotCFG.GithubCom.OauthToken = copilotToken
	}
	if copilotTokenURL != "" {
		copilotCFG.GithubCom.DevOverride.CopilotTokenURL = copilotTokenURL
	}
}

// Copiloter is the interface that wraps the token method.
// token return the access token for github copilot
type copiloter interface {
	token() (string, error)
	refresh() (string, error)
}

func (c *copilot) refresh() (string, error) {
	caches.Delete(c.GithubCom.OauthToken)
	return c.token()
}

func (c *copilot) token() (string, error) {
	if cacheToken, ok := caches.Get(c.GithubCom.OauthToken); ok {
		return cacheToken.(string), nil
	}
	tokenURL := c.GithubCom.DevOverride.CopilotTokenURL
	if tokenURL == "" {
		tokenURL = "https://api.github.com/copilot_internal/v2/token"
	}

	req, err := http.NewRequest(http.MethodGet, tokenURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.GithubCom.OauthToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get token error: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if token := gjson.GetBytes(body, "token").String(); token != "" {
		caches.Set(c.GithubCom.OauthToken, token, 14*time.Minute)
		return token, nil
	}

	return "", fmt.Errorf("get token error")
}
