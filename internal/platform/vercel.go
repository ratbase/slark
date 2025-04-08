package platform

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slark/internal/models"
)

func CreateVercelProject(config models.ProjectConfig, platformData models.PlatformData) error {
	// Base URL for Vercel API
	baseURL := "https://api.vercel.com/v11/projects"

	// Add team parameter to URL if provided
	requestURL := baseURL
	if platformData.TeamId != "" && platformData.TeamId != "team_xxxx" {
		requestURL = fmt.Sprintf("%s?teamId=%s", baseURL, platformData.TeamId)
	}

	// Set outputDirectory based on build folder
	// var outputDir interface{} = nil
	// if config.BuildFolder != "./" {
	// 	outputDir = config.BuildFolder
	// }

	// Updated project data according to v11 API format from provided example
	projectData := map[string]interface{}{
		"name":                              config.Name,
		"buildCommand":                      nil,
		"commandForIgnoringBuildStep":       nil,
		"devCommand":                        nil,
		"environmentVariables":              []map[string]interface{}{},
		"framework":                         platformData.Framework,
		"installCommand":                    nil,
		"outputDirectory":                   nil,
		"publicSource":                      nil,
		"rootDirectory":                     nil,
		"skipGitConnectDuringLink":          true, // Set to true to skip GitHub connection
		"enableAffectedProjectsDeployments": true,
		"oidcTokenConfig": map[string]interface{}{
			"enabled":    true,
			"issuerMode": "global",
		},
	}

	jsonData, err := json.Marshal(projectData)
	if err != nil {
		return fmt.Errorf("failed to marshal project data: %w", err)
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+platformData.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Handle response status codes
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		// Parse error response
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
			if errMsg, ok := errorResponse["error"].(map[string]interface{})["message"]; ok {
				return fmt.Errorf("failed to create project, status code: %d, message: %v", resp.StatusCode, errMsg)
			}

		}
		return fmt.Errorf("failed to create project, status code: %d", resp.StatusCode)
	}

	return nil
}

func DeleteVercelProject(projectName, deployBranch, buildFolder string) error {
	return nil
}

func GetVercelProject(projectName, deployBranch, buildFolder string) error {
	return nil
}
