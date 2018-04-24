package jira

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func verifyURL(jiraURL string) (string, error) {
	jiraURL = strings.TrimSpace(strings.ToLower(jiraURL))
	var err error
	parsedURL, err := url.ParseRequestURI(jiraURL)
	if err != nil {
		return "", err
	}
	return parsedURL.String(), nil

}

func generateNonce(length int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	var nonce strings.Builder
	for i := 0; i < length; i++ {
		d := r1.Intn(9)
		nonce.WriteString(strconv.Itoa(d))
	}
	return nonce.String()
}
