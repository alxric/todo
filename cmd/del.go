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
	"todo/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Args:  cobra.ExactArgs(1),
	Short: "Deletes the task with the index supplied from your todo list",
	RunE: func(cmd *cobra.Command, args []string) error {
		return del(args)
	},
}

func init() {
	rootCmd.AddCommand(delCmd)
}

func del(args []string) error {
	var i int
	var err error
	if i, err = verifyIndex(args); err != nil {
		return err
	}
	t := cfg.Tasks[i]
	fmt.Printf("Deleting task: '%s'\n", t.Text)
	err = j.DeleteIssue(t.JiraID)
	if err != nil {
		return fmt.Errorf("Unable to delete Jira issue: %v", err)
	}
	cfg.Tasks = append(cfg.Tasks[:i-1], cfg.Tasks[i:]...)
	config.Write(viper.ConfigFileUsed(), &cfg)

	return nil
}
