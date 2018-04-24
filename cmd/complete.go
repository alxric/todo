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

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete",
	Args:  cobra.ExactArgs(1),
	Short: "Completes the task with the index supplied",
	RunE: func(cmd *cobra.Command, args []string) error {
		return complete(args)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

}

func complete(args []string) error {
	var i int
	var err error
	if i, err = verifyIndex(args); err != nil {
		return err
	}
	t := cfg.Tasks[i]
	if t.Done == true {
		return fmt.Errorf("Task '%s' is already completed", t.Text)
	}
	fmt.Printf("Completing task: '%s'\n", t.Text)
	err = j.ChangeIssueStatus(t.JiraID, cfg.Jira.Project.DoneID, "Closed by Todo")
	if err != nil {
		return fmt.Errorf("Unable to change Jira status: %v", err)
	}
	t.Done = true
	t.Completed = time.Now()
	config.Write(viper.ConfigFileUsed(), &cfg)
	return nil
}
