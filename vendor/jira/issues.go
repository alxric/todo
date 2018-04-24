package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type searchIssue struct {
	Issues []Issue
}

// Issue describes a Jira issue
type Issue struct {
	Fields Fields `json:"fields"`
	ID     string `json:"id"`
}

// Fields describes Issue-> Fields
type Fields struct {
	Summary   string       `json:"summary"`
	Project   IssueProject `json:"project"`
	IssueType IssueType    `json:"issuetype"`
}

// IssueProject describes Issue->Fields->Project
type IssueProject struct {
	ID string `json:"id"`
}

// IssueType describes Issue->Fields->IssueType
type IssueType struct {
	ID string `json:"id"`
}

// CreateIssue from supplied Issue
func (c *Client) CreateIssue(issue *Issue) ([]byte, error) {
	j, err := json.Marshal(issue)
	if err != nil {
		return nil, err
	}
	b, status, err := c.apiCall(
		fmt.Sprintf("%s/rest/api/2/issue", c.BaseURL),
		"POST",
		bytes.NewReader(j))
	if err != nil {
		return nil, err
	}
	if status != 201 {
		return nil, fmt.Errorf("Unable to create issue: %s", string(b))
	}
	return b, nil
}

// DeleteIssue  using supplied issueID
func (c *Client) DeleteIssue(issueID string) error {
	b, status, err := c.apiCall(
		fmt.Sprintf("%s/rest/api/2/issue/%s", c.BaseURL, issueID),
		"DELETE",
		nil)
	if err != nil {
		return err
	}
	if status != 204 {
		return fmt.Errorf("Could not delete Jira issue: %s", string(b))
	}
	return nil
}

// ChangeIssueStatus will change status of the supplied issue ID
func (c *Client) ChangeIssueStatus(issueID string, statusID string, msg string) error {
	jsonB := []byte(fmt.Sprintf("{\"transition\":{\"id\":\"%s\"}}", statusID))
	b, status, err := c.apiCall(
		fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.BaseURL, issueID),
		"POST",
		bytes.NewReader(jsonB))
	if err != nil {
		return err
	}
	if status != 204 {
		return fmt.Errorf("Could not change Jira issue status: %s", string(b))
	}
	return nil
}

// SearchIssues given the supplied search string
func (c *Client) SearchIssues(searchJSON []byte) ([]Issue, error) {
	b, status, err := c.apiCall(
		fmt.Sprintf("%s/rest/api/2/search", c.BaseURL),
		"POST",
		bytes.NewReader(searchJSON))
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("Could not change search Jira issues: %s", string(b))
	}
	si := &searchIssue{}
	err = json.Unmarshal(b, &si)
	if err != nil {
		return nil, err
	}
	return si.Issues, nil
}

// IssueTransitions describes possible issue transitions
type IssueTransitions struct {
	Expand      string       `json:"expand"`
	Transitions []Transition `json:"transitions"`
}

// Transition describes an issue transition
type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	To   struct {
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
	} `json:"to"`
}

// ListTransitions will show all available transitions for a supplied issue ID
func (c *Client) ListTransitions(issueID string) ([]Transition, error) {
	b, status, err := c.apiCall(
		fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.BaseURL, issueID),
		"GET",
		nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("Could not list Jira transitions: %s", string(b))
	}
	it := &IssueTransitions{}
	err = json.Unmarshal(b, it)
	if err != nil {
		return nil, err
	}
	return it.Transitions, nil
}
