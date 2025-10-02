package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	RepoOwner   = "NoaTamburrini"
	RepoName    = "portman"
	CheckPeriod = 24 * time.Hour
)

// Version is set via ldflags during build
var Version = "dev"

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

// CheckForUpdate checks if a newer version is available and prints a message
func CheckForUpdate() {
	// Check if we should skip (last check was recent)
	if shouldSkipCheck() {
		return
	}

	// Check GitHub for latest release
	latestVersion, err := getLatestVersion()
	if err != nil {
		// Silently fail - don't bother user with network errors
		return
	}

	// Update last check time
	updateLastCheck()

	// Compare versions
	if latestVersion != "" && isNewer(latestVersion, Version) {
		fmt.Fprintf(os.Stderr, "\n⚠️  Update available: %s (current: v%s)\n", latestVersion, Version)
		fmt.Fprintf(os.Stderr, "Run: curl -fsSL https://raw.githubusercontent.com/%s/%s/main/install.sh | sh\n\n", RepoOwner, RepoName)
	}
}

func getLatestVersion() (string, error) {
	client := &http.Client{Timeout: 3 * time.Second}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func shouldSkipCheck() bool {
	cacheFile := getCacheFile()
	info, err := os.Stat(cacheFile)
	if err != nil {
		return false // No cache file, should check
	}

	return time.Since(info.ModTime()) < CheckPeriod
}

func updateLastCheck() {
	cacheFile := getCacheFile()
	os.MkdirAll(filepath.Dir(cacheFile), 0755)
	os.WriteFile(cacheFile, []byte(time.Now().Format(time.RFC3339)), 0644)
}

func getCacheFile() string {
	cacheDir, _ := os.UserCacheDir()
	return filepath.Join(cacheDir, "portman", "last_update_check")
}

// isNewer compares two semantic versions and returns true if latest > current
func isNewer(latest, current string) bool {
	// Remove 'v' prefix if present
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	// Split into parts
	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	// Compare each part
	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		var latestNum, currentNum int
		fmt.Sscanf(latestParts[i], "%d", &latestNum)
		fmt.Sscanf(currentParts[i], "%d", &currentNum)

		if latestNum > currentNum {
			return true
		} else if latestNum < currentNum {
			return false
		}
	}

	// If all parts are equal, not newer
	return false
}
