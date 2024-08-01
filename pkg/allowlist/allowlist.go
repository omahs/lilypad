package allowlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lilypad-tech/lilypad/pkg/data"
	"github.com/rs/zerolog/log"
)

// Allowlist represents a map of module IDs to their allowed versions
type Allowlist map[string]string

// GlobalAllowlist is the global allowlist variable
var GlobalAllowlist Allowlist

type AllowlistItem struct {
	ModuleId string `json:"ModuleId"`
	Version  string `json:"Version"`
	Enabled  bool   `json:"Enabled"`
}

func PullAllowlist() error {
	url := "https://raw.githubusercontent.com/Lilypad-Tech/module-allowlist/main/allowlist.json"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("Failed to make HTTP request")
		return fmt.Errorf("failed to fetch allowlist: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Int("statusCode", resp.StatusCode).Str("url", url).Msg("Received non-OK status code")
		return fmt.Errorf("failed to fetch allowlist: HTTP %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var allowlistItems []AllowlistItem
	if err := json.Unmarshal(body, &allowlistItems); err != nil {
		log.Error().Err(err).Str("body", string(body)).Msg("Failed to unmarshal JSON")
		return fmt.Errorf("failed to unmarshal allowlist: %v", err)
	}

	GlobalAllowlist = make(Allowlist)
	for _, item := range allowlistItems {
		if item.Enabled {
			GlobalAllowlist[item.ModuleId] = item.Version
		}
	}

	if len(GlobalAllowlist) == 0 {
		log.Warn().Msg("Allowlist is empty")
	}

	saveDir := filepath.Join(os.TempDir(), "lilypad-allowlist")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		log.Error().Err(err).Str("dir", saveDir).Msg("Failed to create directory")
		return fmt.Errorf("failed to create directory: %v", err)
	}

	savePath := filepath.Join(saveDir, "allowlist.json")
	if err := ioutil.WriteFile(savePath, body, 0644); err != nil {
		log.Error().Err(err).Str("path", savePath).Msg("Failed to write allowlist to file")
		return fmt.Errorf("failed to save allowlist: %v", err)
	}

	log.Info().
		Str("path", savePath).
		Int("moduleCount", len(GlobalAllowlist)).
		Msg("Allowlist saved successfully")

	return nil
}

// IsModuleAllowed checks if a given module is allowed based on the allowlist
func IsModuleAllowed(module data.ModuleConfig) bool {
	moduleID := module.Name // Assuming ModuleConfig has a Name field that represents the module ID

	allowedVersion, isAllowed := GlobalAllowlist[moduleID]
	if !isAllowed {
		log.Debug().
			Str("module", moduleID).
			Msg("module not in allowlist")
		return false
	}

	moduleVersion := extractVersion(module)
	if moduleVersion == "" {
		log.Error().Interface("module", module).Msg("unable to extract version from module")
		return false
	}

	if compareVersions(moduleVersion, allowedVersion) < 0 {
		log.Debug().
			Str("module", moduleID).
			Str("allowedVersion", allowedVersion).
			Str("moduleVersion", moduleVersion).
			Msg("module version is less than allowed version")
		return false
	}

	return true
}

func extractVersion(module data.ModuleConfig) string {
	parts := strings.Split(module.Name, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func compareVersions(v1, v2 string) int {
	v1Parts := strings.Split(strings.TrimPrefix(v1, "v"), ".")
	v2Parts := strings.Split(strings.TrimPrefix(v2, "v"), ".")

	for i := 0; i < len(v1Parts) && i < len(v2Parts); i++ {
		n1, err1 := strconv.Atoi(v1Parts[i])
		n2, err2 := strconv.Atoi(v2Parts[i])

		if err1 != nil || err2 != nil {
			// If we can't parse the version numbers, fall back to string comparison
			if v1Parts[i] < v2Parts[i] {
				return -1
			} else if v1Parts[i] > v2Parts[i] {
				return 1
			}
			continue
		}

		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}

	if len(v1Parts) < len(v2Parts) {
		return -1
	} else if len(v1Parts) > len(v2Parts) {
		return 1
	}

	return 0
}
