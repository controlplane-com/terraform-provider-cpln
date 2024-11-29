package cpln

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imroc/req/v3"
)

// UnixTime represents a Unix timestamp.
type UnixTime int64

// AppConfig contains configuration details for the application, including Firebase settings.
type AppConfig struct {
	Firebase FirebaseConfig `json:"firebase"` // Firebase configuration.
}

// FirebaseConfig holds the API key for Firebase services.
type FirebaseConfig struct {
	APIKey string `json:"apiKey"` // The API key for Firebase services.
}

// TokenData represents the structure for storing an access token.
var TokenData struct {
	AccessToken string `json:"access_token"`
}

// MinRemaining defines the minimum remaining time before token refresh.
const MinRemaining UnixTime = 10 * 60

// TokenRefreshURL is the endpoint for refreshing Firebase ID tokens.
const TokenRefreshURL string = "https://securetoken.googleapis.com/v1/token?key="

// MakeAuthorizationHeader creates an authorization header for the given profile.
func (c *Client) MakeAuthorizationHeader() error {

	// Check if the refresh token is empty; if so, return an error.
	if c.RefreshToken == "" {
		return errors.New("empty refresh token")
	}

	// Parse the JWT access token to validate its structure and signature.
	token, err := jwt.Parse(c.Token, func(token *jwt.Token) (interface{}, error) {
		// Define the method of signing, for example, if it's HMAC:
		// return []byte("your-256-bit-secret"), nil

		// If the signing method is unknown, return an error.
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	})

	// If there's an error parsing the token, attempt to update the access token.
	if err != nil {
		if err := c.updateAccessToken(); err != nil {
			return err
		}
		return nil
	}

	// Extract claims from the token and assert their type.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Retrieve the expiration time ('exp' claim) from the token.
	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.New("invalid expiration time in token")
	}

	// Calculate the token's time-to-live (TTL) by comparing its expiration to the current time.
	expires := UnixTime(exp)
	ttl := expires - UnixNow()

	// If the token's TTL is greater than or equal to the minimum required, reuse it.
	if ttl >= MinRemaining {
		log.Printf("Reusing still-valid access token. Expiring in %ds.\n", ttl)
	} else {
		// If the token is nearing expiration, refresh it.
		log.Println("Refreshing token")
		if err := c.updateAccessToken(); err != nil {
			return err
		}
	}

	return nil
}

// updateAccessToken updates the access token for the given profile.
func (c *Client) updateAccessToken() error {

	// Initialize a new HTTP client using the request context.
	client := req.C()

	// Construct the discovery URL by appending the discovery endpoint to the host URL.
	discoveryURL := c.HostURL + "/discovery"

	// Send a GET request to the discovery URL to retrieve configuration data.
	resp, err := client.R().Get(discoveryURL)
	if err != nil {
		return err
	}

	// Define a variable to hold the application configuration data.
	var data AppConfig

	// Unmarshal the JSON response into the AppConfig struct.
	if err := resp.UnmarshalJson(&data); err != nil {
		return err
	}

	// Prepare the JSON payload for the token refresh request.
	jsonBody, err := json.Marshal(map[string]string{
		"refresh_token": c.RefreshToken,
		"grant_type":    "refresh_token",
	})
	if err != nil {
		return err
	}

	// Send a POST request to the token refresh URL with the JSON payload.
	// The URL is constructed by appending the Firebase API key to the TokenRefreshURL.
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonBody).
		Post(TokenRefreshURL + data.Firebase.APIKey)
	if err != nil {
		return errors.New("your session has expired. A new refresh token is required")
	}

	// Unmarshal the JSON response into the TokenData struct to extract the access token.
	if err := response.UnmarshalJson(&TokenData); err != nil {
		return err
	}

	// Update the client's token with the new access token, prefixed with "Bearer ".
	c.Token = "Bearer " + TokenData.AccessToken

	return nil
}

// UnixNow returns the current Unix time.
func UnixNow() UnixTime {
	return UnixTime(time.Now().Unix())
}
