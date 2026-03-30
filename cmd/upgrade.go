package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/NoaTamburrini/portman/internal/version"
)

func executeUpgrade() {
	fmt.Println("Checking for updates...")

	info := version.GetUpdateInfoNow()
	if !info.Available {
		fmt.Printf("You're already on the latest version (v%s)\n", version.Version)
		return
	}

	fmt.Printf("Updating portman: v%s → %s\n", info.CurrentVersion, info.LatestVersion)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command",
			fmt.Sprintf("irm https://raw.githubusercontent.com/%s/%s/main/install.ps1 | iex", version.RepoOwner, version.RepoName))
	} else {
		cmd = exec.Command("sh", "-c",
			fmt.Sprintf("curl -fsSL https://raw.githubusercontent.com/%s/%s/main/install.sh | sh", version.RepoOwner, version.RepoName))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Upgrade failed: %v\n", err)
		os.Exit(1)
	}
}
