package cpln

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

// Client - Simple API Client
type Client struct {
	HostURL      string
	Org          string
	HTTPClient   *http.Client
	Token        string
	RefreshToken string
}

// NewClient - Instantiate a new API Client
func NewClient(org, host, profile, token, refreshToken *string) (*Client, error) {

	c := Client{
		HTTPClient:   &http.Client{Timeout: 90 * time.Second},
		HostURL:      *host,
		Org:          *org,
		Token:        *token,
		RefreshToken: *refreshToken,
	}

	if c.RefreshToken != "" {

		err := c.MakeAuthorizationHeader()

		// Handle error
		if err != nil {
			return nil, fmt.Errorf("unable to obtain access token using the refresh token. Error: %s", err)
		}
	} else if c.Token == "" {
		// Create command
		cmd := exec.Command("cpln", "profile", "token", *profile)

		// Create buffers for stdout and stderr
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Run the command
		err := cmd.Run()

		// Handle error
		if err != nil {
			return nil, fmt.Errorf("unable to obtain access token. Verify cpln is installed and added to PATH. Error: %s. Stderr: %s. Stderr: %s", err, stderr.String(), stderr.String())
		}

		// Handle token
		out := stdout.String()

		if strings.TrimSpace(string(out)) == "" {
			return nil, fmt.Errorf("empty access token")
		}

		c.Token = strings.TrimSuffix(string(out), "\n")
	}

	// log.Printf("[INFO] New Client instantiated. Endpoint: %s. Org: %s. Profile: %s", *host, *org, *profile)

	return &c, nil
}

func (c *Client) doRequest(req *http.Request, contentType string, optionalTokens ...string) ([]byte, int, error) {

	// WSL TO GET IP: cat /etc/resolv.conf
	// os.Setenv("HTTP_PROXY", "http://172.17.80.1:8888")

	clientToken := c.Token

	if len(optionalTokens) > 0 {
		clientToken = optionalTokens[0]
	}

	req.Header.Set("Authorization", clientToken)

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	res, err := c.HTTPClient.Do(req)

	if err != nil {

		if res != nil {
			return nil, res.StatusCode, err
		}

		return nil, 0, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, res.StatusCode, err
	}

	// log.Printf("[INFO] Status Code: %d. URL: %s. Method: %s", res.StatusCode, req.URL, req.Method)

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusAccepted {
		return nil, res.StatusCode, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, res.StatusCode, err
}

func (c *Client) Get(link string, resource interface{}) (interface{}, int, error) {

	// Remove leading slash
	if link[0] == '/' {
		link = link[1:]
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.HostURL, link), nil)

	if err != nil {
		return nil, 0, err
	}

	body, code, err := c.doRequest(req, "")

	if err != nil {
		return nil, code, err
	}

	vp := reflect.New(reflect.TypeOf(resource).Elem())
	err = json.Unmarshal(body, vp.Interface())

	if err != nil {
		return nil, code, err
	}

	return vp.Interface(), code, nil
}

func (c *Client) GetResource(id string, resource interface{}) (interface{}, int, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, id), nil)

	if err != nil {
		return nil, 0, err
	}

	body, code, err := c.doRequest(req, "")
	if err != nil {
		return nil, code, err
	}

	vp := reflect.New(reflect.TypeOf(resource).Elem())

	err = json.Unmarshal(body, vp.Interface())
	if err != nil {
		return nil, code, err
	}

	return vp.Interface(), code, nil
}

func (c *Client) CreateResource(resourceType, id string, resource interface{}) (int, error) {

	c.ForceCreatedByTerraformTag(resource)

	g, err := json.Marshal(resource)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, resourceType), strings.NewReader(string(g)))
	if err != nil {
		return 0, err
	}

	_, code, err := c.doRequest(req, "application/json")
	if err != nil {
		return code, err
	}

	return code, nil
}

func (c *Client) CreateResourceAgent(resource Agent) (*Agent, int, error) {

	c.ForceCreatedByTerraformTag(resource)

	g, err := json.Marshal(resource)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/org/%s/agent", c.HostURL, c.Org), strings.NewReader(string(g)))
	if err != nil {
		return nil, 0, err
	}

	body, code, err := c.doRequest(req, "application/json")
	if err != nil {
		return nil, code, err
	}

	output := Agent{}

	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, code, err
	}

	return &output, code, nil
}

func (c *Client) UpdateResource(id string, resource interface{}) (int, error) {

	c.ForceCreatedByTerraformTag(resource)

	g, err := json.Marshal(resource)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, id), strings.NewReader(string(g)))
	if err != nil {
		return 0, err
	}

	_, code, err := c.doRequest(req, "application/json")
	if err != nil {
		return code, err
	}

	return code, nil
}

func (c *Client) DeleteResource(id string) error {

	// Add a delay to allow any referenced resources to be deleted.
	time.Sleep(5 * time.Second)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, id), nil)
	if err != nil {
		return err
	}

	_, _, err = c.doRequest(req, "")
	if err != nil {
		return err
	}

	return nil
}

// Force a tag indicating resource was created by Terraform
func (c *Client) ForceCreatedByTerraformTag(resource interface{}) {
	// Use reflection to get the value of the resource
	val := reflect.ValueOf(resource)

	// Dereference pointer if necessary
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Look for the embedded Base struct
	baseField := val.FieldByName("Base")

	// Check if it is valid
	if baseField.IsValid() {
		// Get the Tags field from the Base struct
		tagsField := baseField.FieldByName("Tags")

		// Check if Tags is nil, and initialize if necessary
		if tagsField.IsNil() {
			newTags := make(map[string]interface{})
			tagsField.Set(reflect.ValueOf(&newTags))
		}

		// Add a new key-value pair to the Tags map
		tags := tagsField.Interface().(*map[string]interface{})
		(*tags)["cpln/managedByTerraform"] = "true"
	}
}
