package jira

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
)

func (c *Client) apiCall(path string, method string, data io.Reader) ([]byte, int, error) {
	apiURL, err := c.generateAPIURL(path, method)
	if err != nil {
		return nil, 500, err
	}
	req, err := http.NewRequest(strings.ToUpper(method), apiURL, data)
	if err != nil {
		return nil, 500, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 500, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, err
	}
	return body, resp.StatusCode, nil
}

func (c *Client) generateAPIURL(path string, method string) (string, error) {
	var apiURL, signString strings.Builder
	ts := strconv.Itoa(int(time.Now().Unix()))
	hash := "2jmj7l5rSw0yVb/vlWAYkK/YBwk="
	nonce := generateNonce(8)
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	apiURL.WriteString(parsedURL.Scheme)
	apiURL.WriteString("://")
	apiURL.WriteString(parsedURL.Hostname())
	if parsedURL.Port() != "" {
		apiURL.WriteString(":")
		apiURL.WriteString(parsedURL.Port())
	}
	apiURL.WriteString(parsedURL.EscapedPath())

	apiURL.WriteString("?oauth_body_hash=")
	apiURL.WriteString(oauth1.PercentEncode(hash))
	apiURL.WriteString("&oauth_nonce=")
	apiURL.WriteString(nonce)
	apiURL.WriteString("&oauth_timestamp=")
	apiURL.WriteString(ts)
	apiURL.WriteString("&oauth_consumer_key=")
	apiURL.WriteString(c.OauthCfg.ConsumerKey)
	apiURL.WriteString("&oauth_signature_method=RSA-SHA1")
	apiURL.WriteString("&oauth_version=1.0")
	apiURL.WriteString("&oauth_token=")
	apiURL.WriteString(c.Token)
	apiURL.WriteString("&oauth_signature=")

	signString.WriteString(strings.ToUpper(method))
	signString.WriteString("&")
	signString.WriteString(parsedURL.Scheme)
	signString.WriteString("://")
	signString.WriteString(parsedURL.Hostname())
	if parsedURL.Port() != "" {
		signString.WriteString(":")
		signString.WriteString(parsedURL.Port())
	}
	signString.WriteString(parsedURL.EscapedPath())
	signString.WriteString("&oauth_body_hash=")
	signString.WriteString("2jmj7l5rSw0yVb%252FvlWAYkK%252FYBwk%253D")
	signString.WriteString("&oauth_consumer_key=")
	signString.WriteString(c.OauthCfg.ConsumerKey)
	signString.WriteString("&oauth_nonce=")
	signString.WriteString(nonce)
	signString.WriteString("&oauth_signature_method=RSA-SHA1&oauth_timestamp=")
	signString.WriteString(ts)
	signString.WriteString("&oauth_token=")
	signString.WriteString(c.Token)
	signString.WriteString("&oauth_version=1.0")
	for key, val := range parsedURL.Query() {
		signString.WriteString("&")
		signString.WriteString(key)
		signString.WriteString("=")
		signString.WriteString(strings.Join(val, ""))
	}
	ss := strings.Replace(signString.String(), "&", "%26", -1)
	ss = strings.Replace(ss, "%26", "&", 2)
	ss = strings.Replace(ss, ":", "%3A", -1)
	ss = strings.Replace(ss, "/", "%2F", -1)
	ss = strings.Replace(ss, "=", "%3D", -1)
	signedString, err := c.OauthCfg.Signer.Sign("", ss)
	if err != nil {
		return "", err
	}
	signedString = oauth1.PercentEncode(signedString)
	apiURL.WriteString(signedString)
	for key, val := range parsedURL.Query() {
		apiURL.WriteString("&")
		apiURL.WriteString(key)
		apiURL.WriteString("=")
		apiURL.WriteString(strings.Join(val, ""))
	}
	return apiURL.String(), nil
}
