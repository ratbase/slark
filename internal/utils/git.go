package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetGitHubRepoInfo retrieves the current GitHub repository information
// Returns the repository in "username/repo" format
func GetGitHubRepoInfo() (string, error) {
	// Get the remote URL
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git remote: %w", err)
	}

	remoteURL := strings.TrimSpace(string(output))

	// Parse the GitHub repository from the URL
	// Handle formats like:
	// https://github.com/username/repo.git
	// git@github.com:username/repo.git
	var repo string

	if strings.HasPrefix(remoteURL, "https://github.com/") {
		// HTTPS format
		repo = strings.TrimPrefix(remoteURL, "https://github.com/")
		repo = strings.TrimSuffix(repo, ".git")
	} else if strings.HasPrefix(remoteURL, "git@github.com:") {
		// SSH format
		repo = strings.TrimPrefix(remoteURL, "git@github.com:")
		repo = strings.TrimSuffix(repo, ".git")
	} else {
		return "", fmt.Errorf("unsupported git remote format: %s", remoteURL)
	}

	return repo, nil
}

// GetCurrentBranch returns the name of the current git branch
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
