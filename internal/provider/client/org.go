package cpln

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Org - Organization
type Org struct {
	Base
	Spec   *OrgSpec   `json:"spec,omitempty"`
	Status *OrgStatus `json:"status,omitempty"`
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

// Logging - Logging
type Logging struct {
	S3        *S3Logging        `json:"s3,omitempty"`
	Coralogix *CoralogixLogging `json:"coralogix,omitempty"`
	Datadog   *DatadogLogging   `json:"datadog,omitempty"`
	Logzio    *LogzioLogging    `json:"logzio,omitempty"`
	Elastic   *ElasticLogging   `json:"elastic,omitempty"`
}

// OrgSpec - Organization Spec
type OrgSpec struct {
	Logging      *Logging   `json:"logging,omitempty"`
	ExtraLogging *[]Logging `json:"extraLogging,omitempty"`
	Tracing      *Tracing   `json:"tracing,omitempty"`
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

// GetOrg - Get Organziation By Name
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
