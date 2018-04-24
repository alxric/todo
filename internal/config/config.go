package config

import (
	"io/ioutil"
	"os"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Cfg contains the whole yaml config
type Cfg struct {
	Tasks []*Task `yaml:"tasks"`
	Jira  Jira    `yaml:"jira"`
}

// Jira contains the Jira config
type Jira struct {
	URL     string  `yaml:"url"`
	Token   string  `yaml:"token"`
	Session string  `yaml:"session"`
	Project Project `yaml:"project"`
}

// Project contains the active project the user has chosen
type Project struct {
	Name      string `yaml:"name"`
	DoneID    string `yaml:"done_id"`
	BacklogID string `yaml:"backlog_id"`
	ID        string `yaml:"id"`
	Key       string `yaml:"key"`
}

// Task defines a todo task
type Task struct {
	Text      string    `yaml:"text"`
	JiraID    string    `yaml:"jira_id"`
	JiraKey   string    `yaml:"jira_key"`
	Done      bool      `yaml:"done"`
	Created   time.Time `yaml:"created"`
	Completed time.Time `yaml:"completed"`
}

// Read supplied yaml file and parse to Cfg
func Read(path string) (*Cfg, error) {
	var cfg Cfg
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Write supplied cfg to the yaml file
func Write(path string, cfg interface{}) error {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	dPath := strings.Split(path, "/")
	if len(dPath) > 0 {
		err = os.MkdirAll(strings.Join(dPath[0:len(dPath)-1], "/"), 0755)
		if err != nil {
			return err
		}

	}
	err = ioutil.WriteFile(path, b, 0666)
	if err != nil {
		return err
	}
	return nil
}
