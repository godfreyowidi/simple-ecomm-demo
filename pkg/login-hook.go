package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Auth0LoginResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"` // add this
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope,omitempty"`
}

func CustomerLogin(ctx context.Context, email, password string) (*Auth0LoginResponse, error) {
	domain := os.Getenv("AUTH0_DOMAIN")
	clientID := os.Getenv("AUTH0_LOGIN_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_LOGIN_CLIENT_SECRET")
	audience := os.Getenv("AUTH0_API_AUDIENCE")
	grantType := "password"
	scope := "openid profile email offline_access" // ✅ include offline_access

	loginData := map[string]string{
		"grant_type":    grantType,
		"client_id":     clientID,
		"client_secret": clientSecret,
		"username":      email,
		"password":      password,
		"audience":      audience,
		"scope":         scope,
	}

	body, _ := json.Marshal(loginData)

	resp, err := http.Post(fmt.Sprintf("https://%s/oauth/token", domain), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var msg map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&msg)
		return nil, fmt.Errorf("login failed: %v", msg)
	}

	var loginResp Auth0LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %w", err)
	}

	return &loginResp, nil
}

func RefreshCustomerToken(ctx context.Context, refreshToken string) (*Auth0LoginResponse, error) {
	domain := os.Getenv("AUTH0_DOMAIN")
	clientID := os.Getenv("AUTH0_LOGIN_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_LOGIN_CLIENT_SECRET")

	refreshData := map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
	}

	body, _ := json.Marshal(refreshData)

	resp, err := http.Post(fmt.Sprintf("https://%s/oauth/token", domain), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var msg map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&msg)
		return nil, fmt.Errorf("refresh failed: %v", msg)
	}

	var refreshResp Auth0LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	return &refreshResp, nil
}
