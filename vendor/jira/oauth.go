package jira

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
)

// GenerateOauthToken will go through the oauth dance to deliver an oaut token
func (c *Client) GenerateOauthToken() error {
	// Get Jira URL from user
	jiraURL, err := readURL()
	if err != nil {
		return fmt.Errorf("Could not generate Jira URL: %v", err)
	}
	c.BaseURL = jiraURL
	c.OauthCfg.Endpoint.RequestTokenURL = fmt.Sprintf(
		"%s/plugins/servlet/oauth/request-token", c.BaseURL)
	c.OauthCfg.Endpoint.AccessTokenURL = fmt.Sprintf(
		"%s/plugins/servlet/oauth/access-token", c.BaseURL)
	c.OauthCfg.Endpoint.AuthorizeURL = fmt.Sprintf(
		"%s/plugins/servlet/oauth/authorize", c.BaseURL)
	requestToken, _, err := c.requestOauthToken()
	if err != nil {
		return fmt.Errorf("Could not request oauth token: %v", err)
	}
	err = userAuth(c.OauthCfg, requestToken)
	if err != nil {
		return err
	}

	err = c.AccessToken(requestToken)
	if err != nil {
		return err
	}

	return nil
}

// AccessToken will set a new oauth token for the client
func (c *Client) AccessToken(requestToken string) error {
	accessURL, err := c.generateAccessURL(requestToken)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", accessURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "'application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return err
	}
	c.Token = values.Get(oauthTokenParam)
	c.Secret = values.Get(oauthTokenSecretParam)
	if c.Token == "" || c.Secret == "" {
		return errors.New(`Response missing oauth_token or oauth_token_secret
Did you remember to click 'Allow' on the website?`)
	}
	return nil
}

func userAuth(cfg *oauth1.Config, requestToken string) error {
	fmt.Printf("\nVisit %s?oauth_token=%s and give the tool access\n\n",
		cfg.Endpoint.AuthorizeURL, requestToken)
	for {
		fmt.Printf("Are you done? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		done, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("Unable to read input to are you done?: %v", err)
		}
		if strings.TrimSpace(strings.ToLower(done)) == "y" {
			break
		}
	}
	return nil
}

func readURL() (string, error) {
	fmt.Printf("Enter Jira URL: ")
	reader := bufio.NewReader(os.Stdin)
	jiraURL, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("Unable to read Jira URL: %v", err)
	}
	if jiraURL, err = verifyURL(jiraURL); err != nil {
		return "", fmt.Errorf("Invalid Jira URL: %v", err)
	}
	return jiraURL, nil
}

const (
	oauthTokenSecretParam       = "oauth_token_secret"
	oauthCallbackConfirmedParam = "oauth_callback_confirmed"
	oauthTokenParam             = "oauth_token"
	oauthSessionParam           = "oauth_session_handle"
)

func (c *Client) requestOauthToken() (requestToken, requestSecret string, err error) {
	requestURL, err := c.generateRequestURL()
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Content-Type", "'application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", "", err
	}
	requestToken = values.Get(oauthTokenParam)
	requestSecret = values.Get(oauthTokenSecretParam)

	if requestToken == "" || requestSecret == "" {
		return "", "", errors.New("oauth1: Response missing oauth_token or oauth_token_secret")
	}
	return requestToken, requestSecret, nil
}

func (c *Client) generateAccessURL(requestToken string) (string, error) {
	var accessURL, signString strings.Builder
	nonce := generateNonce(8)
	ts := strconv.Itoa(int(time.Now().Unix()))
	accessURL.WriteString(c.OauthCfg.Endpoint.AccessTokenURL)
	accessURL.WriteString("?oauth_consumer_key=")
	accessURL.WriteString(c.OauthCfg.ConsumerKey)
	accessURL.WriteString("&oauth_nonce=")
	accessURL.WriteString(nonce)
	accessURL.WriteString("&oauth_signature=")

	signString.WriteString("POST&")
	signString.WriteString(c.OauthCfg.Endpoint.AccessTokenURL)
	signString.WriteString("&oauth_consumer_key=")
	signString.WriteString(c.OauthCfg.ConsumerKey)
	signString.WriteString("&oauth_nonce=")
	signString.WriteString(nonce)
	signString.WriteString("&oauth_signature_method=RSA-SHA1&oauth_timestamp=")
	signString.WriteString(ts)
	signString.WriteString("&oauth_token=")
	signString.WriteString(requestToken)
	signString.WriteString("&oauth_version=1.0")
	ss := oauth1.PercentEncode(signString.String())
	ss = strings.Replace(ss, "%26", "&", 2)

	signedString, err := c.OauthCfg.Signer.Sign("", ss)
	if err != nil {
		return "", err
	}
	signedString = oauth1.PercentEncode(signedString)
	accessURL.WriteString(signedString)
	accessURL.WriteString("&oauth_signature_method=RSA-SHA1")
	accessURL.WriteString("&oauth_timestamp=")
	accessURL.WriteString(ts)
	accessURL.WriteString("&oauth_token=")
	accessURL.WriteString(requestToken)
	accessURL.WriteString("&oauth_version=1.0")
	return accessURL.String(), nil
	return "", nil
}

func (c *Client) generateRequestURL() (string, error) {
	var requestURL, signString strings.Builder
	nonce := generateNonce(8)
	ts := strconv.Itoa(int(time.Now().Unix()))
	requestURL.WriteString(c.OauthCfg.Endpoint.RequestTokenURL)
	requestURL.WriteString("?oauth_consumer_key=")
	requestURL.WriteString(c.OauthCfg.ConsumerKey)
	requestURL.WriteString("&oauth_nonce=")
	requestURL.WriteString(nonce)
	requestURL.WriteString("&oauth_signature=")

	signString.WriteString("POST&")
	signString.WriteString(c.OauthCfg.Endpoint.RequestTokenURL)
	signString.WriteString("&oauth_consumer_key=")
	signString.WriteString(c.OauthCfg.ConsumerKey)
	signString.WriteString("&oauth_nonce=")
	signString.WriteString(nonce)
	signString.WriteString("&oauth_signature_method=RSA-SHA1&oauth_timestamp=")
	signString.WriteString(ts)
	signString.WriteString("&oauth_version=1.0")
	ss := oauth1.PercentEncode(signString.String())
	ss = strings.Replace(ss, "%26", "&", 2)

	signedString, err := c.OauthCfg.Signer.Sign("", ss)
	if err != nil {
		return "", err
	}
	signedString = oauth1.PercentEncode(signedString)
	requestURL.WriteString(signedString)
	requestURL.WriteString("&oauth_signature_method=RSA-SHA1")
	requestURL.WriteString("&oauth_timestamp=")
	requestURL.WriteString(ts)
	requestURL.WriteString("&oauth_version=1.0")
	return requestURL.String(), nil

}
