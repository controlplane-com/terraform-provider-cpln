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
