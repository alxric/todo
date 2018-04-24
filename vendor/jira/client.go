package jira

import (
	"crypto/rsa"

	"github.com/dghubble/oauth1"
)

// Client is the Jira client clients can use to talk to the jira api
type Client struct {
	BaseURL    string
	Token      string
	Secret     string
	Session    string
	OauthCfg   *oauth1.Config
	PrivateKey *rsa.PrivateKey
}

// NewClient will return a client able to perform Jira api calls
func NewClient(config *oauth1.Config) *Client {
	c := &Client{
		OauthCfg: config,
	}
	return c
}
