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
	"fmt"
	"jira"
	"todo/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// toggleCmd represents the toggle command
var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Args:  cobra.ExactArgs(1),
	Short: "Toggles the supplied task between offline/online",
	RunE: func(cmd *cobra.Command, args []string) error {
		return t.Toggle(args)
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)
}

// Toggle a task online/offline
func (t *Todo) Toggle(args []string) error {
	var i int
	var err error
	if i, err = t.verifyIndex(args); err != nil {
		return err
	}
	task := t.Config.Tasks[i]
	switch task.JiraID {
	case "":
		fmt.Printf("Creating Jira issue for '%s'\n", task.Text)
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
		t.createJira(issue, task)
	default:
		fmt.Printf("Taking todo '%s' offline\n", task.Text)
		task.JiraID = ""
		task.JiraKey = ""
	}
	config.Write(viper.ConfigFileUsed(), &t.Config)

	return nil
}
