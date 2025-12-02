package cpln

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	// MarketplaceProductionURL is the production marketplace service endpoint
	MarketplaceProductionURL = "https://marketplace.marketplace.services.controlplane.site"
	// MarketplaceStagingURL is the staging marketplace service endpoint
	MarketplaceStagingURL = "https://marketplace.staging.marketplace.services.controlplane.site"
)

// MarketplaceRelease represents information about an installed helm release.
type MarketplaceRelease struct {
	Name      string
	Template  string
	Version   string
	Gvc       *string
	Values    string
	Resources *[]HelmReleaseResource
}

// MarketplaceTemplate represents a catalog template item from the marketplace.
type MarketplaceTemplate struct {
	// Name of the template (e.g., "postgres", "redis", "nginx")
	Name *string `json:"name,omitempty"`
	// Versions is a map of version strings to version metadata
	Versions *map[string]MarketplaceVersion `json:"versions,omitempty"`
}

// MarketplaceVersion represents a specific version of a catalog template.
type MarketplaceVersion struct {
	// CreatesGvc indicates whether this template version creates its own GVC.
	// When true, the GVC parameter should be sent as empty string
	CreatesGvc *bool `json:"createsGvc,omitempty"`
}

// MarketplaceInstallRequest represents a request to install or upgrade a helm release.
type MarketplaceInstallRequest struct {
	Org      *string `json:"org,omitempty"`
	Gvc      *string `json:"gvc,omitempty"`
	Name     *string `json:"name,omitempty"`
	Template *string `json:"template,omitempty"`
	Version  *string `json:"version,omitempty"`
	Values   *string `json:"values,omitempty"`
	Action   *string `json:"action,omitempty"`
}

// MarketplaceTemplateRequest represents a request to template a release without installing.
type MarketplaceTemplateRequest struct {
	Org      *string `json:"org,omitempty"`
	Gvc      *string `json:"gvc,omitempty"`
	Name     *string `json:"name,omitempty"`
	Template *string `json:"template,omitempty"`
	Version  *string `json:"version,omitempty"`
	Values   *string `json:"values,omitempty"`
}

// MarketplaceUninstallRequest represents a request to uninstall a helm release.
type MarketplaceUninstallRequest struct {
	Org  *string `json:"org,omitempty"`
	Name *string `json:"name,omitempty"`
}

// MarketplaceHelmResponse represents the response from all Helm operations.
type MarketplaceHelmResponse struct {
	// Message contains the complete helm command output
	Message *string `json:"message,omitempty"`
}

// HelmRelease represents the decoded helm release secret data structure.
type HelmRelease struct {
	Info        HelmReleaseInfo `json:"info"`
	ValuesFiles *[]string       `json:"valuesFiles"`
}

type HelmReleaseInfo struct {
	Resources []HelmReleaseResource `json:"resources"`
}

type HelmReleaseResource struct {
	ID       string `json:"id"`
	Kind     string `json:"kind"`
	Version  int    `json:"version"`
	Link     string `json:"link"`
	Template Base   `json:"template"`
}

// getMarketplaceURL returns the appropriate marketplace URL based on the client configuration.
// If the client is configured for the test environment (https://api.test.cpln.io) or
// if the TF_ACC environment variable is set, it returns the staging marketplace URL.
// Otherwise, it returns the production marketplace URL.
func (c *Client) getMarketplaceURL() string {
	// Check if using test API endpoint
	if c.HostURL == TestClientEndpoint {
		return MarketplaceStagingURL
	}

	// Check if TF_ACC environment variable is set (acceptance testing)
	if os.Getenv("TF_ACC") != "" {
		return MarketplaceStagingURL
	}

	// Default to production marketplace URL
	return MarketplaceProductionURL
}

// GetMarketplaceTemplate gets details for a specific catalog template from the marketplace.
func (c *Client) GetMarketplaceTemplate(templateName string) (*MarketplaceTemplate, error) {
	// Construct the HTTP request to fetch template details
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/template/%s", c.getMarketplaceURL(), templateName), nil)
	if err != nil {
		// Return error if request construction fails
		return nil, err
	}

	// Execute the request through the client's doRequest method
	body, _, err := c.doRequest(req, "")
	if err != nil {
		// Return error if the API call fails
		return nil, err
	}

	// Parse the JSON response into a MarketplaceTemplate struct
	var template MarketplaceTemplate
	err = json.Unmarshal(body, &template)
	if err != nil {
		// Return error if JSON parsing fails
		return nil, err
	}

	// Return the populated template structure
	return &template, nil
}

// InstallMarketplaceRelease installs or upgrades a marketplace release.
func (c *Client) InstallMarketplaceRelease(request MarketplaceInstallRequest) (*MarketplaceHelmResponse, error) {
	// Marshal the request struct into JSON bytes
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		// Return error if JSON marshaling fails
		return nil, err
	}

	// Construct the HTTP POST request with the JSON body
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/helm/install", c.getMarketplaceURL()),
		bytes.NewReader(bodyBytes),
	)

	// Return error if request construction fails
	if err != nil {
		return nil, err
	}

	// Execute the request with application/json content type
	respBody, _, err := c.doRequest(req, "application/json")
	if err != nil {
		// Return error if the API call fails
		return nil, err
	}

	// Parse the JSON response into a MarketplaceHelmResponse struct
	var response MarketplaceHelmResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		// Return error if JSON parsing fails
		return nil, err
	}

	// Return the helm install/upgrade output
	return &response, nil
}

// UninstallMarketplaceRelease uninstalls a marketplace release.
func (c *Client) UninstallMarketplaceRelease(request MarketplaceUninstallRequest) (*MarketplaceHelmResponse, error) {
	// Marshal the request struct into JSON bytes
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		// Return error if JSON marshaling fails
		return nil, err
	}

	// Construct the HTTP POST request with the JSON body
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/helm/uninstall", c.getMarketplaceURL()),
		bytes.NewReader(bodyBytes),
	)

	// Return error if request construction fails
	if err != nil {
		return nil, err
	}

	// Execute the request with application/json content type
	respBody, _, err := c.doRequest(req, "application/json")
	if err != nil {
		// Return error if the API call fails
		return nil, err
	}

	// Parse the JSON response into a MarketplaceHelmResponse struct
	var response MarketplaceHelmResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		// Return error if JSON parsing fails
		return nil, err
	}

	// Return the helm uninstall output
	return &response, nil
}

// GetMarketplaceRelease queries the Control Plane API for helm release secrets to get the current state of an installed marketplace release.
func (c *Client) GetMarketplaceRelease(releaseName string, query Query) (*MarketplaceRelease, int, error) {
	// Find the latest helm release secret using the provided query
	latestSecret, code, err := c.findLatestHelmReleaseSecret(query)
	if err != nil || code == 404 {
		return nil, code, err
	}

	// Dereference the tags pointer
	tags := *latestSecret.Tags

	// Initialize the release info structure with basic data from tags
	info := &MarketplaceRelease{
		Name:     releaseName,
		Template: tags["cpln/marketplace-template"].(string),
		Version:  tags["cpln/marketplace-template-version"].(string),
	}

	// Extract GVC name from cpln/marketplace-gvc tag if it exists
	if gvc, ok := tags["cpln/marketplace-gvc"].(string); ok {
		info.Gvc = &gvc
	}

	// Reveal the secret to access its encoded data
	revealedSecret, code, err := c.revealSecret(*latestSecret.Name)
	if err != nil {
		return nil, code, err
	}

	// Decode the helm release data from the revealed secret
	helmRelease, err := c.decodeHelmReleaseData(revealedSecret.Data)
	if err != nil {
		return nil, 0, err
	}

	// Set resources from the decoded helm release
	info.Resources = &helmRelease.Info.Resources

	// Set values from the first ValuesFiles entry if available
	if helmRelease.ValuesFiles != nil && len(*helmRelease.ValuesFiles) > 0 {
		// Normalize trailing whitespace to match the plan modifier behavior
		info.Values = strings.TrimRight((*helmRelease.ValuesFiles)[0], "\n\r\t ")
	}

	// Return the complete release information including metadata, resources, and values
	return info, 0, nil
}

// findLatestHelmReleaseSecret queries for helm release secrets and returns the latest one by version tag.
func (c *Client) findLatestHelmReleaseSecret(query Query) (*Secret, int, error) {
	// Marshal the query struct into JSON bytes for the POST request
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	// Construct HTTP POST request to the secret query endpoint
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/org/%s/secret/-query", c.HostURL, c.Org),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, 0, err
	}

	// Execute the query request to find helm release secrets
	body, _, err := c.doRequest(req, "application/json")
	if err != nil {
		return nil, 0, err
	}

	// Parse the query response into a Secrets struct (list of secrets)
	var result Secrets
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, 0, err
	}

	// Check if any secrets were found for this release
	if len(result.Items) == 0 {
		return nil, 404, nil
	}

	// Find the latest secret by version tag
	latestSecret := findLatestSecretByVersion(result.Items)

	return &latestSecret, 0, nil
}

// findLatestSecretByVersion sorts secrets by version tag and returns the one with the highest version.
func findLatestSecretByVersion(secrets []Secret) Secret {
	// Sort secrets by version tag to get the latest release
	sort.Slice(secrets, func(i, j int) bool {
		versionI := extractVersionFromTags(secrets[i].Tags)
		versionJ := extractVersionFromTags(secrets[j].Tags)
		return versionI < versionJ
	})

	// Return the latest release secret (highest version number after sorting)
	return secrets[len(secrets)-1]
}

// revealSecret reveals a secret by name and returns the revealed secret data.
func (c *Client) revealSecret(secretName string) (*Secret, int, error) {
	// Construct HTTP GET request to reveal the secret
	revealReq, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/org/%s/secret/%s/-reveal", c.HostURL, c.Org, secretName),
		nil,
	)
	if err != nil {
		return nil, 0, err
	}

	// Execute the reveal request to get the secret's decrypted data
	revealBody, code, err := c.doRequest(revealReq, "")
	if err != nil {
		return nil, code, err
	}

	// Parse the revealed secret response
	var revealedSecret Secret
	err = json.Unmarshal(revealBody, &revealedSecret)
	if err != nil {
		return nil, 0, err
	}

	return &revealedSecret, code, nil
}

// decodeHelmReleaseData decodes the helm release data from a secret payload.
func (c *Client) decodeHelmReleaseData(secretData *interface{}) (*HelmRelease, error) {
	// Return nil if secret data is not present
	if secretData == nil {
		return nil, fmt.Errorf("secret data is nil")
	}

	// Cast the data interface to a map to access the "payload" field
	dataMap, ok := (*secretData).(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("secret data is not a map")
	}

	// Get the base64-encoded payload data
	releaseDataEncoded, ok := dataMap["payload"].(string)
	if !ok {
		return nil, fmt.Errorf("payload field is missing or not a string")
	}

	// Decode the base64 string to get the gzipped data
	releaseDataBytes, err := base64.StdEncoding.DecodeString(releaseDataEncoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 payload: %w", err)
	}

	// Create a gzip reader to decompress the data
	gzipReader, err := gzip.NewReader(bytes.NewReader(releaseDataBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	// Read the gunzipped data
	gunzippedData, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read gunzipped data: %w", err)
	}

	// Parse the gunzipped JSON into HelmRelease struct
	var helmRelease HelmRelease
	if err := json.Unmarshal(gunzippedData, &helmRelease); err != nil {
		return nil, fmt.Errorf("failed to unmarshal helm release data: %w", err)
	}

	return &helmRelease, nil
}

// extractVersionFromTags extracts the version number from a secret's tags.
func extractVersionFromTags(tags *map[string]interface{}) int {
	if tags == nil {
		return 0
	}

	// Get the "version" tag value
	versionValue, ok := (*tags)["version"]
	if !ok {
		return 0
	}

	// Try to convert to int directly if it's already a number
	switch v := versionValue.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		// Try to parse string as integer
		if version, err := strconv.Atoi(v); err == nil {
			return version
		}
	}

	return 0
}
