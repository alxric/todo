package main

import (
	"jira"
	"testing"
	"todo/cmd"
	"todo/internal/config"
)

func TestAddTodo(t *testing.T) {
	todo := cmd.Todo{}
	todo.JC = &jira.Client{BaseURL: "1"}
	todo.Config = &config.Cfg{}
	if err := todo.Add("1", true); err != nil {
		t.Errorf("Could not create issue: %v", err)
	}
	if err := todo.Add("", true); err == nil {
		t.Errorf("Should not be able to create empty issue")
	}
}
