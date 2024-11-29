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
	HTTPClient   *http.Client
	HostURL      string
	Org          string
	Token        string
	RefreshToken string
}

// NewClient - Instantiate a new API Client
func NewClient(org, host, profile, token, refreshToken string) (*Client, error) {

	// Initialize a new client with the provided HTTP client, host URL, organization, token, and refresh token.
	c := Client{
		HTTPClient:   &http.Client{Timeout: 90 * time.Second},
		HostURL:      host,
		Org:          org,
		Token:        token,
		RefreshToken: refreshToken,
	}

	// Check if a refresh token is available.
	if c.RefreshToken != "" {
		// Attempt to generate the authorization header using the refresh token.
		if err := c.MakeAuthorizationHeader(); err != nil {
			return nil, fmt.Errorf("failed to generate authorization header using the provided refresh token: %v. Verify the refresh token and try again.", err)
		}
	} else if c.Token == "" { // If no access token or refresh token, attempt to retrieve it using CLI.

		// Construct the CLI command to fetch the access token based on the specified profile.
		cmd := exec.Command("cpln", "profile", "token", profile)

		// Prepare buffers to capture standard output and error.
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Execute the command and capture any errors.
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("unable to retrieve access token via 'cpln' command. Ensure 'cpln' is correctly installed and available in your PATH. Command error: %v. Additional details: %s", err, stderr.String())
		}

		// Retrieve the command output (token) and trim any whitespace.
		outputToken := strings.TrimSpace(stdout.String())
		if outputToken == "" {
			return nil, fmt.Errorf("retrieved an empty access token from 'cpln' command")
		}

		// Assign the retrieved token to the client.
		c.Token = outputToken
	}

	// Print client configuration if needed for debugging
	// log.Printf("[INFO] New Client instantiated. Endpoint: %s. Org: %s. Profile: %s", *host, *org, *profile)

	return &c, nil
}

// Get fetches the resource from the specified link, unmarshaling the response into the given structure.
func (c *Client) Get(link string, resource interface{}) (interface{}, int, error) {

	// If the link starts with a leading slash, remove it to prevent double slashes in the URL.
	if link[0] == '/' {
		link = link[1:]
	}

	// Construct a new HTTP GET request with the specified link appended to the host URL.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.HostURL, link), nil)
	if err != nil {
		return nil, 0, err // Return an error if request creation fails.
	}

	// Execute the HTTP request using a helper method, and capture the response body and status code.
	body, code, err := c.doRequest(req, "")
	if err != nil {
		return nil, code, err // Return an error if the request fails.
	}

	// Initialize a new instance of the resource type dynamically using reflection.
	vp := reflect.New(reflect.TypeOf(resource).Elem())

	// Unmarshal the JSON response body into the resource structure.
	err = json.Unmarshal(body, vp.Interface())
	if err != nil {
		return nil, code, err // Return an error if JSON unmarshaling fails.
	}

	// Return the populated resource instance and status code.
	return vp.Interface(), code, nil
}

// GetResource sends an HTTP GET request to retrieve a specific resource by its ID from the organization.
// It unmarshals the JSON response into the provided resource structure using reflection.
func (c *Client) GetResource(id string, resource interface{}) (interface{}, int, error) {

	// Construct the request URL with the specified ID and organization, using the client's HostURL.
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, id), nil)
	if err != nil {
		return nil, 0, err // Return an error if request creation fails.
	}

	// Execute the HTTP request and retrieve the response body and status code.
	body, code, err := c.doRequest(req, "")
	if err != nil {
		return nil, code, err // Return an error if the request fails.
	}

	// Create a new instance of the resource's type using reflection for dynamic unmarshaling.
	vp := reflect.New(reflect.TypeOf(resource).Elem())

	// Unmarshal the JSON response body into the new resource instance.
	err = json.Unmarshal(body, vp.Interface())
	if err != nil {
		return nil, code, err // Return an error if JSON unmarshaling fails.
	}

	// Return the populated resource instance and the HTTP status code.
	return vp.Interface(), code, nil
}

// CreateResource sends an HTTP POST request to create a new resource of the specified type and ID in the organization.
func (c *Client) CreateResource(resourceType, id string, resource interface{}) (int, error) {

	// Add a Terraform-specific tag to the resource to indicate it was created by Terraform.
	c.ForceCreatedByTerraformTag(resource)

	// Marshal the resource into JSON format for the request body.
	g, err := json.Marshal(resource)
	if err != nil {
		return 0, err // Return an error if JSON marshalling fails.
	}

	// Construct the HTTP POST request to create the resource at the target URL.
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, resourceType), strings.NewReader(string(g)))
	if err != nil {
		return 0, err // Return an error if request creation fails.
	}

	// Execute the request with JSON content type and capture the status code.
	_, code, err := c.doRequest(req, "application/json")
	if err != nil {
		return code, err // Return an error if the request fails.
	}

	return code, nil // Return the HTTP status code upon successful resource creation.
}

// CreateResourceAgent creates a new Agent resource by sending an HTTP POST request and returns the created agent.
func (c *Client) CreateResourceAgent(resource Agent) (*Agent, int, error) {

	// Add a Terraform-specific tag to the Agent resource.
	c.ForceCreatedByTerraformTag(resource)

	// Marshal the Agent into JSON format for the request body.
	g, err := json.Marshal(resource)
	if err != nil {
		return nil, 0, err // Return an error if JSON marshalling fails.
	}

	// Construct the HTTP POST request to create the agent at the target URL.
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/org/%s/agent", c.HostURL, c.Org), strings.NewReader(string(g)))
	if err != nil {
		return nil, 0, err // Return an error if request creation fails.
	}

	// Execute the request with JSON content type, capturing the response body and status code.
	body, code, err := c.doRequest(req, "application/json")
	if err != nil {
		return nil, code, err // Return an error if the request fails.
	}

	// Unmarshal the JSON response into an Agent struct.
	output := Agent{}
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, code, err // Return an error if unmarshaling fails.
	}

	return &output, code, nil // Return the created agent and status code.
}

// UpdateResource sends an HTTP PATCH request to update an existing resource with the given ID.
func (c *Client) UpdateResource(id string, resource interface{}) (int, error) {

	// Add a Terraform-specific tag to the resource to indicate it was created by Terraform.
	c.ForceCreatedByTerraformTag(resource)

	// Marshal the resource into JSON format for the request body.
	g, err := json.Marshal(resource)
	if err != nil {
		return 0, err // Return an error if JSON marshalling fails.
	}

	// Construct the HTTP PATCH request to update the resource at the target URL.
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, id), strings.NewReader(string(g)))
	if err != nil {
		return 0, err // Return an error if request creation fails.
	}

	// Execute the request with JSON content type and capture the status code.
	_, code, err := c.doRequest(req, "application/json")
	if err != nil {
		return code, err // Return an error if the request fails.
	}

	return code, nil // Return the HTTP status code upon successful resource update.
}

// DeleteResource sends an HTTP DELETE request to remove a resource by ID from the organization.
func (c *Client) DeleteResource(id string) error {

	// Introduce a delay to allow any related resources to complete deletion before proceeding.
	time.Sleep(5 * time.Second)

	// Construct the HTTP DELETE request to remove the resource at the target URL.
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, id), nil)
	if err != nil {
		return err // Return an error if request creation fails.
	}

	// Execute the delete request.
	_, _, err = c.doRequest(req, "")
	if err != nil {
		return err // Return an error if the request fails.
	}

	return nil // Return nil if the deletion was successful.
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
		tagsReplaceField := baseField.FieldByName("TagsReplace")

		// Ensure the Tags field is valid, can be set, and is addressable
		if !tagsField.IsNil() {
			// Set a new key-value pair in the Tags map
			tags := tagsField.Interface().(*map[string]interface{})
			(*tags)["cpln/managedByTerraform"] = "true"
		}

		if !tagsReplaceField.IsNil() {
			// Set a new key-value pair in the TagsReplace map
			tags := tagsReplaceField.Interface().(*map[string]interface{})
			(*tags)["cpln/managedByTerraform"] = "true"
		}
	}
}

func (c *Client) doRequest(req *http.Request, contentType string, optionalTokens ...string) ([]byte, int, error) {

	// To retrieve the IP address for WSL, use: cat /etc/resolv.conf
	// Uncomment the following line to set the HTTP proxy, if required.
	// os.Setenv("HTTP_PROXY", "http://172.17.80.1:8888")

	// Set the token for authorization; if an optional token is provided, use it instead.
	clientToken := c.Token
	if len(optionalTokens) > 0 {
		clientToken = optionalTokens[0]
	}

	// Add the Authorization header to the request.
	req.Header.Set("Authorization", clientToken)

	// If a Content-Type is specified, set it in the request headers.
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Send the HTTP request using the configured HTTP client.
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		// If the response is available, return its status code with the error.
		if res != nil {
			return nil, res.StatusCode, err
		}
		// Return a generic error if the response is unavailable.
		return nil, 0, err
	}
	defer res.Body.Close()

	// Read the entire response body.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	// Log the response status code, URL, and method if needed for debugging.
	// log.Printf("[INFO] Status Code: %d. URL: %s. Method: %s", res.StatusCode, req.URL, req.Method)

	// Check for successful response codes (200 OK, 201 Created, 202 Accepted).
	// Return an error if the response code indicates a failure.
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusAccepted {
		return nil, res.StatusCode, fmt.Errorf("received unexpected status code: %d, response body: %s", res.StatusCode, body)
	}

	// Return the response body and status code with no errors.
	return body, res.StatusCode, nil
}
