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
		return add(strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

type jiraReply struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func add(text string) error {
	if j.BaseURL == "" {
		return fmt.Errorf("Could not read configuration. Run 'todo init'")
	}
	t := &config.Task{
		Text:    text,
		Done:    false,
		Created: time.Now(),
	}
	issue := &jira.Issue{
		Fields: jira.Fields{
			Summary: fmt.Sprintf("TODO: %s", t.Text),
			Project: jira.IssueProject{
				ID: cfg.Jira.Project.ID,
			},
			IssueType: jira.IssueType{
				ID: "10002",
			},
		},
	}

	b, err := j.CreateIssue(issue)
	if err != nil {
		return fmt.Errorf("Unable to create Jira issue; %v", err)
	}
	jr := &jiraReply{}
	err = json.Unmarshal(b, jr)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal Jira reply: %v", err)
	}
	t.JiraID = jr.ID
	t.JiraKey = jr.Key
	cfg.Tasks = append(cfg.Tasks, t)
	config.Write(viper.ConfigFileUsed(), &cfg)
	fmt.Println(fmt.Sprintf("Task added: '%s'", t.Text))
	return nil
}
