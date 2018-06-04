// Copyright © 2018 Alexander Rickardsson <alex@rickardsson.se>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"jira"
	"strings"
	"time"
	"todo/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new task",
	RunE: func(cmd *cobra.Command, args []string) error {
		offline, err := cmd.Flags().GetBool("offline")
		if err != nil {
			return fmt.Errorf("Invalid value for offline")
		}
		return t.Add(strings.Join(args, " "), offline)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolP("offline", "o", false, "Offline mode (Don't create JIRA task)")
}

type jiraReply struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

// Add a new task
func (t *Todo) Add(text string, offline bool) error {
	if text == "" {
		return fmt.Errorf("Your todo can't be empty")
	}
	if t.JC.BaseURL == "" {
		return fmt.Errorf("Could not read configuration. Run 'todo init'")
	}
	task := &config.Task{
		Text:    text,
		Done:    false,
		Created: time.Now(),
	}
	if !offline {
		issue := &jira.Issue{
			Fields: jira.Fields{
				Summary: fmt.Sprintf("TODO: %s", task.Text),
				Project: jira.IssueProject{
					ID: t.Config.Jira.Project.ID,
				},
				IssueType: jira.IssueType{
					ID: t.Config.Jira.Project.IssueType,
				},
			},
		}
		err := t.createJira(issue, task)
		if err != nil {
			return fmt.Errorf("Unable to create Jira issue: %v", err)
		}
	}

	t.Config.Tasks = append(t.Config.Tasks, task)
	config.Write(viper.ConfigFileUsed(), &t.Config)
	fmt.Println(fmt.Sprintf("Task added: '%s'", task.Text))
	return nil
}

func (t *Todo) createJira(issue *jira.Issue, task *config.Task) error {
	b, err := t.JC.CreateIssue(issue)
	if err != nil {
		return fmt.Errorf("Unable to create Jira issue; %v", err)
	}
	jr := &jiraReply{}
	err = json.Unmarshal(b, jr)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal Jira reply: %v", err)
	}
	task.JiraID = jr.ID
	task.JiraKey = jr.Key
	return nil
}
