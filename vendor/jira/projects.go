package jira

import (
	"encoding/json"
	"fmt"
)

// Project describes a Jira project
type Project struct {
	Expand     string `json:"expand"`
	Self       string `json:"self"`
	ID         string `json:"id"`
	Key        string `json:"key"`
	Name       string `json:"name"`
	AvatarUrls struct {
		Four8X48  string `json:"48x48"`
		Two4X24   string `json:"24x24"`
		One6X16   string `json:"16x16"`
		Three2X32 string `json:"32x32"`
	} `json:"avatarUrls"`
	ProjectTypeKey  string `json:"projectTypeKey"`
	ProjectCategory struct {
		Self        string `json:"self"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"projectCategory,omitempty"`
}

// ListProjects that the oauth token has access to
func (c *Client) ListProjects() ([]Project, error) {
	b, _, err := c.apiCall(
		fmt.Sprintf("%s/rest/api/2/project?recent=10", c.BaseURL),
		"GET",
		nil)
	if err != nil {
		return nil, err
	}
	projects := []Project{}
	err = json.Unmarshal(b, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// ProjectStatuses describes possible statuses for a project
type ProjectStatuses struct {
	Self     string   `json:"self"`
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Subtask  bool     `json:"subtask"`
	Statuses []Status `json:"statuses"`
}

// Status describes a Jira status
type Status struct {
	Self           string `json:"self"`
	Description    string `json:"description"`
	IconURL        string `json:"iconUrl"`
	Name           string `json:"name"`
	ID             string `json:"id"`
	StatusCategory struct {
		Self      string `json:"self"`
		ID        int    `json:"id"`
		Key       string `json:"key"`
		ColorName string `json:"colorName"`
		Name      string `json:"name"`
	} `json:"statusCategory"`
}

// ListProjectStatuses will return all possible statuses for a given project ID
func (c *Client) ListProjectStatuses(projectID string) ([]ProjectStatuses, error) {
	b, _, err := c.apiCall(
		fmt.Sprintf("%s/rest/api/2/project/%s/statuses", c.BaseURL, projectID),
		"GET",
		nil)
	if err != nil {
		return nil, err
	}
	statuses := []ProjectStatuses{}
	err = json.Unmarshal(b, &statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
