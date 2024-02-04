package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
)

var copilotConfig CopilotConfig

type CopilotConfig struct {
	GithubCom struct {
		User        string `json:"user"`
		OauthToken  string `json:"oauth_token"`
		DevOverride struct {
			CopilotTokenURL string `json:"copilot_token_url"`
		} `json:"dev_override"`
	} `json:"github.com"`
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// ~/.config/github-copilot/hosts.json

	// get user home directory
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	hostsJson := dir + "/.config/github-copilot/hosts.json"
	log.Println(hostsJson)
	bts, err := os.ReadFile(hostsJson)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bts, &copilotConfig); err != nil {
		panic(err)
	}
}

func getAccToken() (string, error) {
	cache := cache.New(15*time.Minute, 20*time.Minute)
	if cacheToken, ok := cache.Get(copilotConfig.GithubCom.OauthToken); ok {
		return cacheToken.(string), nil
	}
	req, err := http.NewRequest("GET", copilotConfig.GithubCom.DevOverride.CopilotTokenURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", copilotConfig.GithubCom.OauthToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body failed: %s", err)
	}
	if resp.StatusCode == http.StatusOK {
		if token := gjson.GetBytes(body, "token").String(); token != "" {
			cache.Set(copilotConfig.GithubCom.OauthToken, token, 14*time.Minute)
			return token, nil
		}
	}

	return "", fmt.Errorf("get token error: %d", resp.StatusCode)
}
