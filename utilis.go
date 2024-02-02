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

var (
	copilotConfig CopilotConfig
	client        = &http.Client{}
)

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
	accToken := ""

	cache := cache.New(15*time.Minute, 60*time.Minute)
	cacheToken, found := cache.Get(copilotConfig.GithubCom.OauthToken)
	if found {
		accToken = cacheToken.(string)
	} else {
		req, err := http.NewRequest("GET", copilotConfig.GithubCom.DevOverride.CopilotTokenURL, nil)
		if err != nil {
			return accToken, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("token %s", copilotConfig.GithubCom.OauthToken))

		resp, err := client.Do(req)
		if err != nil {
			return accToken, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return accToken, fmt.Errorf("数据读取失败")
		}
		if resp.StatusCode == http.StatusOK {
			accToken = gjson.GetBytes(body, "token").String()
			if accToken == "" {
				return accToken, fmt.Errorf("acc_token 未返回")
			}
			cache.Set(copilotConfig.GithubCom.OauthToken, accToken, 14*time.Minute)
		} else {
			log.Printf("获取 acc_token 请求失败：%d, %s ", resp.StatusCode, string(body))
			return accToken, fmt.Errorf("获取 acc_token 请求失败： %d", resp.StatusCode)
		}
	}
	return accToken, nil
}
