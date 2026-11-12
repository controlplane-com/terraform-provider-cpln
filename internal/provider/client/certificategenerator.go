package cpln

import "fmt"

type CertificateGenerator struct {
	Base
	Spec        *CertificateGeneratorSpec   `json:"spec,omitempty"`
	SpecReplace *CertificateGeneratorSpec   `json:"$replace/spec,omitempty"`
	Status      *CertificateGeneratorStatus `json:"status,omitempty"`
}

/*** Spec ***/

type CertificateGeneratorSpec struct {
	Provider *string                   `json:"provider,omitempty"` // Enum: "acme"
	Acme     *CertificateGeneratorAcme `json:"acme,omitempty"`
}

/*** Spec Related ***/

type CertificateGeneratorAcme struct {
	Server                 *string                                     `json:"server,omitempty"`
	Email                  *string                                     `json:"email,omitempty"`
	ExternalAccountBinding *CertificateGeneratorExternalAccountBinding `json:"externalAccountBinding,omitempty"`
	Challenge              *CertificateGeneratorAcmeChallenge          `json:"challenge,omitempty"`
}

type CertificateGeneratorExternalAccountBinding struct {
	KeyId      *string `json:"keyId,omitempty"`
	SecretLink *string `json:"secretLink,omitempty"`
}

type CertificateGeneratorAcmeChallenge struct {
	Http01 *CertificateGeneratorHttp01Challenge `json:"http01,omitempty"`
	Dns01  *CertificateGeneratorDns01Challenge  `json:"dns01,omitempty"`
}

type CertificateGeneratorHttp01Challenge struct{}

type CertificateGeneratorDns01Challenge struct {
	Provider   *string                              `json:"provider,omitempty"` // Enum: "cpln", "route53", "cloudDns", "azureDns", "cloudflare"
	Route53    *CertificateGeneratorRoute53Dns01    `json:"route53,omitempty"`
	CloudDns   *CertificateGeneratorCloudDnsDns01   `json:"cloudDns,omitempty"`
	AzureDns   *CertificateGeneratorAzureDnsDns01   `json:"azureDns,omitempty"`
	Cloudflare *CertificateGeneratorCloudflareDns01 `json:"cloudflare,omitempty"`
}

type CertificateGeneratorRoute53Dns01 struct {
	Region       *string `json:"region,omitempty"`
	HostedZoneId *string `json:"hostedZoneId,omitempty"`
	SecretLink   *string `json:"secretLink,omitempty"`
}

type CertificateGeneratorCloudDnsDns01 struct {
	Project    *string `json:"project,omitempty"`
	SecretLink *string `json:"secretLink,omitempty"`
}

type CertificateGeneratorAzureDnsDns01 struct {
	ResourceGroupName *string `json:"resourceGroupName,omitempty"`
	HostedZoneName    *string `json:"hostedZoneName,omitempty"`
	SecretLink        *string `json:"secretLink,omitempty"`
}

type CertificateGeneratorCloudflareDns01 struct {
	SecretLink *string `json:"secretLink,omitempty"`
}

/*** Status ***/

type CertificateGeneratorStatus struct {
	Ready       *bool   `json:"ready,omitempty"`
	Message     *string `json:"message,omitempty"`
	LastUpdated *string `json:"lastUpdated,omitempty"`
}

/*** Client Functions ***/

// GetCertificateGenerator - Get Certificate Generator by name
func (c *Client) GetCertificateGenerator(name string) (*CertificateGenerator, int, error) {

	certificateGenerator, code, err := c.GetResource(fmt.Sprintf("certificategenerator/%s", name), new(CertificateGenerator))

	if err != nil {
		return nil, code, err
	}

	return certificateGenerator.(*CertificateGenerator), code, err
}

// CreateCertificateGenerator - Create a new Certificate Generator
func (c *Client) CreateCertificateGenerator(certificateGenerator CertificateGenerator) (*CertificateGenerator, int, error) {

	code, err := c.CreateResource("certificategenerator", *certificateGenerator.Name, certificateGenerator)

	if err != nil {
		return nil, code, err
	}

	return c.GetCertificateGenerator(*certificateGenerator.Name)
}

// UpdateCertificateGenerator - Update an existing Certificate Generator
func (c *Client) UpdateCertificateGenerator(certificateGenerator CertificateGenerator) (*CertificateGenerator, int, error) {

	code, err := c.UpdateResource(fmt.Sprintf("certificategenerator/%s", *certificateGenerator.Name), certificateGenerator)

	if err != nil {
		return nil, code, err
	}

	return c.GetCertificateGenerator(*certificateGenerator.Name)
}

// DeleteCertificateGenerator - Delete Certificate Generator by name
func (c *Client) DeleteCertificateGenerator(name string) error {
	return c.DeleteResource(fmt.Sprintf("certificategenerator/%s", name))
}
