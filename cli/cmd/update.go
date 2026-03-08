package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/formatho/agent-todo/cli/config"
	"github.com/inconshreveable/go-update"
	"github.com/spf13/cobra"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the CLI to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		checkOnly, _ := cmd.Flags().GetBool("check")

		// Fetch latest release from GitHub
		latestRelease, err := getLatestRelease()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		currentVersion := config.GetVersion()
		latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")

		fmt.Printf("Current version: %s\n", currentVersion)
		fmt.Printf("Latest version:  %s\n", latestVersion)

		if currentVersion == latestVersion {
			fmt.Println("✓ You're already on the latest version!")
			return nil
		}

		fmt.Printf("\n✓ New version available: %s\n", latestVersion)

		if checkOnly {
			return nil
		}

		// Perform the update
		fmt.Println("\nDownloading update...")

		// Determine the correct binary URL for this platform
		goos := runtime.GOOS
		goarch := runtime.GOARCH

		var assetName string
		if goos == "windows" {
			assetName = fmt.Sprintf("agent-todo_%s_%s_%s.exe", latestVersion, goos, goarch)
		} else {
			assetName = fmt.Sprintf("agent-todo_%s_%s_%s", latestVersion, goos, goarch)
		}

		// Get the download URL for the asset
		downloadURL, err := getAssetDownloadURL(latestRelease.TagName, assetName)
		if err != nil {
			return fmt.Errorf("failed to get download URL: %w", err)
		}

		// Download and apply update
		resp, err := http.Get(downloadURL)
		if err != nil {
			return fmt.Errorf("failed to download update: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download update: HTTP %d", resp.StatusCode)
		}

		// Get the executable path
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		// Apply the update
		err = update.Apply(resp.Body, update.Options{
			TargetPath: execPath,
		})
		if err != nil {
			return fmt.Errorf("failed to apply update: %w", err)
		}

		fmt.Printf("✓ Successfully updated to version %s\n", latestVersion)
		fmt.Println("\nPlease restart the CLI to use the new version.")
		return nil
	},
}

func getLatestRelease() (*GitHubRelease, error) {
	resp, err := http.Get("https://api.github.com/repos/formatho/agent-todo/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

func getAssetDownloadURL(tag, assetName string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/formatho/agent-todo/releases/tags/%s", tag)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release struct {
		Assets []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.URL, nil
		}
	}

	return "", fmt.Errorf("asset not found: %s", assetName)
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("check", "c", false, "Only check for updates, don't install")
}
