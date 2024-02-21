## cli

```
./copi 
or
./copi -token your_token -token_url https://api.github.com/copilot_internal/v2/token
```

## docker

```bash
docker run --rm -id -p 8081:8081 ghcr.io/cubatic45/copi:latest -token your_token -token_url https://api.github.com/copilot_internal/v2/token

or

docker run --rm -it -p 8081:8081 -v"$(echo $HOME)".config/github-copilot:/root/.config/github-copilot ghcr.io/cubatic45/copi:latest
```

