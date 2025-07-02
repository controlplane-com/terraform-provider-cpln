package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Org - Organization
type Org struct {
	Base
	Spec        *OrgSpec   `json:"spec,omitempty"`
	SpecReplace *OrgSpec   `json:"$replace/spec,omitempty"`
	Status      *OrgStatus `json:"status,omitempty"`
}

type CreateOrgRequest struct {
	Org      *Org      `json:"org,omitempty"`
	Invitees *[]string `json:"invitees,omitempty"`
}

// OrgStatus - Organization Status
type OrgStatus struct {
	AccountLink *string `json:"accountLink,omitempty"`
	Active      *bool   `json:"active,omitempty"`
}

// S3Logging - S3Logging
type S3Logging struct {
	Bucket      *string `json:"bucket,omitempty"`
	Region      *string `json:"region,omitempty"`
	Prefix      *string `json:"prefix,omitempty"`
	Credentials *string `json:"credentials,omitempty"`
}

// CoralogixLogging - CoralogixLogging
type CoralogixLogging struct {
	Cluster     *string `json:"cluster,omitempty"`
	Credentials *string `json:"credentials,omitempty"`
	App         *string `json:"app,omitempty"`
	Subsystem   *string `json:"subsystem,omitempty"`
}

// DatadogLogging - DatadogLogging
type DatadogLogging struct {
	Host        *string `json:"host,omitempty"`
	Credentials *string `json:"credentials,omitempty"`
}

// LogzioLogging - LogzioLogging
type LogzioLogging struct {
	ListenerHost *string `json:"listenerHost,omitempty"`
	Credentials  *string `json:"credentials,omitempty"`
}

// ElasticLogging - ElasticLogging
type ElasticLogging struct {
	AWS          *AWSLogging          `json:"aws,omitempty"`
	ElasticCloud *ElasticCloudLogging `json:"elasticCloud,omitempty"`
	Generic      *GenericLogging      `json:"generic,omitempty"`
}

type AWSLogging struct {
	Host        *string `json:"host,omitempty"`
	Port        *int    `json:"port,omitempty"`
	Index       *string `json:"index,omitempty"`
	Type        *string `json:"type,omitempty"`
	Credentials *string `json:"credentials,omitempty"`
	Region      *string `json:"region,omitempty"`
}

type ElasticCloudLogging struct {
	Index       *string `json:"index,omitempty"`
	Type        *string `json:"type,omitempty"`
	Credentials *string `json:"credentials,omitempty"`
	CloudID     *string `json:"cloudId,omitempty"`
}

type GenericLogging struct {
	Host        *string `json:"host,omitempty"`
	Port        *int    `json:"port,omitempty"`
	Path        *string `json:"path,omitempty"`
	Index       *string `json:"index,omitempty"`
	Type        *string `json:"type,omitempty"`
	Credentials *string `json:"credentials,omitempty"`
}

type CloudWatchLogging struct {
	Region        *string                 `json:"region,omitempty"`
	Credentials   *string                 `json:"credentials,omitempty"`
	RetentionDays *int                    `json:"retentionDays,omitempty"`
	GroupName     *string                 `json:"groupName,omitempty"`
	StreamName    *string                 `json:"streamName,omitempty"`
	ExtractFields *map[string]interface{} `json:"extractFields,omitempty"`
}

type FluentdLogging struct {
	Host *string `json:"host,omitempty"`
	Port *int    `json:"port,omitempty"`
}

type StackdriverLogging struct {
	Credentials *string `json:"credentials,omitempty"`
	Location    *string `json:"location,omitempty"`
}

type SyslogLogging struct {
	Host     *string `json:"host,omitempty"`
	Port     *int    `json:"port,omitempty"`
	Mode     *string `json:"mode,omitempty"`
	Format   *string `json:"format,omitempty"`
	Severity *int    `json:"severity,omitempty"`
}

// Logging - Logging
type Logging struct {
	S3          *S3Logging          `json:"s3,omitempty"`
	Coralogix   *CoralogixLogging   `json:"coralogix,omitempty"`
	Datadog     *DatadogLogging     `json:"datadog,omitempty"`
	Logzio      *LogzioLogging      `json:"logzio,omitempty"`
	Elastic     *ElasticLogging     `json:"elastic,omitempty"`
	CloudWatch  *CloudWatchLogging  `json:"cloudWatch,omitempty"`
	Fluentd     *FluentdLogging     `json:"fluentd,omitempty"`
	Stackdriver *StackdriverLogging `json:"stackdriver,omitempty"`
	Syslog      *SyslogLogging      `json:"syslog,omitempty"`
}

// AuthConfig - AuthConfig
type AuthConfig struct {
	DomainAutoMembers *[]string `json:"domainAutoMembers,omitempty"`
	SamlOnly          *bool     `json:"samlOnly,omitempty"`
}

// Observability - Observability
type Observability struct {
	LogsRetentionDays    *int      `json:"logsRetentionDays,omitempty"`
	MetricsRetentionDays *int      `json:"metricsRetentionDays,omitempty"`
	TracesRetentionDays  *int      `json:"tracesRetentionDays,omitempty"`
	DefaultAlertEmails   *[]string `json:"defaultAlertEmails,omitempty"`
}

type OrgThreatDetectionSyslog struct {
	Transport *string `json:"transport,omitempty"`
	Host      *string `json:"host,omitempty"`
	Port      *int    `json:"port,omitempty"`
}

type OrgThreatDetection struct {
	Enabled         *bool                     `json:"enabled,omitempty"`
	MinimumSeverity *string                   `json:"minimumSeverity,omitempty"`
	Syslog          *OrgThreatDetectionSyslog `json:"syslog,omitempty"`
}

type OrgSecurity struct {
	ThreatDetection *OrgThreatDetection `json:"threatDetection,omitempty"`
}

// OrgSpec - Organization Spec
type OrgSpec struct {
	Logging               *Logging       `json:"logging,omitempty"`
	ExtraLogging          *[]Logging     `json:"extraLogging,omitempty"`
	Tracing               *Tracing       `json:"tracing,omitempty"`
	SessionTimeoutSeconds *int           `json:"sessionTimeoutSeconds,omitempty"`
	AuthConfig            *AuthConfig    `json:"authConfig,omitempty"`
	Observability         *Observability `json:"observability,omitempty"`
	Security              *OrgSecurity   `json:"security,omitempty"`
}

type UpdateSpec struct {
	Spec interface{} `json:"spec"`
}

type ReplaceLogging struct {
	Logging      *Logging   `json:"$replace/logging"`
	ExtraLogging *[]Logging `json:"$replace/extraLogging"`
}

type ReplaceTracing struct {
	Tracing *Tracing `json:"$replace/tracing"`
}

// GetOrg - Get Organization By Name
func (c *Client) GetOrg() (*Org, int, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s", c.HostURL, c.Org), nil)

	if err != nil {
		return nil, 0, err
	}

	body, code, err := c.doRequest(req, "")
	if err != nil {
		return nil, code, err
	}

	org := Org{}
	err = json.Unmarshal(body, &org)
	if err != nil {
		return nil, code, err
	}

	return &org, code, nil
}

// GetSpecificOrg - Get Organization By Name
func (c *Client) GetSpecificOrg(name string) (*Org, int, error) {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s", c.HostURL, name), nil)

	if err != nil {
		return nil, 0, err
	}

	body, code, err := c.doRequest(req, "")
	if err != nil {
		return nil, code, err
	}

	org := Org{}
	err = json.Unmarshal(body, &org)
	if err != nil {
		return nil, code, err
	}

	return &org, code, nil
}

func (c *Client) GetOrgAccount(orgName string) (*Account, int, error) {

	billingNgEndpoint, code, err := c.GetBillingNgEndpoint()
	if err != nil {
		return nil, code, err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/org/%s/account", billingNgEndpoint, orgName), nil)

	if err != nil {
		return nil, 0, err
	}

	// Add Bearer prefix if it doesn't exist on the token
	tokenWithBearer := c.Token
	if !strings.HasPrefix(strings.ToLower(c.Token), "bearer ") {
		tokenWithBearer = "Bearer " + c.Token
	}

	body, code, err := c.doRequest(req, "", tokenWithBearer)
	if err != nil {
		return nil, code, err
	}

	account := Account{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, code, err
	}

	return &account, code, nil
}

// CreateOrg - Create Organization
func (c *Client) CreateOrg(accountId string, createOrg CreateOrgRequest) (*Org, int, error) {

	g, err := json.Marshal(createOrg)
	if err != nil {
		return nil, 0, err
	}

	billingNgEndpoint, code, err := c.GetBillingNgEndpoint()
	if err != nil {
		return nil, code, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/account/%s/org", billingNgEndpoint, accountId), strings.NewReader(string(g)))
	if err != nil {
		return nil, 0, err
	}

	// Add Bearer prefix if it doesn't exist on the token
	tokenWithBearer := c.Token
	if !strings.HasPrefix(strings.ToLower(c.Token), "bearer ") {
		tokenWithBearer = "Bearer " + c.Token
	}

	_, code, err = c.doRequest(req, "application/json", tokenWithBearer)
	if err != nil {
		return nil, code, err
	}

	return c.GetSpecificOrg(*createOrg.Org.Name)
}

// UpdateOrg - Update Organization
func (c *Client) UpdateOrg(org Org) (*Org, int, error) {

	g, err := json.Marshal(org)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/org/%s", c.HostURL, c.Org), strings.NewReader(string(g)))
	if err != nil {
		return nil, 0, err
	}

	_, code, err := c.doRequest(req, "application/json")
	if err != nil {
		return nil, code, err
	}

	return c.GetSpecificOrg(c.Org)
}

// UpdateOrgLogging - Update an existing Org Logging
func (c *Client) UpdateOrgLogging(extraLogging *[]Logging) (*Org, int, error) {

	var logging *Logging

	if extraLogging != nil && len(*extraLogging) > 0 {
		logging = &(*extraLogging)[0]
		*extraLogging = (*extraLogging)[1:]
	}

	spec := UpdateSpec{
		Spec: ReplaceLogging{
			Logging:      logging,
			ExtraLogging: extraLogging,
		},
	}

	code, err := c.UpdateResource("", spec)
	if err != nil {
		return nil, code, err
	}

	return c.GetOrg()
}

// UpdateOrgLogging - Update an existing Org Tracing
func (c *Client) UpdateOrgTracing(tracing *Tracing) (*Org, int, error) {

	spec := UpdateSpec{
		Spec: ReplaceTracing{
			Tracing: tracing,
		},
	}

	code, err := c.UpdateResource("", spec)
	if err != nil {
		return nil, code, err
	}

	return c.GetOrg()
}
