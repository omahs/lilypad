package allowlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

type AllowlistItem struct {
	ModuleId string `json:"ModuleId"`
	Version  string `json:"Version"`
	Enabled  bool   `json:"Enabled"`
}

type Allowlist map[string]string

func PullAllowlist() (Allowlist, error) {
	url := "https://raw.githubusercontent.com/Lilypad-Tech/module-allowlist/main/allowlist.json"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Failed to make HTTP request")
		return nil, fmt.Errorf("failed to fetch allowlist: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("statusCode", resp.StatusCode).Str("url", url).Msg("Received non-OK status code")
		return nil, fmt.Errorf("failed to fetch allowlist: HTTP %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var allowlistItems []AllowlistItem
	if err := json.Unmarshal(body, &allowlistItems); err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg("Failed to unmarshal JSON")
		return nil, fmt.Errorf("failed to unmarshal allowlist: %v", err)
	}

	allowlist := make(Allowlist)
	for _, item := range allowlistItems {
		if item.Enabled {
			allowlist[item.ModuleId] = item.Version
		}
	}

	if len(allowlist) == 0 {
		log.Warn().Msg("Allowlist is empty")
	}

	saveDir := filepath.Join(os.TempDir(), "lilypad-allowlist")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		log.Error().Err(err).Str("dir", saveDir).Msg("Failed to create directory")
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	savePath := filepath.Join(saveDir, "allowlist.json")
	if err := ioutil.WriteFile(savePath, body, 0644); err != nil {
		log.Error().Err(err).Str("path", savePath).Msg("Failed to write allowlist to file")
		return nil, fmt.Errorf("failed to save allowlist: %v", err)
	}

	log.Info().
		Str("path", savePath).
		Int("moduleCount", len(allowlist)).
		Msg("Allowlist saved successfully")

	return allowlist, nil
}
