package lr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	LrURL        = "https://devapi.lrinternal.com"
	tokenURL     = "https://demotesting.devhub.lrinternal.com/service/oauth/token"
	clientID     = ""
	clientSecret = ""
	audience     = "https://devapi.lrinternal.com/identity/v2/manage"
	CustomerID   = "7c9254057e2044c5b3fadf8bf0b3dd31"
	AppId        = "99207378"
)

// Token storage
var (
	m2mToken    string
	tokenExpiry int64
	mu          sync.Mutex
)

func GetM2MToken() string {
	mu.Lock()
	defer mu.Unlock()

	// Check if token is still valid
	if time.Now().Unix() < tokenExpiry-10 {
		return m2mToken
	}

	// Generate new token
	token, expiresIn, err := generateM2MToken()
	if err != nil {
		return ""
	}

	// Store the new token globally
	m2mToken = token
	tokenExpiry = time.Now().Unix() + expiresIn

	return m2mToken
}

// Calls the API to generate a new M2M token
func generateM2MToken() (string, int64, error) {
	payload := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"grant_type":    "client_credentials",
		"audience":      audience,
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(data))
	if err != nil {
		return "", 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", 0, fmt.Errorf("failed to get token: %s", string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", 0, err
	}

	return tokenResp.AccessToken, tokenResp.ExpiresIn, nil
}
