package allowlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Allowlist []string

func PullAllowlist() (Allowlist, error) {
	// Hardcoded URL for the Lilypad-Tech module allowlist
	url := "https://raw.githubusercontent.com/Lilypad-Tech/module-allowlist/main/allowlist.json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch allowlist: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch allowlist: HTTP %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var allowlist Allowlist
	if err := json.Unmarshal(body, &allowlist); err != nil {
		return nil, fmt.Errorf("failed to unmarshal allowlist: %v", err)
	}

	// Save the allowlist to a local file
	saveDir := filepath.Join(os.TempDir(), "lilypad-allowlist")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	savePath := filepath.Join(saveDir, "allowlist.json")
	if err := ioutil.WriteFile(savePath, body, 0644); err != nil {
		return nil, fmt.Errorf("failed to save allowlist: %v", err)
	}

	log.Info().Msgf("Allowlist saved to: %s", savePath)

	return allowlist, nil
}
