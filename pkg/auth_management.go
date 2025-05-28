package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Auth0TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func GetManagementToken(ctx context.Context) (string, error) {
	domain := os.Getenv("AUTH0_DOMAIN")
	clientID := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	audience := os.Getenv("AUTH0_MANAGEMENT_API_AUDIENCE")

	fmt.Println("[Auth0] Requesting management token with:")
	fmt.Println("  Domain   :", domain)
	fmt.Println("  Client ID:", clientID)
	fmt.Println("  Audience :", audience)

	if domain == "" || clientID == "" || clientSecret == "" || audience == "" {
		return "", errors.New("missing Auth0 configuration environment variables")
	}

	data := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     clientID,
		"client_secret": clientSecret,
		"audience":      audience,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token request: %w", err)
	}

	url := fmt.Sprintf("https://%s/oauth/token", domain)
	fmt.Println("[Auth0] POST", url)

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("HTTP error getting token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var msg map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&msg)
		fmt.Printf("[Auth0] Token request failed (%d): %+v\n", resp.StatusCode, msg)
		return "", fmt.Errorf("failed to get management token: %v", msg)
	}

	var tokenResp Auth0TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	fmt.Println("[Auth0] Successfully obtained token")
	fmt.Println("  Token Type :", tokenResp.TokenType)
	fmt.Println("  Token (truncated) :", tokenResp.AccessToken[:20]+"...")

	return tokenResp.AccessToken, nil
}

type CreateUserPayload struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Connection string `json:"connection"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

type Auth0UserResponse struct {
	UserID string `json:"user_id"`
}

func CreateUser(ctx context.Context, token string, payload CreateUserPayload) (string, error) {
	domain := os.Getenv("AUTH0_DOMAIN")
	url := fmt.Sprintf("https://%s/api/v2/users", domain)

	fmt.Println("[Auth0] Creating user at:", url)
	fmt.Printf("[Auth0] Payload: %+v\n", payload)

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("[Auth0] Authorization Header: Bearer", token[:20]+"...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP error creating user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var msg map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&msg)
		fmt.Printf("[Auth0] User creation failed (%d): %+v\n", resp.StatusCode, msg)
		return "", fmt.Errorf("auth0 user creation failed: %v", msg)
	}

	var user Auth0UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", fmt.Errorf("failed to decode user creation response: %w", err)
	}

	fmt.Println("[Auth0] User created successfully")
	fmt.Println("  User ID:", user.UserID)

	return user.UserID, nil
}
