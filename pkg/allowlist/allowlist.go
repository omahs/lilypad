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

type AllowlistJSON struct {
	Modules []string `json:"modules"`
}

type Allowlist []string

func PullAllowlist() (Allowlist, error) {
	// Hardcoded URL for the Lilypad-Tech module allowlist
	url := "https://raw.githubusercontent.com/Lilypad-Tech/module-allowlist/main/allowlist.json"

	// Create a client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the HTTP request
	resp, err := client.Get(url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Failed to make HTTP request")
		return nil, fmt.Errorf("failed to fetch allowlist: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		log.Error().Int("statusCode", resp.StatusCode).Str("url", url).Msg("Received non-OK status code")
		return nil, fmt.Errorf("failed to fetch allowlist: HTTP %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Unmarshal the JSON
	var allowlistJSON AllowlistJSON
	if err := json.Unmarshal(body, &allowlistJSON); err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg("Failed to unmarshal JSON")
		return nil, fmt.Errorf("failed to unmarshal allowlist: %v", err)
	}

	// Check if the modules list is empty
	if len(allowlistJSON.Modules) == 0 {
		log.Warn().Msg("Allowlist is empty")
	}

	// Convert AllowlistJSON to Allowlist
	allowlist := Allowlist(allowlistJSON.Modules)

	// Save the allowlist to a local file
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
