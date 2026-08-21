package folders

import (
	"fmt"
	"strings"
)

func requireProjectID(projectID string) (string, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return "", fmt.Errorf("project id is empty")
	}
	return projectID, nil
}

func requireProjectFolderIDs(projectID, folderID string) (string, string, error) {
	projectID, err := requireProjectID(projectID)
	if err != nil {
		return "", "", err
	}
	folderID = strings.TrimSpace(folderID)
	if folderID == "" {
		return "", "", fmt.Errorf("folder id is empty")
	}
	return projectID, folderID, nil
}
