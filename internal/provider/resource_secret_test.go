package cpln

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneSecret_basic performs an acceptance test for the resource.
func TestAccControlPlaneSecret_basic(t *testing.T) {
	// Initialize the test
	resourceTest := NewSecretResourceTest()

	// Run the acceptance test case for the resource, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "SECRET") },
		ProtoV6ProviderFactories: GetProviderServer(),
		CheckDestroy:             resourceTest.CheckDestroy,
		Steps:                    resourceTest.InitializeSteps(),
	})
}

/*** Resource Test ***/

// SecretResourceTest defines the necessary functionality to test the resource.
type SecretResourceTest struct {
	Cases *[]SecretResourceTestCase
}

// NewSecretResourceTest creates a SecretResourceTest with initialized test cases for all supported secret types.
func NewSecretResourceTest() SecretResourceTest {
	// Create a resource test instance
	resourceTest := SecretResourceTest{}

	// Initialize the test cases to cover all secret types
	cases := []SecretResourceTestCase{
		resourceTest.NewOpaqueScenario(),
		resourceTest.NewTlsScenario(),
		resourceTest.NewGcpScenario(),
		resourceTest.NewAwsScenario(),
		resourceTest.NewEcrScenario(),
		resourceTest.NewDockerScenario(),
		resourceTest.NewUserpassScenario(),
		resourceTest.NewKeypairScenario(),
		resourceTest.NewDictionaryScenario(),
		resourceTest.NewAzureSdkScenario(),
		resourceTest.NewAzureConnectorScenario(),
		resourceTest.NewNatsAccountScenario(),
	}

	// Set the cases for the resource test
	resourceTest.Cases = &cases

	// Return the resource test
	return resourceTest
}

// InitializeSteps defines the ordered test steps for creating, updating, importing, and reverting secret resources.
func (srt *SecretResourceTest) InitializeSteps() []resource.TestStep {
	// Build combined HCL configs for initial state
	initialConfig := srt.CollectHcl(func(s SecretResourceTestCase) string { return s.InitialHCL })

	// Build combined HCL configs for updated state
	updateConfig := srt.CollectHcl(func(s SecretResourceTestCase) string { return s.UpdateHCL })

	// Collect all test check functions for initial verification
	initialChecks := srt.CollectChecks(func(s SecretResourceTestCase) []resource.TestCheckFunc { return s.InitialChecks })

	// Collect all test check functions for updated verification
	updateChecks := srt.CollectChecks(func(s SecretResourceTestCase) []resource.TestCheckFunc { return s.UpdateChecks })

	// Build import steps for each scenario
	importSteps := make([]resource.TestStep, len(*srt.Cases))
	for i, c := range *srt.Cases {
		// Create an import test step for the current case
		importSteps[i] = resource.TestStep{
			ResourceName: c.Scenario.ResourceAddress,
			ImportState:  true,
		}
	}

	// Declare the full set of test steps
	var steps []resource.TestStep

	// Add step to create the resource and validate initial state
	steps = append(steps, resource.TestStep{
		Config: initialConfig,
		Check:  resource.ComposeAggregateTestCheckFunc(initialChecks...),
	})

	// Add import steps to validate import functionality
	steps = append(steps, importSteps...)

	// Add step to update the resource and validate updated state
	steps = append(steps, resource.TestStep{
		Config: updateConfig,
		Check:  resource.ComposeAggregateTestCheckFunc(updateChecks...),
	})

	// Add step to revert the resource to its initial state
	steps = append(steps, resource.TestStep{
		Config: initialConfig,
		Check:  resource.ComposeAggregateTestCheckFunc(initialChecks...),
	})

	// Return the full sequence of test steps
	return steps
}

// CollectHcl concatenates HCL blocks from all scenarios.
func (srt *SecretResourceTest) CollectHcl(hclFunc func(SecretResourceTestCase) string) string {
	// Initialize a slice to store HCL blocks
	parts := []string{}

	// Iterate over the test cases and generate HCL for each
	for _, c := range *srt.Cases {
		// Append the HCL generated from the current case
		parts = append(parts, hclFunc(c))
	}

	// Join all HCL blocks with double newline separator and return the result
	return strings.Join(parts, "\n\n")
}

// CollectChecks flattens all TestCheckFuncs.
func (srt *SecretResourceTest) CollectChecks(checksFunc func(SecretResourceTestCase) []resource.TestCheckFunc) []resource.TestCheckFunc {
	// Initialize a slice to collect all test check functions
	var checks []resource.TestCheckFunc

	// Iterate over the test cases and extract checks from each
	for _, c := range *srt.Cases {
		// Append all checks from the current case to the result slice
		checks = append(checks, checksFunc(c)...)
	}

	// Return the full list of test check functions
	return checks
}

// CheckDestroy verifies that all resources have been destroyed.
func (srt *SecretResourceTest) CheckDestroy(s *terraform.State) error {
	// Log the start of the destroy check with the count of resources in the root module
	tflog.Info(TestLoggerContext, fmt.Sprintf("Starting CheckDestroy for cpln_secret resources. Total resources: %d", len(s.RootModule().Resources)))

	// If no resources are present in the Terraform state, log and return early
	if len(s.RootModule().Resources) == 0 {
		return errors.New("CheckDestroy error: no resources found in the state to verify")
	}

	// Iterate through each resource in the state
	for _, rs := range s.RootModule().Resources {
		// Log the resource type being checked
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking resource type: %s", rs.Type))

		// Continue only if the resource is as expected
		if rs.Type != "cpln_secret" {
			continue
		}

		// Retrieve the name for the current resource
		secretName := rs.Primary.ID
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of secret with name: %s", secretName))

		// Use the TestProvider client to check if the API resource still exists in the data service
		secret, code, err := TestProvider.client.GetSecret(secretName)

		// If a 404 status code is returned, it indicates the API resource was deleted
		if code == 404 {
			continue
		}

		// If an error occurs during the request, return an error
		if err != nil {
			return fmt.Errorf("error occurred while checking if secret %s exists: %w", secretName, err)
		}

		// If the API resource is found, return an error indicating it still exists
		if secret != nil {
			return fmt.Errorf("CheckDestroy failed: secret %s still exists in the system", *secret.Name)
		}
	}

	// Log successful completion of the destroy check
	tflog.Info(TestLoggerContext, "All cpln_secret resources have been successfully destroyed")
	return nil
}

// Test Cases //

// NewOpaqueScenario creates a test case for opaque secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewOpaqueScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-opaque-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	payload := "opaque_secret_payload"
	encoding := "plain"

	// Define the updated config
	payloadUpdate := "b3BhcXVlX3NlY3JldF9wYXlsb2FkX3VwZGF0ZQ=="
	encodingUpdate := "base64"

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.opaque",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "opaque description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.OpaqueRequiredOnly(payload),
		UpdateHCL:  scenario.OpaqueUpdateWithOptionals(payloadUpdate, encodingUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			scenario.TestCheckNestedBlocks("opaque", []map[string]interface{}{
				{
					"payload":  payload,
					"encoding": encoding,
				},
			}),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			scenario.TestCheckNestedBlocks("opaque", []map[string]interface{}{
				{
					"payload":  payloadUpdate,
					"encoding": encodingUpdate,
				},
			}),
		},
	}
}

// NewTlsScenario creates a test case for TLS secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewTlsScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-tls-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	privateKey := MustLoadTestData("private_key.pem")
	certificate := MustLoadTestData("certificate.pem")

	// Define the updated config
	privateKeyUpdate := MustLoadTestData("private_key_update.pem")
	certificateUpdate := MustLoadTestData("certificate_update.pem")

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.tls",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "tls description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.TlsRequiredOnly(privateKey, certificate),
		UpdateHCL:  scenario.TlsUpdateWithOptionals(privateKeyUpdate, certificateUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			scenario.TestCheckNestedBlocks("tls", []map[string]interface{}{
				{
					"key":  privateKey,
					"cert": certificate,
				},
			}),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			scenario.TestCheckNestedBlocks("tls", []map[string]interface{}{
				{
					"key":   privateKeyUpdate,
					"cert":  certificateUpdate,
					"chain": certificateUpdate,
				},
			}),
		},
	}
}

// NewGcpScenario creates a test case for GCP secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewGcpScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-gcp-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	data := "{   \"type\":   \"gcp\",\"project_id\":\"cpln12345\",\"private_key_id\":\"pvt_key\",\"private_key\":\"key\",\"client_email\":\"support@cpln.io\",\"client_id\":\"12744\",\"auth_uri\":\"cloud.google.com\",\"token_uri\":\"token.cloud.google.com\",\"auth_provider_x509_cert_url\":\"cert.google.com\",\"client_x509_cert_url\":\"cert.google.com\"}"

	// Define the updated config
	dataUpdate := "{   \"type\":   \"gcp\",\"project_id\":\"cpln12345-update\",\"private_key_id\":\"pvt_key\",\"private_key\":\"key\",\"client_email\":\"support@cpln.io\",\"client_id\":\"12744\",\"auth_uri\":\"cloud.google.com\",\"token_uri\":\"token.cloud.google.com\",\"auth_provider_x509_cert_url\":\"cert.google.com\",\"client_x509_cert_url\":\"cert.google.com\"}"

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.gcp",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "gcp description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.GcpRequiredOnly(data),
		UpdateHCL:  scenario.GcpUpdateWithOptionals(dataUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			resource.TestCheckResourceAttr(scenario.ResourceAddress, "gcp", data),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			resource.TestCheckResourceAttr(scenario.ResourceAddress, "gcp", dataUpdate),
		},
	}
}

// NewAwsScenario creates a test case for AWS secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewAwsScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-aws-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	secretKey := "AKIAIOSFODNN7EXAMPLE"
	accessKey := "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	// Define the updated config
	secretKeyUpdate := "AKIAIOSFODNN7EXAMPLEUPDATE"
	accessKeyUpdate := "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEYUPDATE"
	roleArnUpdate := "arn:awskeyupdate"
	externalIdUpdate := "ExampleExternalID-2024-02-09-abc123XYZ-update"

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.aws",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "aws description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.AwsRequiredOnly(secretKey, accessKey),
		UpdateHCL:  scenario.AwsUpdateWithOptionals(secretKeyUpdate, accessKeyUpdate, roleArnUpdate, externalIdUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			scenario.TestCheckNestedBlocks("aws", []map[string]interface{}{
				{
					"secret_key": secretKey,
					"access_key": accessKey,
				},
			}),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			scenario.TestCheckNestedBlocks("aws", []map[string]interface{}{
				{
					"secret_key":  secretKeyUpdate,
					"access_key":  accessKeyUpdate,
					"role_arn":    roleArnUpdate,
					"external_id": externalIdUpdate,
				},
			}),
		},
	}
}

// NewEcrScenario creates a test case for ECR secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewEcrScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-ecr-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	secretKey := "AKIAIOSFODNN7EXAMPLE"
	accessKey := "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	repos := []string{"915716931765.dkr.ecr.us-west-2.amazonaws.com/env-test", "015716931765.dkr.ecr.us-west-2.amazonaws.com/cpln-test"}

	// Define the updated config
	secretKeyUpdate := "AKIAIOSFODNN7EXAMPLEUPDATE"
	accessKeyUpdate := "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEYUPDATE"
	roleArnUpdate := "arn:awskeyupdate"
	externalIdUpdate := "ExampleExternalID-2024-02-09-abc123XYZ-update"
	reposUpdate := []string{"915716931765.dkr.ecr.us-west-2.amazonaws.com/env-test-update", "015716931765.dkr.ecr.us-west-2.amazonaws.com/cpln-test-update", "015716931765.dkr.ecr.us-west-2.amazonaws.com/cpln-test-new"}

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.ecr",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "ecr description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.EcrRequiredOnly(secretKey, accessKey, repos),
		UpdateHCL:  scenario.EcrUpdateWithOptionals(secretKeyUpdate, accessKeyUpdate, roleArnUpdate, externalIdUpdate, reposUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			scenario.TestCheckNestedBlocks("ecr", []map[string]interface{}{
				{
					"secret_key": secretKey,
					"access_key": accessKey,
					"repos":      repos,
				},
			}),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			scenario.TestCheckNestedBlocks("ecr", []map[string]interface{}{
				{
					"secret_key":  secretKeyUpdate,
					"access_key":  accessKeyUpdate,
					"role_arn":    roleArnUpdate,
					"external_id": externalIdUpdate,
					"repos":       reposUpdate,
				},
			}),
		},
	}
}

// NewDockerScenario creates a test case for Docker secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewDockerScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-docker-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	data := "{\"auths\":{\"your-registry-server\":{\"username\":\"your-name\",\"password\":\"your-pword\",\"email\":\"your-email\",\"auth\":\"<Secret>\"}  }  }"

	// Define the updated config
	dataUpdate := "{\"auths\":{\"your-registry-server-update\":{\"username\":\"your-name-update\",\"password\":\"your-pword-update\",\"email\":\"your-email-update\",\"auth\":\"<Secret>\"}  }  }"

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.docker",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "docker description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.DockerRequiredOnly(data),
		UpdateHCL:  scenario.DockerUpdateWithOptionals(dataUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			resource.TestCheckResourceAttr(scenario.ResourceAddress, "docker", data),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			resource.TestCheckResourceAttr(scenario.ResourceAddress, "docker", dataUpdate),
		},
	}
}

// NewUserpassScenario creates a test case for UserPass secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewUserpassScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-userpass-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	username := "cpln_username"
	password := "cpln_password"
	encoding := "plain"

	// Define the updated config
	usernameUpdate := "cpln_username_update"
	passwordUpdate := "cpln_password_update"
	encodingUpdate := "base64"

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.userpass",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "userpass description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.UserpassRequiredOnly(username, password),
		UpdateHCL:  scenario.UserpassUpdateWithOptionals(usernameUpdate, passwordUpdate, encodingUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			scenario.TestCheckNestedBlocks("userpass", []map[string]interface{}{
				{
					"username": username,
					"password": password,
					"encoding": encoding,
				},
			}),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			scenario.TestCheckNestedBlocks("userpass", []map[string]interface{}{
				{
					"username": usernameUpdate,
					"password": passwordUpdate,
					"encoding": encodingUpdate,
				},
			}),
		},
	}
}

// NewKeypairScenario creates a test case for KeyPair secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewKeypairScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-keypair-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	secretKey := MustLoadTestData("secret_key.pem")

	// Define the updated config
	secretKeyUpdate := MustLoadTestData("secret_key_update.pem")
	publicKeyUpdate := MustLoadTestData("public_key.pem")
	passphraseUpdate := "cpln"

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.keypair",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "keypair description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.KeypairRequiredOnly(secretKey),
		UpdateHCL:  scenario.KeypairUpdateWithOptionals(secretKeyUpdate, publicKeyUpdate, passphraseUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			scenario.TestCheckNestedBlocks("keypair", []map[string]interface{}{
				{
					"secret_key": secretKey,
				},
			}),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			scenario.TestCheckNestedBlocks("keypair", []map[string]interface{}{
				{
					"secret_key": secretKeyUpdate,
					"public_key": publicKeyUpdate,
					"passphrase": passphraseUpdate,
				},
			}),
		},
	}
}

// NewDictionaryScenario creates a test case for Dictionary secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewDictionaryScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-dictionary-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	dictionary := map[string]interface{}{
		"key01": "value-01",
		"key02": "value-02",
	}

	// Define the updated config
	dictionaryUpdate := map[string]interface{}{
		"key01":       "value-01",
		"key02update": "value-02-update",
		"key03":       "value-03",
	}

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.dictionary",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "dictionary description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.DictionaryRequiredOnly(dictionary),
		UpdateHCL:  scenario.DictionaryUpdateWithOptionals(dictionaryUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			scenario.TestCheckMapAttr("dictionary", ConvertMapToStringMap(dictionary)),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			scenario.TestCheckMapAttr("dictionary", ConvertMapToStringMap(dictionaryUpdate)),
		},
	}
}

// NewAzureSdkScenario creates a test case for AzureSdk secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewAzureSdkScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-azure-sdk-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	data := "{     \"subscriptionId\":   \"2cd8674e-4f89-4a1f-b420-7a1361b46ef7\",\"tenantId\":\"292f5674-c8b0-488b-9ff8-6d30d77f38d9\",\"clientId\":\"649846ce-d862-49d5-a5eb-7d5aad90f54e\",\"clientSecret\":\"cpln\"}"

	// Define the updated config
	dataUpdate := "{\"subscriptionId\": \"2cd8674e-4f89-4a1f-b420-7a1361b46ef7\", \"tenantId\": \"292f5674-c8b0-488b-9ff8-6d30d77f38d9\", \"clientId\": \"649846ce-d862-49d5-a5eb-7d5aad90f54e\", \"clientSecret\": \"cplnupdate\"}"

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.azure-sdk",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "azure-sdk description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.AzureSdkRequiredOnly(data),
		UpdateHCL:  scenario.AzureSdkUpdateWithOptionals(dataUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			resource.TestCheckResourceAttr(scenario.ResourceAddress, "azure_sdk", data),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			resource.TestCheckResourceAttr(scenario.ResourceAddress, "azure_sdk", dataUpdate),
		},
	}
}

// NewAzureConnectorScenario creates a test case for AzureConnector secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewAzureConnectorScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-azure-connector-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	url := "https://example.azurewebsites.net/api/iam-broker"
	code := "iH0wQjWdAai3oE1C7XrC3t1BBaD7N7foapAylbMaR7HXOmGNYzM3QA=="

	// Define the updated config
	urlUpdate := "https://example.azurewebsites.net/api/iam-broker-update"
	codeUpdate := "iH0wQjWdAai3oE1C7XrC3t1BBaD7N7foapAylbMaR7HXOmGUPDATE=="

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.azure-connector",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "azure-connector description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.AzureConnectorRequiredOnly(url, code),
		UpdateHCL:  scenario.AzureConnectorUpdateWithOptionals(urlUpdate, codeUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			scenario.TestCheckNestedBlocks("azure_connector", []map[string]interface{}{
				{
					"url":  url,
					"code": code,
				},
			}),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			scenario.TestCheckNestedBlocks("azure_connector", []map[string]interface{}{
				{
					"url":  urlUpdate,
					"code": codeUpdate,
				},
			}),
		},
	}
}

// NewNatsAccountScenario creates a test case for NatsAccount secret type with initial and updated configurations.
func (srt *SecretResourceTest) NewNatsAccountScenario() SecretResourceTestCase {
	// Generate a unique name for the secret resource
	name := fmt.Sprintf("secret-nats-account-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	// Define the initial config
	accountId := "AB7JJPKAYKNQOKRKIOS5UCCLALTUAAXCC7FR2QGC4V5UFCAKW4EBIFVZ"
	privateKey := "SAABRA7OGVHKARDQLUQ6THIABW5PMOHJVPSOPTWZRP4WD5LPVOLGTU6ONQ"

	// Define the updated config
	accountIdUpdate := "AB7JJPKAYKNQOKRKIOS5UCCLALTUAAXCC7FR2QGC4V5UFCAKW4EBIFVZ"
	privateKeyUpdate := "SAABRA7OGVHKARDQLUQ6THIABW5PMOHJVPSOPTWZRP4WD5LPVOLGTU6ONQ"

	// Create the secret test scenario with metadata and descriptions
	scenario := SecretResourceTestScenario{
		ProviderTestCase: ProviderTestCase{
			Kind:              "secret",
			ResourceAddress:   "cpln_secret.nats-account",
			Name:              name,
			Description:       name,
			DescriptionUpdate: "nats-account description updated",
		},
	}

	// Return the complete test case for the opaque secret
	return SecretResourceTestCase{
		Scenario:   scenario,
		InitialHCL: scenario.NatsAccountRequiredOnly(accountId, privateKey),
		UpdateHCL:  scenario.NatsAccountUpdateWithOptionals(accountIdUpdate, privateKeyUpdate),
		InitialChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.Description, "0"),
			scenario.TestCheckNestedBlocks("nats_account", []map[string]interface{}{
				{
					"account_id":  accountId,
					"private_key": privateKey,
				},
			}),
		},
		UpdateChecks: []resource.TestCheckFunc{
			scenario.Exists(),
			scenario.GetDefaultChecks(scenario.DescriptionUpdate, "3"),
			scenario.TestCheckNestedBlocks("nats_account", []map[string]interface{}{
				{
					"account_id":  accountIdUpdate,
					"private_key": privateKeyUpdate,
				},
			}),
		},
	}
}

/*** Resource Test Case ***/

// SecretResourceTestScenario defines a specific resource test case.
type SecretResourceTestCase struct {
	Scenario      SecretResourceTestScenario
	InitialHCL    string
	UpdateHCL     string
	InitialChecks []resource.TestCheckFunc
	UpdateChecks  []resource.TestCheckFunc
}

/*** Resource Test Scenario ***/

// SecretResourceTestScenario defines a specific resource test scenario.
type SecretResourceTestScenario struct {
	ProviderTestCase
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (srts *SecretResourceTestScenario) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Log the start of the existence check with the resource count
		tflog.Info(TestLoggerContext, fmt.Sprintf("Checking existence of secret: %s. Total resources: %d", srts.Name, len(s.RootModule().Resources)))

		// Retrieve the resource from the Terraform state
		rs, ok := s.RootModule().Resources[srts.ResourceAddress]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", srts.ResourceAddress)
		}

		// Ensure the resource ID matches the expected API resource name
		if rs.Primary.ID != srts.Name {
			return fmt.Errorf("resource ID %s does not match expected secret name %s", rs.Primary.ID, srts.Name)
		}

		// Retrieve the API resource from the external system using the provider client
		remoteSecret, _, err := TestProvider.client.GetSecret(srts.Name)
		if err != nil {
			return fmt.Errorf("error retrieving secret from external system: %w", err)
		}

		// Verify the API resource name from the external system matches the expected resource name
		if *remoteSecret.Name != srts.Name {
			return fmt.Errorf("mismatch in secret name: expected %s, got %s", srts.Name, *remoteSecret.Name)
		}

		// Log successful verification of API resource existence
		tflog.Info(TestLoggerContext, fmt.Sprintf("secret %s verified successfully in both state and external system.", srts.Name))
		return nil
	}
}

// Configs //

// OpaqueRequiredOnly returns a minimal HCL block for an opaque secret using only required fields.
func (srts *SecretResourceTestScenario) OpaqueRequiredOnly(payload string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "opaque" {
  name = "%s"

  opaque {
    payload  = "%s"
  }
}
`, srts.Name, payload)
}

// OpaqueUpdateWithOptionals returns an HCL block for an opaque secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) OpaqueUpdateWithOptionals(payload string, encoding string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "opaque" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
  }

  opaque {
    payload  = "%s"
    encoding = "%s"
  }
}
`, srts.Name, srts.DescriptionUpdate, payload, encoding)
}

// TlsRequiredOnly returns a minimal HCL block for a TLS secret using only required fields.
func (srts *SecretResourceTestScenario) TlsRequiredOnly(key string, cert string) string {
	return fmt.Sprintf(`
variable "testcertprivate" {
  type = string
  default = <<EOT
%s
EOT
}

variable "testcert" {
  type = string
  default = <<EOT
%s
EOT
}

resource "cpln_secret" "tls" {
  name = "%s"

  tls {
    key   = chomp(var.testcertprivate)
    cert  = chomp(var.testcert)
  }
}
`, key, cert, srts.Name)
}

// TlsUpdateWithOptionals returns an HCL block for a TLS secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) TlsUpdateWithOptionals(key string, cert string) string {
	return fmt.Sprintf(`
variable "testcertprivate" {
  type = string
  default = <<EOT
%s
EOT
}

variable "testcert" {
  type = string
  default = <<EOT
%s
EOT
}

resource "cpln_secret" "tls" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "tls"
  }

  tls {
    key   = chomp(var.testcertprivate)
    cert  = chomp(var.testcert)
    chain = chomp(var.testcert)
  }
}
`, key, cert, srts.Name, srts.DescriptionUpdate)
}

// GcpRequiredOnly returns a minimal HCL block for a GCP secret using only required fields.
func (srts *SecretResourceTestScenario) GcpRequiredOnly(data string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "gcp" {
  name = "%s"
  gcp  = trimspace(<<EOT
%s
EOT
  )
}
`, srts.Name, data)
}

// GcpUpdateWithOptionals returns an HCL block for a GCP secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) GcpUpdateWithOptionals(data string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "gcp" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "gcp"
  }

  gcp = trimspace(<<EOT
%s
EOT
  )
}
`, srts.Name, srts.DescriptionUpdate, data)
}

// AwsRequiredOnly returns a minimal HCL block for a AWS secret using only required fields.
func (srts *SecretResourceTestScenario) AwsRequiredOnly(secretKey string, accessKey string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name = "%s"

  aws {
    secret_key  = "%s"
    access_key  = "%s"
  }
}
`, srts.Name, secretKey, accessKey)
}

// AwsUpdateWithOptionals returns an HCL block for a AWS secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) AwsUpdateWithOptionals(secretKey string, accessKey string, roleArn string, externalId string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "aws" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "aws"
  }
	
  aws {
    secret_key  = "%s"
    access_key  = "%s"
    role_arn    = "%s"
    external_id = "%s"
  }
}
`, srts.Name, srts.DescriptionUpdate, secretKey, accessKey, roleArn, externalId)
}

// EcrRequiredOnly returns a minimal HCL block for a ECR secret using only required fields.
func (srts *SecretResourceTestScenario) EcrRequiredOnly(secretKey string, accessKey string, repos []string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "ecr" {
  name   = "%s"

  ecr {
    secret_key  = "%s"
    access_key  = "%s"
    repos       = %s
  }
}
`, srts.Name, secretKey, accessKey, StringSliceToString(repos))
}

// EcrUpdateWithOptionals returns an HCL block for a ECR secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) EcrUpdateWithOptionals(secretKey string, accessKey string, roleArn string, externalId string, repos []string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "ecr" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "ecr"
	}

  ecr {
    secret_key  = "%s"
    access_key  = "%s"
    role_arn    = "%s"
    external_id = "%s"
    repos       = %s
	}
}
`, srts.Name, srts.DescriptionUpdate, secretKey, accessKey, roleArn, externalId, StringSliceToString(repos))
}

// DockerRequiredOnly returns a minimal HCL block for a Docker secret using only required fields.
func (srts *SecretResourceTestScenario) DockerRequiredOnly(data string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "docker" {
  name   = "%s"
  docker = trimspace(<<EOT
%s
EOT
  )
}
`, srts.Name, data)
}

// DockerUpdateWithOptionals returns an HCL block for a Docker secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) DockerUpdateWithOptionals(data string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "docker" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "docker"
  }

  docker = trimspace(<<EOT
%s
EOT
  )
}
`, srts.Name, srts.DescriptionUpdate, data)
}

// UserpassRequiredOnly returns a minimal HCL block for a UserPass secret using only required fields.
func (srts *SecretResourceTestScenario) UserpassRequiredOnly(username string, password string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "userpass" {
  name = "%s"

  userpass {
    username = "%s"
    password = "%s"
  }
}
`, srts.Name, username, password)
}

// UserpassUpdateWithOptionals returns an HCL block for a UserPass secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) UserpassUpdateWithOptionals(username string, password string, encoding string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "userpass" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "userpass"
  }

  userpass {
    username = "%s"
    password = "%s"
    encoding = "%s"
  }
}
`, srts.Name, srts.DescriptionUpdate, username, password, encoding)
}

// KeypairRequiredOnly returns a minimal HCL block for a KeyPair secret using only required fields.
func (srts *SecretResourceTestScenario) KeypairRequiredOnly(data string) string {
	return fmt.Sprintf(`
variable "test-secret-key" {
  type = string
  default = <<EOT
%s
EOT
}

resource "cpln_secret" "keypair" {
  name = "%s"

  keypair {
    secret_key = chomp(var.test-secret-key)
  }
}
`, data, srts.Name)
}

// KeypairUpdateWithOptionals returns an HCL block for a UserPass secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) KeypairUpdateWithOptionals(secretKey string, publicKey string, passphrase string) string {
	return fmt.Sprintf(`
variable "test-secret-key" {
  type = string
  default = <<EOT
%s
EOT
}

variable "test-public-key" {
  type = string
  default = <<EOT
%s
EOT
}

resource "cpln_secret" "keypair" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "keypair"
  }

  keypair {
    secret_key = chomp(var.test-secret-key)
    public_key = chomp(var.test-public-key)
    passphrase = "%s"
  }
}
`, secretKey, publicKey, srts.Name, srts.DescriptionUpdate, passphrase)
}

// DictionaryRequiredOnly returns a minimal HCL block for a Dictionary secret using only required fields.
func (srts *SecretResourceTestScenario) DictionaryRequiredOnly(dictionary map[string]interface{}) string {
	return fmt.Sprintf(`
resource "cpln_secret" "dictionary" {
  name   = "%s"

  dictionary = %s
}
`, srts.Name, MapToHCL(dictionary, 2))
}

// DictionaryUpdateWithOptionals returns an HCL block for a Dictionary secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) DictionaryUpdateWithOptionals(dictionary map[string]interface{}) string {
	return fmt.Sprintf(`
resource "cpln_secret" "dictionary" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "dictionary"
  }

  dictionary = %s
}
`, srts.Name, srts.DescriptionUpdate, MapToHCL(dictionary, 2))
}

// AzureSdkRequiredOnly returns a minimal HCL block for a AzureSdk secret using only required fields.
func (srts *SecretResourceTestScenario) AzureSdkRequiredOnly(data string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "azure-sdk" {
  name      = "%s"
  azure_sdk = trimspace(<<EOT
%s
EOT
  )
}
`, srts.Name, data)
}

// AzureSdkUpdateWithOptionals returns an HCL block for a AzureSdk secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) AzureSdkUpdateWithOptionals(data string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "azure-sdk" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "azure-sdk"
  }

  azure_sdk = trimspace(<<EOT
%s
EOT
  )
}
`, srts.Name, srts.DescriptionUpdate, data)
}

// AzureConnectorRequiredOnly returns a minimal HCL block for a AzureConnector secret using only required fields.
func (srts *SecretResourceTestScenario) AzureConnectorRequiredOnly(url string, code string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "azure-connector" {
  name = "%s"

  azure_connector {
    url  = "%s"
    code = "%s"
  }
}
`, srts.Name, url, code)
}

// AzureConnectorUpdateWithOptionals returns an HCL block for a AzureConnector secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) AzureConnectorUpdateWithOptionals(url string, code string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "azure-connector" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "azure-connector"
  }

  azure_connector {
    url  = "%s"
    code = "%s"
  }
}
`, srts.Name, srts.DescriptionUpdate, url, code)
}

// NatsAccountRequiredOnly returns a minimal HCL block for a NatsAccount secret using only required fields.
func (srts *SecretResourceTestScenario) NatsAccountRequiredOnly(accountId string, privateKey string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "nats-account" {
  name = "%s"

  nats_account {
    account_id  = "%s"
    private_key = "%s"
  }
}
`, srts.Name, accountId, privateKey)
}

// NatsAccountUpdateWithOptionals returns an HCL block for a NatsAccount secret including optional fields like description and tags.
func (srts *SecretResourceTestScenario) NatsAccountUpdateWithOptionals(accountId string, privateKey string) string {
	return fmt.Sprintf(`
resource "cpln_secret" "nats-account" {
  name        = "%s"
  description = "%s"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "nats-account"
  }

  nats_account {
    account_id  = "%s"
    private_key = "%s"
  }
}
`, srts.Name, srts.DescriptionUpdate, accountId, privateKey)
}
