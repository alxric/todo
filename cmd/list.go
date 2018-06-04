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
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:       "list",
	Short:     "Lists all of your tasks",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"open", "completed"},
	RunE: func(cmd *cobra.Command, args []string) error {
		filter, err := cmd.Flags().GetString("filter")

		if err != nil {
			return fmt.Errorf("Invalid filter. Valid options are: [open, completed]")
		}
		filter = strings.TrimSpace(strings.ToLower(filter))
		if filter != "open" && filter != "completed" && filter != "" {
			return fmt.Errorf("Invalid filter. Valid options are: [open, completed]")
		}
		return t.List(filter)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("filter", "f", "", "Filter list [open, completed]")

}

// List tasks
func (t *Todo) List(filter string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Done", "Task", "Created", "Completed", "URL"})
	table.SetBorder(true)
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	for index, task := range t.Config.Tasks {
		switch {
		case filter == "open" && task.Done:
			continue
		case filter == "completed" && !task.Done:
			continue
		}
		var done, completed string
		taskName := task.Text
		if len(taskName) >= 20 {
			taskName = fmt.Sprintf("%s...", taskName[:16])
		}
		if task.Done {
			done = " X "
		}
		if task.Completed.Unix() > 0 {
			completed = task.Completed.Format("2006-01-02 15:04:05")
		}
		var jiraURL string
		if task.JiraKey != "" {
			jiraURL = fmt.Sprintf("%s/browse/%s", t.Config.Jira.URL, task.JiraKey)
		}
		row := []string{

			fmt.Sprintf("%d", index+1),
			done,
			taskName,
			task.Created.Format("2006-01-02 15:04:05"),
			completed,
			jiraURL,
		}
		table.Append(row)
	}
	table.Render()
	return nil
}
