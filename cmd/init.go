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
	"bufio"
	"fmt"
	"jira"
	"os"
	"strconv"
	"strings"
	"todo/internal/config"

	"github.com/dghubble/oauth1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalizes a new setup of todo",
	Long: `init will guide you through the setup to get your todo tool to work.

You will be asked to supply the URL to the JIRA installation you want to work
against. Your credentials will not be stored anywhere. Only thing getting stored
is an OAuth token which will then be used for any future authentications against
Jira`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initTodo()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

var jiraURL, jiraUser, jiraPassword string

func initTodo() error {
	err := setupJira()
	if err != nil {
		return err
	}
	return nil
}

func setupJira() error {
	j = jira.NewClient(&oauth1.Config{
		CallbackURL:    "oob",
		ConsumerKey:    "Todo",
		ConsumerSecret: "dont_care",
		Endpoint:       oauth1.Endpoint{},
		Signer: &oauth1.RSASigner{
			PrivateKey: privateKey,
		},
	})
	j.PrivateKey = privateKey
	fmt.Printf("Initalizing todo...\n\n")
	err := j.GenerateOauthToken()
	if err != nil {
		return err
	}
	fmt.Printf("\nGreat! Authentication succesful! Let's proceed...\n")

	p, err := chooseProject()
	if err != nil {
		return fmt.Errorf("Unable to choose Jira project: %v", err)
	}

	s, backlog, err := chooseDoneStatus(p.Key)
	if err != nil {
		return fmt.Errorf("Unable to choose Jira project: %v", err)
	}
	token, err := config.Encrypt(&privateKey.PublicKey, j.Token)
	if err != nil {
		return err
	}
	session, err := config.Encrypt(&privateKey.PublicKey, j.Session)
	if err != nil {
		return err
	}
	cfg.Jira.URL = j.BaseURL
	cfg.Jira.Token = token
	cfg.Jira.Session = session
	cfg.Jira.Project = config.Project{
		Name:      p.Name,
		DoneID:    s.ID,
		BacklogID: backlog.ID,
		ID:        p.ID,
		Key:       p.Key,
	}
	err = config.Write(viper.ConfigFileUsed(), cfg)
	if err != nil {
		return err
	}
	fmt.Printf("\nInitialization done! Enjoy the tool\n")
	return nil
}

func chooseProject() (*jira.Project, error) {
	var p jira.Project
	for {
		fmt.Printf("\nChoose which project you want Todo to use:\n\n")
		projects, err := j.ListProjects()
		if err != nil {
			return nil, err
		}
		for index, project := range projects {
			fmt.Println(fmt.Sprintf("%v) %s", index+1, project.Name))
		}
		fmt.Printf("\nPick a project: ")
		reader := bufio.NewReader(os.Stdin)
		projectChoice, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("Could not read input from picking a project ")
		}
		pInt, err := strconv.Atoi(strings.TrimSpace(projectChoice))
		if err != nil {
			fmt.Printf("\nInvalid project choice! Only numbers between 1 and %v is available :%v\n",
				len(projects), err)
			continue
		}
		if pInt <= 0 || pInt > len(projects) {
			fmt.Printf("\nInvalid project choice! Only numbers between 1 and %v is available\n",
				len(projects))
			continue
		}
		fmt.Printf("\nProject '%s' chosen.\n",
			projects[pInt-1].Name)
		p = projects[pInt-1]
		break
	}
	return &p, nil
}

func chooseDoneStatus(projectKey string) (*jira.Transition, *jira.Transition, error) {
	var p, bl jira.Transition
	for {
		fmt.Printf("\nChoose which status you want to use when marking an issue as completed:\n\n")
		searchJSON := []byte(fmt.Sprintf(
			`{"jql":"project = %s", "startAt": 0, "maxResults": 1}`, projectKey))
		si, err := j.SearchIssues(searchJSON)
		if err != nil || len(si) == 0 {
			return nil, nil, err
		}
		transitions, err := j.ListTransitions(si[0].ID)
		if err != nil {
			return nil, nil, err
		}
		for index, transition := range transitions {
			if transition.Name == "Backlog" {
				bl = transition
			}
			fmt.Println(fmt.Sprintf("%v) %s", index+1, transition.Name))
		}
		fmt.Printf("\nPick a status: ")
		reader := bufio.NewReader(os.Stdin)
		projectChoice, err := reader.ReadString('\n')
		if err != nil {
			return nil, nil, fmt.Errorf("Could not read input from picking a project status")
		}
		pInt, err := strconv.Atoi(strings.TrimSpace(projectChoice))
		if err != nil {
			fmt.Printf("\nInvalid status choice! Only numbers between 1 and %v is available :%v\n",
				len(transitions), err)
			continue
		}
		if pInt <= 0 || pInt > len(transitions) {
			fmt.Printf("\nInvalid status choice! Only numbers between 1 and %v is available\n",
				len(transitions))
			continue
		}
		fmt.Printf("\nStatus '%s' chosen.\n", transitions[pInt-1].Name)
		p = transitions[pInt-1]
		break
	}
	return &p, &bl, nil
}
