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

// MakeAuthorizationHeader creates an authorization header for the given profile.
func (c *Client) MakeAuthorizationHeader() error {

	if c.RefreshToken == "" {
		return errors.New("empty refresh token")
	}

	token, err := jwt.Parse(c.Token, func(token *jwt.Token) (interface{}, error) {
		// Define the method of signing, for example, if it's HMAC:
		// return []byte("your-256-bit-secret"), nil

		// If the signing method is unknown, return an error.
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	})

	if err != nil {
		err = c.updateAccessToken()
		if err != nil {
			return err
		}
		return nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid expiration time in token")
	}

	expires := UnixTime(exp)
	ttl := expires - UnixNow()
	if ttl >= MinRemaining {
		log.Printf("Reusing still-valid accessToken. Expiring in %ds.\n", ttl)
	} else {
		log.Println("Refreshing token")
		err = c.updateAccessToken()
		if err != nil {
			return err
		}
	}

	return nil
}

// updateAccessToken updates the access token for the given profile.
func (c *Client) updateAccessToken() error {

	client := req.C()
	discoveryURL := c.HostURL + "/discovery"

	resp, err := client.R().Get(discoveryURL)
	if err != nil {
		return err
	}

	var data struct {
		Firebase struct {
			APIKey string `json:"apiKey"`
		} `json:"firebase"`
	}
	err = resp.UnmarshalJson(&data)
	if err != nil {
		return err
	}

	tokenRefreshURL := "https://securetoken.googleapis.com/v1/token?key="

	jsonBody, err := json.Marshal(map[string]string{
		"refresh_token": c.RefreshToken,
		"grant_type":    "refresh_token",
	})
	if err != nil {
		return err
	}

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonBody).
		Post(tokenRefreshURL + data.Firebase.APIKey)

	if err != nil {
		return errors.New("your session has expired. A new refresh token is required")
	}

	var tokenData struct {
		AccessToken string `json:"access_token"`
	}
	err = response.UnmarshalJson(&tokenData)
	if err != nil {
		return err
	}

	c.Token = "Bearer " + tokenData.AccessToken
	return nil
}

// UnixNow returns the current Unix time.
func UnixNow() UnixTime {
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
