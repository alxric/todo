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
	"time"
	"todo/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// oopsCmd represents the oops command
var oopsCmd = &cobra.Command{
	Use:   "oops",
	Args:  cobra.ExactArgs(1),
	Short: "Re-opens an accidentally completed task",
	RunE: func(cmd *cobra.Command, args []string) error {
		return t.Oops(args)
	},
}

func init() {
	rootCmd.AddCommand(oopsCmd)
}

// Oops will re-open a closed issue
func (t *Todo) Oops(args []string) error {
	var i int
	var err error
	if i, err = t.verifyIndex(args); err != nil {
		return err
	}
	task := t.Config.Tasks[i]
	if task.Done == false {
		return fmt.Errorf("Task '%s' is not completed", task.Text)
	}
	fmt.Printf("Re-opening task: '%s'\n", task.Text)
	if task.JiraID != "" {
		err = t.JC.ChangeIssueStatus(task.JiraID, t.Config.Jira.Project.BacklogID, "Re-opened by Todo")
		if err != nil {
			return fmt.Errorf("Unable to change Jira status: %v", err)
		}
	}
	task.Done = false
	var tt time.Time
	task.Completed = tt
	config.Write(viper.ConfigFileUsed(), &t.Config)
	return nil
}
