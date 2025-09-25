package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"justTest/internal/models"
	"net/http"
	"time"
)

type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
func (c *AuthClient) GetUserByID(userID string) (*models.AuthUser, error) {
	if userID == "" {
		return nil, fmt.Errorf("invalid user id")
	}
	url := fmt.Sprintf("%s/api/auth/user/%s", c.baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %s", err)

	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)

	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error executing request: %s %s", resp.Status, string(body))

	}
	var user models.AuthUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %s", err)
	}
	return &user, nil
}

func (c *AuthClient) ValidateToken(token string) (*models.AuthUser, error) {
	if token == "" {
		return nil, fmt.Errorf("invalid token")

	}
	url := fmt.Sprintf("%s/api/auth/validate%s", c.baseURL, token)

	requestData := map[string]string{"token": token}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %s", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error executing request: %s %s", resp.Status, string(body))
	}
	var user models.AuthUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %s", err)
	}
	return &user, nil
}
