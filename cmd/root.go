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
	"crypto/rsa"
	"fmt"
	"jira"
	"os"
	"strconv"
	"strings"
	"todo/internal/config"

	"github.com/dghubble/oauth1"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	privKeyFile string
	privateKey  *rsa.PrivateKey
	cfg         *config.Cfg
	token       string
	j           *jira.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Small CLI todo app that syncs with JIRA",
	Long: `Todo can be used to manage your todo list from the command line.

It will also create your issues in a JIRA project defined during initialization`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := viper.ReadInConfig(); err != nil {
		if len(os.Args) >= 2 {
			switch os.Args[1] {
			case "init":
			default:
				fmt.Println("Configuration not found! Run 'todo init'")
				os.Exit(1)
			}
		}

	}
	if err := rootCmd.Execute(); err != nil {
		switch {
		case strings.Contains(err.Error(), "token_rejected"):
			fmt.Println("Token has expired. Run 'todo init' again.")
		default:
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

func init() {
	cfg = &config.Cfg{}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.config/todo/todo.yaml)",
	)
	rootCmd.PersistentFlags().StringVar(
		&privKeyFile,
		"privkey",
		"",
		"private RSA key to authenticate against Jira (default is $HOME/.ssh/jira_privatekey.pem)",
	)
	initConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(fmt.Sprintf("%s/.config/todo", home))
		viper.SetConfigName("todo")
	}

	if privKeyFile == "" {
		privKeyFile = fmt.Sprintf("%s/.ssh/jira_privatekey.pem", home)
	}
	privateKey, err = config.ReadRSAKey(privKeyFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.ReadInConfig(); err == nil {
		viper.SetConfigType("yaml")
		cfg, err = config.Read(viper.ConfigFileUsed())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if cfg.Jira.Token != "" {
			token, err = config.Decrypt(privateKey, cfg.Jira.Token)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			initJiraClient()
		}
	} else {
		viper.SetConfigFile(fmt.Sprintf("%s/.config/todo/todo.yaml", home))
	}

}

func initJiraClient() {
	j = jira.NewClient(&oauth1.Config{
		CallbackURL:    "oob",
		ConsumerKey:    "Todo",
		ConsumerSecret: "dont_care",
		Signer: &oauth1.RSASigner{
			PrivateKey: privateKey,
		},
	})
	j.PrivateKey = privateKey
	j.BaseURL = cfg.Jira.URL
	j.Token = token
}

func verifyIndex(args []string) (int, error) {
	i, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid index specified: %s", args[0])
	}
	switch {
	case i <= 0:
		return 0, fmt.Errorf("Invalid index specified: %d", i)
	case i > len(cfg.Tasks):
		return 0, fmt.Errorf("Invalid index specified: %d", i)
	}
	return i - 1, nil
}
