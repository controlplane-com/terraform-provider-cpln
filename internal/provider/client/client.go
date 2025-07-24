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

	constants "github.com/controlplane-com/terraform-provider-cpln/internal/provider/constants"
)

// DefaultClientEndpoint is the default data service endpoint.
var DefaultClientEndpoint string = "https://api.cpln.io"

// Client - Simple API Client
type Client struct {
	HostURL         string
	Org             string
	HTTPClient      *http.Client
	Token           string
	RefreshToken    string
	ProviderVersion string
}

// NewClient instantiates a new API Client with optional token refresh
func NewClient(org, host, profile, token, refreshToken *string, providerVersion string) (*Client, error) {
	// If host is nil, attempt to extact the host from the CLI profile, will fall back to default if extraction failed
	if host == nil {
		host = ExtractHostFromProfile(profile)
	}

	// Initialize the Client struct with HTTP client, host, org, and tokens
	c := Client{
		// Set the API host URL
		HostURL: *host,
		// Assign the organization identifier
		Org: *org,
		// Configure HTTP client with a 90 second timeout
		HTTPClient: &http.Client{Timeout: 90 * time.Second},
		// Initialize the access token
		Token: *token,
		// Initialize the refresh token
		RefreshToken: *refreshToken,
		// Set the provider version
		ProviderVersion: providerVersion,
	}

	// Check if a refresh token was provided
	if c.RefreshToken != "" {
		// Attempt to set the Authorization header using the refresh token
		err := c.MakeAuthorizationHeader()

		// Handle error from refresh token flow
		if err != nil {
			return nil, fmt.Errorf("unable to obtain access token using the refresh token. Error: %s", err)
		}

		// If no refresh token but no access token either, extract from CLI profile
	} else if c.Token == "" {
		// Invoke profile-based token extraction
		token, err := c.ExtractTokenFromProfile(*profile)

		// Propagate any extraction errors
		if err != nil {
			return nil, err
		}

		// Assign the extracted token to the client
		c.Token = *token
	}

	// Return the configured client instance
	return &c, nil
}

// Define a method on Client to execute HTTP requests with optional content type and tokens
func (c *Client) doRequest(req *http.Request, contentType string, optionalTokens ...string) ([]byte, int, error) {
	// Provide WSL DNS retrieval tip
	// WSL TO GET IP: cat /etc/resolv.conf
	// Include example proxy configuration for debugging
	// os.Setenv("HTTP_PROXY", "http://172.17.80.1:8888")

	// Default to the client’s token for authorization
	clientToken := c.Token

	// Override the token if an optional one is provided
	if len(optionalTokens) > 0 {
		// Use the first provided override token
		clientToken = optionalTokens[0]
	}

	// Set the Authorization header on the request
	req.Header.Set("Authorization", clientToken)

	// Conditionally set the Content-Type header if specified
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Perform the HTTP request using the client’s HTTPClient
	res, err := c.HTTPClient.Do(req)

	// Handle errors that occur during the HTTP round trip
	if err != nil {
		// If a response exists, return its status code with the error
		if res != nil {
			return nil, res.StatusCode, err
		}

		// Otherwise return zero status code with the error
		return nil, 0, err
	}

	// Ensure the response body is closed when done reading
	defer res.Body.Close()

	// Read the full response body into memory
	body, err := io.ReadAll(res.Body)

	// Handle errors encountered while reading the body
	if err != nil {
		return nil, res.StatusCode, err
	}

	// Reject any status codes outside 200, 201, or 202 as errors
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusAccepted {
		// Include both status code and raw body in the error
		return nil, res.StatusCode, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	// Return the successful body bytes and status code
	return body, res.StatusCode, err
}

// Define a method on Client to perform GET requests and decode JSON into a given resource type
func (c *Client) Get(link string, resource interface{}) (interface{}, int, error) {
	// Normalize the link by removing a leading slash if present
	if link[0] == '/' {
		link = link[1:]
	}

	// Build a new HTTP GET request for the specified link
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.HostURL, link), nil)

	// Return immediately if request creation fails
	if err != nil {
		return nil, 0, err
	}

	// Invoke doRequest without a content type to fetch the raw response
	body, code, err := c.doRequest(req, "")

	// Propagate any errors from the HTTP request
	if err != nil {
		return nil, code, err
	}

	// Create a new instance of the expected resource type
	vp := reflect.New(reflect.TypeOf(resource).Elem())

	// Unmarshal the JSON response into the newly created instance
	err = json.Unmarshal(body, vp.Interface())

	// Return upon JSON decoding errors
	if err != nil {
		return nil, code, err
	}

	// Return the populated resource instance and HTTP status code
	return vp.Interface(), code, nil
}

// Define a method on Client to fetch a specific resource by ID and decode JSON into the provided type
func (c *Client) GetResource(id string, resource interface{}) (interface{}, int, error) {
	// Build a new HTTP GET request for the resource endpoint by ID
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, id), nil)

	// Return immediately if request creation fails
	if err != nil {
		return nil, 0, err
	}

	// Invoke doRequest without a content type to fetch the raw response
	body, code, err := c.doRequest(req, "")

	// Propagate any errors from the HTTP request
	if err != nil {
		return nil, code, err
	}

	// Create a new instance of the expected resource type
	vp := reflect.New(reflect.TypeOf(resource).Elem())

	// Unmarshal the JSON response into the newly created instance
	err = json.Unmarshal(body, vp.Interface())

	// Return upon JSON decoding errors
	if err != nil {
		return nil, code, err
	}

	// Return the populated resource instance and HTTP status code
	return vp.Interface(), code, nil
}

// CreateResource sends a POST to create a resource with retry logic on 429 errors.
func (c *Client) CreateResource(resourceType, id string, resource interface{}) (int, error) {
	// Tag the resource as created by Terraform
	c.ForceCreatedByTerraformTag(resource, false)

	// Convert the resource object into JSON bytes
	bodyBytes, err := json.Marshal(resource)

	// Handle JSON marshalling errors immediately
	if err != nil {
		return 0, err
	}

	// Define how many times we'll retry the request
	const maxRetries = 5

	// Set the initial backoff delay
	backoff := 2 * time.Second

	// Attempt the HTTP request up to maxRetries times
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Build a new HTTP POST request for creating the resource
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, resourceType),
			strings.NewReader(string(bodyBytes)),
		)

		// If request construction fails, abort immediately
		if err != nil {
			return 0, err
		}

		// Execute the HTTP request
		_, code, err := c.doRequest(req, "application/json")

		// On HTTP 429 Too Many Requests...
		if err != nil {
			if code == http.StatusTooManyRequests {
				// If the error message mentions quota, return the original error
				if strings.Contains(err.Error(), "quota") {
					return code, err
				}

				// If we've exhausted all retries, break to report failure
				if attempt == maxRetries {
					break
				}

				// Wait for the current backoff duration
				time.Sleep(backoff)

				// Double the backoff for the next attempt
				backoff *= 2

				// Retry the request
				continue
			}

			// For other errors, return immediately
			return code, err
		}

		// On success, return the HTTP status code
		return code, nil
	}

	// Report final failure after all retry attempts
	return 0, fmt.Errorf("create resource %q failed after %d attempts due to HTTP 429", id, maxRetries)
}

// CreateResourceAgent sends a POST to create an Agent with retry logic on 429 errors.
func (c *Client) CreateResourceAgent(resource Agent) (*Agent, int, error) {
	// Tag the Agent as created by Terraform
	c.ForceCreatedByTerraformTag(resource, false)

	// Convert the Agent object into JSON bytes
	bodyBytes, err := json.Marshal(resource)

	// Handle JSON marshalling errors immediately
	if err != nil {
		return nil, 0, err
	}

	// Define how many times we'll retry the request
	const maxRetries = 5

	// Set the initial backoff delay
	backoff := 2 * time.Second

	// Attempt the HTTP request up to maxRetries times
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Build a new HTTP POST request for creating the Agent
		req, err := http.NewRequest(
			http.MethodPost,
			fmt.Sprintf("%s/org/%s/agent", c.HostURL, c.Org),
			strings.NewReader(string(bodyBytes)),
		)

		// If request construction fails, abort immediately
		if err != nil {
			return nil, 0, err
		}

		// Execute the HTTP request
		respBody, code, err := c.doRequest(req, "application/json")
		// On HTTP 429 Too Many Requests...
		if err != nil {
			if code == http.StatusTooManyRequests {
				// If the error message mentions quota, return the original error
				if strings.Contains(err.Error(), "quota") {
					return nil, code, err
				}

				// If we've exhausted all retries, break to report failure
				if attempt == maxRetries {
					break
				}

				// Wait for the current backoff duration
				time.Sleep(backoff)

				// Double the backoff for the next attempt
				backoff *= 2

				// Retry the request
				continue
			}

			// For other errors, return immediately
			return nil, code, err
		}

		// Parse the successful response body into an Agent
		var output Agent

		// Handle JSON unmarshalling errors immediately
		if err := json.Unmarshal(respBody, &output); err != nil {
			return nil, code, err
		}

		// Return the created Agent and HTTP status code
		return &output, code, nil
	}

	// Report final failure after all retry attempts
	return nil, 0, fmt.Errorf("create Agent failed after %d attempts due to HTTP 429", maxRetries)
}

// UpdateResource sends a PATCH to update a resource with retry logic on 429 errors.
func (c *Client) UpdateResource(id string, resource interface{}) (int, error) {
	// Tag the resource as updated by Terraform
	c.ForceCreatedByTerraformTag(resource, true)

	// Convert the resource object into JSON bytes
	bodyBytes, err := json.Marshal(resource)

	// Handle JSON marshalling errors immediately
	if err != nil {
		return 0, err
	}

	// Define how many times we'll retry the request
	const maxRetries = 5

	// Set the initial backoff delay
	backoff := 2 * time.Second

	// Attempt the HTTP request up to maxRetries times
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Build a new HTTP PATCH request for updating the resource
		req, err := http.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, id),
			strings.NewReader(string(bodyBytes)),
		)

		// If request construction fails, abort immediately
		if err != nil {
			return 0, err
		}

		// Execute the HTTP request
		_, code, err := c.doRequest(req, "application/json")

		// On HTTP 429 Too Many Requests...
		if err != nil {
			if code == http.StatusTooManyRequests {
				// If the error message mentions quota, return the original error
				if strings.Contains(err.Error(), "quota") {
					return code, err
				}

				// If we've exhausted all retries, break to report failure
				if attempt == maxRetries {
					break
				}

				// Wait for the current backoff duration
				time.Sleep(backoff)

				// Double the backoff for the next attempt
				backoff *= 2

				// Retry the request
				continue
			}

			// For other errors, return immediately
			return code, err
		}

		// Return on first successful response
		return code, nil
	}

	// Report final failure after all retry attempts
	return 0, fmt.Errorf("update resource %q failed after %d attempts due to HTTP 429", id, maxRetries)
}

// DeleteResource attempts to delete the specified resource by ID with retry logic on 409 and 429 errors.
func (c *Client) DeleteResource(id string) error {
	// Define how many times we'll retry the delete
	const maxRetries = 5

	// Set the initial backoff delay
	backoff := 2 * time.Second

	// Attempt the HTTP delete up to maxRetries times
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Build a new HTTP DELETE request for the resource
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("%s/org/%s/%s", c.HostURL, c.Org, id),
			nil,
		)

		// If request construction fails, abort immediately
		if err != nil {
			return err
		}

		// Execute the HTTP request
		_, code, err := c.doRequest(req, "")

		// On HTTP 409 Conflict...
		if err != nil && code == http.StatusConflict {
			// If we've exhausted all retries, break to report failure
			if attempt == maxRetries {
				break
			}

			// Wait for the current backoff duration
			time.Sleep(backoff)

			// Double the backoff for the next attempt
			backoff *= 2

			// Retry the request
			continue
		}
		// On HTTP 429 Too Many Requests...
		if err != nil && code == http.StatusTooManyRequests {
			// If the error message mentions quota, return the original error
			if strings.Contains(err.Error(), "quota") {
				return err
			}

			// If we've exhausted all retries, break to report failure
			if attempt == maxRetries {
				break
			}

			// Wait for the current backoff duration
			time.Sleep(backoff)

			// Double the backoff for the next attempt
			backoff *= 2

			// Retry the request
			continue
		}

		// If delete succeeded or a non‐retryable error occurred, return immediately
		if err == nil {
			return nil
		}

		// Return any other error immediately
		return err
	}

	// Report final failure after all retry attempts
	return fmt.Errorf("delete resource %q failed after %d attempts", id, maxRetries)
}

// ForceCreatedByTerraformTag Force a tag indicating resource was created by Terraform
func (c *Client) ForceCreatedByTerraformTag(resource interface{}, isUpdate bool) {
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
		// Determine the tags field name
		tagsFieldName := "Tags"

		// If this is an update, then change the field name
		if isUpdate {
			tagsFieldName = "TagsReplace"
		}

		// Get the Tags field from the Base struct
		tagsField := baseField.FieldByName(tagsFieldName)

		// Check if Tags is nil, and initialize if necessary
		if tagsField.IsNil() {
			newTags := make(map[string]interface{})
			tagsField.Set(reflect.ValueOf(&newTags))
		}

		// Add a new key-value pair to the Tags map
		tags := tagsField.Interface().(*map[string]interface{})
		(*tags)[constants.ManagedByTerraformTagKey] = "true"
		(*tags)[constants.TerraformVersionTagKey] = c.ProviderVersion
	}
}

// ExtractTokenFromProfile runs the cpln CLI to fetch the access token for the specified profile
func (c *Client) ExtractTokenFromProfile(profileName string) (*string, error) {
	// Create the command
	cmd := exec.Command("cpln", "profile", "token", profileName)

	// Create buffers for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Attempt to execute the cpln command and handle errors accordingly
	err := cmd.Run()
	if err != nil {
		errorMessage := ""
		stdoutString := stdout.String()
		stderrString := stderr.String()

		// Add stdout to the error message
		if len(strings.TrimSpace(stdoutString)) != 0 {
			errorMessage = fmt.Sprintf("Stdout: %s", stdoutString)
		}

		// Add stderr to the error message
		if len(strings.TrimSpace(stderrString)) != 0 {
			// Add a white space if error message is not empty so we can cleanly add stderr to the message
			if len(errorMessage) != 0 {
				errorMessage = fmt.Sprintf("%s ", errorMessage)
			}

			errorMessage = fmt.Sprintf("%sStderr: %s", errorMessage, stderrString)
		}

		return nil, fmt.Errorf("unable to obtain access token from profile '%s'. Verify cpln is installed and added to PATH. Error: %s. %s", profileName, err, errorMessage)
	}

	// Convert stdout from bytes to a string
	stringifiedStdout := stdout.String()

	// Handle the case where the token is empty
	if strings.TrimSpace(string(stringifiedStdout)) == "" {
		return nil, fmt.Errorf("empty access token")
	}

	// Finalize the token
	token := strings.TrimSuffix(string(stringifiedStdout), "\n")

	// Return a pointer to the token
	return &token, nil
}

// ExtractHostFromProfile runs the cpln CLI to fetch the endpoint host for the specified profile, fallback to DefaultClientEndpoint
func ExtractHostFromProfile(profileName *string) *string {
	// Return default endpoint if no profile name provided
	if profileName == nil || *profileName == "" {
		return &DefaultClientEndpoint
	}

	// Prepare the cpln profile get command to retrieve profile details
	cmd := exec.Command("cpln", "profile", "get", *profileName, "-o", "json")

	// Initialize buffers to capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Redirect command stdout into stdout buffer
	cmd.Stdout = &stdout

	// Redirect command stderr into stderr buffer
	cmd.Stderr = &stderr

	// Execute the command and capture any error
	err := cmd.Run()

	// On error executing the CLI, fallback to default endpoint
	if err != nil {
		return &DefaultClientEndpoint
	}

	// Convert stdout bytes to string
	output := stdout.String()

	// Unmarshal stdout JSON into a slice of maps
	var profiles []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &profiles); err != nil {
		return &DefaultClientEndpoint
	}

	// Ensure at least one profile object is present
	if len(profiles) == 0 {
		return &DefaultClientEndpoint
	}

	// Access the first profile object
	first := profiles[0]

	// Retrieve the nested request object
	reqObj, ok := first["request"].(map[string]interface{})
	if !ok {
		return &DefaultClientEndpoint
	}

	// Extract the endpoint property from the request object
	endpoint, ok := reqObj["endpoint"].(string)
	if !ok || endpoint == "" {
		return &DefaultClientEndpoint
	}

	// Return the extracted endpoint host
	return &endpoint
}

// RemoveManagedByTerraformTag removes the managed-by-Terraform tag from the given base struct.
func (c *Client) RemoveManagedByTerraformTag(base *Base) {
	// Dereference TagsReplace to work with the actual map
	tagsReplace := *base.TagsReplace

	// Remove the managedByTerraform tag
	delete(tagsReplace, "cpln/managedByTerraform")

	// Update the pointer to reference the modified map
	base.TagsReplace = &tagsReplace
}
