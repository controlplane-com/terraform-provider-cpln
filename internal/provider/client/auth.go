package cpln

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imroc/req/v3"
)

// UnixTime represents a Unix timestamp.
type UnixTime int64

// MinRemaining defines the minimum remaining time before token refresh.
const MinRemaining UnixTime = 10 * 60

// MakeAuthorizationHeader determines whether the current access token is valid and refreshes it if necessary.
func (c *Client) MakeAuthorizationHeader() error {
	// Verify that a refresh token has been provided
	if c.RefreshToken == "" {
		// Return an error if the refresh token is missing
		return errors.New("empty refresh token")
	}

	// Parse the JWT token to inspect its claims
	token, err := jwt.Parse(c.Token, func(token *jwt.Token) (interface{}, error) {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	})

	// Handle parsing errors by initiating a token refresh
	if err != nil {
		// Attempt to update the access token using the refresh token
		err = c.updateAccessToken()

		// Propagate any errors from the token update process
		if err != nil {
			return err
		}

		// Exit early after successfully refreshing the token
		return nil
	}

	// Extract claims from the parsed token object
	claims, ok := token.Claims.(jwt.MapClaims)

	// Return an error if the token claims are not in the expected format
	if !ok {
		return errors.New("invalid token claims")
	}

	// Retrieve the expiration time (in Unix seconds) from the claims
	exp, ok := claims["exp"].(float64)

	// Return an error if the expiration claim is missing or invalid
	if !ok {
		return errors.New("invalid expiration time in token")
	}

	// Convert the expiration timestamp to UnixTime type
	expires := UnixTime(exp)

	// Calculate the time-to-live for the current token
	ttl := expires - UnixNow()

	// Check if the token has sufficient remaining lifespan
	if ttl >= MinRemaining {
		// Log that the token is still valid
		log.Printf("Reusing still-valid accessToken. Expiring in %ds.\n", ttl)
	} else {
		// Log that the token will be refreshed
		log.Println("Refreshing token")

		// Perform the token refresh operation
		err = c.updateAccessToken()

		// Propagate any errors that occur during refresh
		if err != nil {
			return err
		}
	}

	// Indicate successful header creation or refresh
	return nil
}

// updateAccessToken uses the refresh token to obtain a new access token from the authentication service.
func (c *Client) updateAccessToken() error {
	// Initialize a new HTTP client based on default configuration
	client := req.C()

	// Construct the discovery endpoint URL
	discoveryURL := c.HostURL + "/discovery"

	// Send a GET request to retrieve endpoint metadata
	resp, err := client.R().Get(discoveryURL)

	// Return on any network or request error
	if err != nil {
		return err
	}

	// Define a local structure to capture JSON response fields
	var data struct {
		Firebase struct {
			APIKey string `json:"apiKey"`
		} `json:"firebase"`
	}

	// Unmarshal the JSON response into the data structure
	err = resp.UnmarshalJson(&data)

	// Return on JSON parsing errors
	if err != nil {
		return err
	}

	// Initialize the token refresh endpoint with the retrieved API key
	tokenRefreshURL := "https://securetoken.googleapis.com/v1/token?key="

	// Create the JSON body for the refresh request
	jsonBody, err := json.Marshal(map[string]string{
		"refresh_token": c.RefreshToken,
		"grant_type":    "refresh_token",
	})

	// Return on JSON marshaling errors
	if err != nil {
		return err
	}

	// Send a POST request to refresh the access token
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonBody).
		Post(tokenRefreshURL + data.Firebase.APIKey)

	// Handle errors indicating session expiration
	if err != nil {
		return errors.New("your session has expired. A new refresh token is required")
	}

	// Define a local structure to capture the new access token
	var tokenData struct {
		AccessToken string `json:"access_token"`
	}

	// Unmarshal the JSON response into the tokenData structure
	err = response.UnmarshalJson(&tokenData)

	// Return on JSON parsing errors
	if err != nil {
		return err
	}

	// Prefix the token with the Bearer scheme and update the client
	c.Token = "Bearer " + tokenData.AccessToken

	// Indicate successful token update
	return nil
}

// UnixNow returns the current time as a UnixTime value representing seconds since the epoch.
func UnixNow() UnixTime {
	// Retrieve the current Unix timestamp in seconds
	return UnixTime(time.Now().Unix())
}

// CplnClaims represents authentication claims for Cloud Plane users
type CplnClaims struct {
	Name          string `json:"name,omitempty"`
	AuthTime      int64  `json:"auth_time,omitempty"`
	UserId        string `json:"user_id,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`

	jwt.Claims
}

// ServiceAccountToken represents the organization and name extracted from a service account token.
type ServiceAccountToken struct {
	Org  string // The organization part of the token
	Name string // The name part of the token
}

// ParseServiceAccountToken decodes a service account token and extracts the org and name if valid.
func ParseServiceAccountToken(token string) *ServiceAccountToken {
	// Return nil if the input is empty or undefined
	if token == "" {
		return nil
	}

	// Check if the token starts with 's', indicating it's possibly a service account
	if strings.HasPrefix(token, "s") {
		// Split the token by '.' and take the second part (index 1)
		parts := strings.Split(token, ".")
		if len(parts) < 2 {
			// Return nil if the format is invalid
			return nil
		}

		// Decode the base64url-encoded part
		decoded, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			// Return nil if decoding fails
			return nil
		}

		// Split the decoded string by '.' to get org and name
		keyParts := strings.Split(string(decoded), ".")
		if len(keyParts) < 2 {
			// Return nil if the decoded format is invalid
			return nil
		}

		// Return the parsed token as a struct
		return &ServiceAccountToken{
			Org:  keyParts[0], // First part is the organization
			Name: keyParts[1], // Second part is the key name
		}
	}

	// Return nil if the token doesn't start with 's'
	return nil
}
