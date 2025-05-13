package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"net/http"
	"slark/internal/models"
)

func CreateVercelProject(config models.ProjectConfig, platformData models.PlatformData) (string, error) {
	baseURL := "https://api.vercel.com/v11/projects"

	requestURL := baseURL
	if platformData.TeamId != "" && platformData.TeamId != "team_xxxx" {
		requestURL = fmt.Sprintf("%s?teamId=%s", baseURL, platformData.TeamId)
	}

	projectData := map[string]any{
		"name":                              config.Name,
		"buildCommand":                      nil,
		"commandForIgnoringBuildStep":       nil,
		"devCommand":                        nil,
		"environmentVariables":              []map[string]any{},
		"framework":                         platformData.Framework,
		"installCommand":                    nil,
		"outputDirectory":                   nil,
		"publicSource":                      nil,
		"enableAffectedProjectsDeployments": true,
		"oidcTokenConfig": map[string]any{
			"enabled":    true,
			"issuerMode": "global",
		},
	}

	if strings.Compare(config.BuildFolder, "./") == 0{
		projectData["rootDirectory"] = config.BuildFolder
	}

	jsonData, err := json.Marshal(projectData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal project data: %w", err)
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+platformData.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)

	// Handle response status codes
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		var errorResponse map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
			if errMsg, ok := errorResponse["error"].(map[string]any)["message"]; ok {
				return "", fmt.Errorf("failed to create project, status code: %d, message: %v", resp.StatusCode, errMsg)
			}

		}
		return "", fmt.Errorf("failed to create project, status code: %d", resp.StatusCode)
	}

	fmt.Printf("response body: %s\n", resp.Body)

	return "", nil
}

func DeleteVercelProject(projectName, deployBranch, buildFolder string) error {
	return nil
}

func GetVercelProject(projectName, deployBranch, buildFolder string) error {
	return nil
}
