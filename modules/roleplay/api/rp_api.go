package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ulog "github.com/S42yt/tuubaa-bot/utils/logger"
)

type gifResponse struct {
	URL string `json:"url"`
}

func GetGifURL(kind string) (string, error) {
	cli := &http.Client{Timeout: 8 * time.Second}
	url := fmt.Sprintf("https://api.otakugifs.xyz/gif?reaction=%s", kind)
	ulog.Debug("Fetching GIF for kind=%s url=%s", kind, url)

	resp, err := cli.Get(url)
	if err != nil {
		ulog.Error("GetGifURL: HTTP get failed for kind=%s: %v", kind, err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ulog.Warn("GetGifURL: unexpected status %d for kind=%s", resp.StatusCode, kind)
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var g gifResponse
	if err := json.NewDecoder(resp.Body).Decode(&g); err != nil {
		ulog.Error("GetGifURL: json decode failed for kind=%s: %v", kind, err)
		return "", err
	}
	ulog.Debug("GetGifURL: got gif url=%s for kind=%s", g.URL, kind)
	return g.URL, nil
}
