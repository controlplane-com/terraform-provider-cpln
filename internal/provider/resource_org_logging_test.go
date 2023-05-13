package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneOrgLogging_basic(t *testing.T) {

	var testLogging client.Logging

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "ORG_LOGGING") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneOrgCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneOrgS3(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging, "s3_logging"),
				),
			},
			{
				Config: testAccControlPlaneOrgCoralogix(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging, "coralogix_logging"),
				),
			},
			{
				Config: testAccControlPlaneOrgDatadog(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging, "datadog_logging"),
				),
			},
			{
				Config: testAccControlPlaneOrgLogzio(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging, "logzio_logging"),
				),
			},
			{
				Config: testAccControlPlaneOrgLogzioWithDifferentListenerHost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging, "logzio_logging-different_listener_host"),
				),
			},
			{
				Config: testAccControlPlaneOrgElasticAWS(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging, "elastic_logging-aws"),
				),
			},
			{
				Config: testAccControlPlaneOrgElasticCloud(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging, "elastic_logging-elastic_cloud"),
				),
			},
		},
	})
}

func testAccControlPlaneOrgS3() string {

	TestLogger.Printf("Inside testAccControlPlaneOrg")

	return `

    resource "cpln_secret" "aws" {
        name = "aws-random-tbd"
        description = "aws description aws-random-tbd" 
                
        tags = {
            terraform_generated = "true"
            acceptance_test = "true"
            secret_type = "aws"
        } 
        
        aws {
            secret_key = "AKIAIOSFODNN7EXAMPLE"
            access_key = "AKIAwJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
            role_arn = "arn:awskey" 
        }
    }

    resource "cpln_org_logging" "tf-logging" {

        s3_logging {

            bucket = "test-bucket"
            region = "us-east1"
            prefix = "/"

            // AWS Secret Only
            credentials = cpln_secret.aws.self_link
        }   
    }
    `
}

func testAccControlPlaneOrgCoralogix() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgCoralogix")

	return `

    resource "cpln_secret" "opaque" {

        name = "opaque-random-coralogix-tbd"
        description = "opaque description opaque-random-tbd" 
        
        tags = {
            terraform_generated = "true"
            acceptance_test = "true"
            secret_type = "opaque"
        }

        opaque {
            payload = "opaque_secret_payload"
            encoding = "plain"
        }
    }

    resource "cpln_org_logging" "tf-logging" {

        coralogix_logging {

            // Valid clusters
            // coralogix.com, coralogix.us, app.coralogix.in, app.eu2.coralogix.com, app.coralogixsg.com
            cluster = "coralogix.com"

            // Opaque Secret Only
            credentials = cpln_secret.opaque.self_link
            
            // Supported variables for App and Subsystem are:
            // {org}, {gvc}, {workload}, {location}
            app = "{workload}"
            subsystem = "{org}"
        }
    }       
    `
}

func testAccControlPlaneOrgDatadog() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgDatadog")

	return `

    resource "cpln_secret" "opaque" {

        name = "opaque-random-datadog-tbd"
        description = "opaque description" 
        
        tags = {
            terraform_generated = "true"
            acceptance_test = "true"
            secret_type = "opaque"
        }

        opaque {
            payload = "opaque_secret_payload"
            encoding = "plain"
        }
    }

    resource "cpln_org_logging" "tf-logging" {

        datadog_logging {

            // Valid Host
            // http-intake.logs.datadoghq.com, http-intake.logs.us3.datadoghq.com, 
            // http-intake.logs.us5.datadoghq.com, http-intake.logs.datadoghq.eu
            host = "http-intake.logs.datadoghq.com"

            // Opaque Secret Only
            credentials = cpln_secret.opaque.self_link  
        }
    }       
    `
}

func testAccControlPlaneOrgLogzio() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgLogzio")

	return `

    resource "cpln_secret" "opaque" {

        name = "opaque-random-datadog-tbd"
        description = "opaque description" 
        
        tags = {
            terraform_generated = "true"
            acceptance_test = "true"
            secret_type = "opaque"
        }

        opaque {
            payload = "opaque_secret_payload"
            encoding = "plain"
        }
    }

    resource "cpln_org_logging" "tf-logging" {

        logzio_logging {

            // Valid Listener Hosts
            // listener.logz.io, listener-nl.logz.io 
            listener_host = "listener.logz.io"

            // Opaque Secret Only
            credentials = cpln_secret.opaque.self_link  
        }
    }       
    `
}

func testAccControlPlaneOrgLogzioWithDifferentListenerHost() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgLogzio1")

	return `

    resource "cpln_secret" "opaque" {

        name = "opaque-random-datadog-tbd"
        description = "opaque description" 
        
        tags = {
            terraform_generated = "true"
            acceptance_test = "true"
            secret_type = "opaque"
        }

        opaque {
            payload = "opaque_secret_payload"
            encoding = "plain"
        }
    }

    resource "cpln_org_logging" "tf-logging" {

        logzio_logging {

            // Valid Hosts
            // listener.logz.io, listener-nl.logz.io 
            listener_host = "listener-nl.logz.io"

            // Opaque Secret Only
            credentials = cpln_secret.opaque.self_link  
        }
    }       
    `
}

func testAccControlPlaneOrgElasticAWS() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgElasticAWS")

	return `

	resource "cpln_secret" "opaque" {

        name = "opaque-random-elastic-logging-aws-tbd"
        description = "opaque description" 
        
        tags = {
            terraform_generated = "true"
            acceptance_test = "true"
            secret_type = "opaque"
        }

        opaque {
            payload = "opaque_secret_payload"
            encoding = "plain"
        }
    }

    resource "cpln_org_logging" "tf-logging" {

        elastic_logging {
			aws {
				host = "es.amazonaws.com"
				port = 8080
				index = "my-index"
				type = "my-type"
				credentials = cpln_secret.opaque.self_link
				region = "us-east-1"
			}
        }
    }
	`
}

func testAccControlPlaneOrgElasticCloud() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgElasticCloud")

	return `

	resource "cpln_secret" "opaque" {

        name = "opaque-random-elastic-logging-elastic-cloud-tbd"
        description = "opaque description" 
        
        tags = {
            terraform_generated = "true"
            acceptance_test = "true"
            secret_type = "opaque"
        }

        opaque {
            payload = "opaque_secret_payload"
            encoding = "plain"
        }
    }

    resource "cpln_org_logging" "tf-logging" {

        elastic_logging {
			elastic_cloud {
				index = "my-index"
				type = "my-type"
				credentials = cpln_secret.opaque.self_link
				cloud_id = "my-cloud-id"
			}
        }
    }
    `
}

func testAccCheckControlPlaneLoggingExists(resourceName string, logging *client.Logging) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		TestLogger.Printf("Inside testAccCheckControlPlaneLoggingExists. Resources Length: %d", len(s.RootModule().Resources))

		_, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Resource not found: %s", s)
		}

		client := testAccProvider.Meta().(*client.Client)
		org, _, err := client.GetOrg()

		if err != nil {
			return err
		}

		*logging = *org.Spec.Logging

		return nil
	}
}

func testAccCheckControlPlaneLoggingAttributes(logging *client.Logging, loggingType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		var expectedValue interface{}
		var toTestValue interface{}

		switch loggingType {
		case "s3_logging":
			value, _, _ := generateTestS3Logging()
			expectedValue = value
			toTestValue = logging.S3

		case "coralogix_logging":
			value, _, _ := generateTestCoralogixLogging()
			expectedValue = value
			toTestValue = logging.Coralogix

		case "datadog_logging":
			value, _, _ := generateTestDatadogLogging()
			expectedValue = value
			toTestValue = logging.Datadog

		case "logzio_logging":
			value, _, _ := generateTestLogzioLogging(loggingType)
			expectedValue = value
			toTestValue = logging.Logzio

		case "logzio_logging-different_listener_host":
			value, _, _ := generateTestLogzioLogging(loggingType)
			expectedValue = value
			toTestValue = logging.Logzio

		case "elastic_logging-aws":
			value, _, _ := generateTestElasticLogging()
			expectedValue = value.AWS
			toTestValue = logging.Elastic.AWS

		case "elastic_logging-elastic_cloud":
			value, _, _ := generateTestElasticLogging()
			expectedValue = value.ElasticCloud
			toTestValue = logging.Elastic.ElasticCloud

		default:
			return nil
		}

		if diff := deep.Equal(expectedValue, toTestValue); diff != nil {
			return fmt.Errorf("%s attributes do not match. Diff: %s", loggingType, diff)
		}

		return nil
	}
}

func testAccCheckControlPlaneOrgCheckDestroy(s *terraform.State) error {

	// TestLogger.Printf("Inside testAccCheckControlPlaneOrgCheckDestroy. Resources Length: %d", len(s.RootModule().Resources))

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		TestLogger.Printf("Inside testAccCheckControlPlaneOrgCheckDestroy: rs.Type: %s", rs.Type)

		if rs.Type != "cpln_org_logging" {
			continue
		}

		orgName := rs.Primary.ID

		TestLogger.Printf("Inside testAccCheckControlPlaneOrgCheckDestroy: Org name: %s", orgName)

		org, _, _ := c.GetOrg()

		if org.Spec.Logging != nil {
			return fmt.Errorf("Org Spec Logging still exists. Org Name: %s", *org.Name)
		}
	}

	return nil
}

/*** Unit Tests ***/
// Build Functions //
func TestControlPlane_BuildS3Logging(t *testing.T) {
	s3Logging, expectedS3Logging, _ := generateTestS3Logging()
	if diff := deep.Equal(s3Logging, &expectedS3Logging); diff != nil {
		t.Errorf("Coralogix Logging was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildCoralogixLogging(t *testing.T) {
	coralogixLogging, expectedCoralogixLogging, _ := generateTestCoralogixLogging()
	if diff := deep.Equal(coralogixLogging, &expectedCoralogixLogging); diff != nil {
		t.Errorf("Coralogix Logging was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildDatadogLogging(t *testing.T) {
	datadogLogging, expectedDatadogLogging, _ := generateTestDatadogLogging()
	if diff := deep.Equal(datadogLogging, &expectedDatadogLogging); diff != nil {
		t.Errorf("Datadog Logging was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildLogzioLogging(t *testing.T) {
	logzioLogging, expectedLogzioLogging, _ := generateTestLogzioLogging("logzio_logging")
	if diff := deep.Equal(logzioLogging, &expectedLogzioLogging); diff != nil {
		t.Errorf("Logzio Logging was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildLogzioLogging_DifferentListenerHost(t *testing.T) {
	logzio, expectedLogzioLogging, _ := generateTestLogzioLogging("logzio_logging-different_listener_host")
	if diff := deep.Equal(logzio, &expectedLogzioLogging); diff != nil {
		t.Errorf("Logzio Logging with different listener host was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildElasticLogging(t *testing.T) {
	elasticLogging, expectedElasticLogging, _ := generateTestElasticLogging()
	if diff := deep.Equal(elasticLogging, &expectedElasticLogging); diff != nil {
		t.Errorf("AWS Logging was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildAWSLogging(t *testing.T) {
	awsLogging, expectedAWSLogging, _ := generateTestAWSLogging()
	if diff := deep.Equal(awsLogging, &expectedAWSLogging); diff != nil {
		t.Errorf("AWS Logging was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildElasticCloudLogging(t *testing.T) {
	elasticCloudLogging, expectedElasticCloudLogging, _ := generateTestElasticCloudLogging()
	if diff := deep.Equal(elasticCloudLogging, &expectedElasticCloudLogging); diff != nil {
		t.Errorf("Elastic Cloud Logging was not built correctly. Diff: %s", diff)
	}
}

/*** Generate Functions ***/
func generateTestS3Logging() (*client.S3Logging, client.S3Logging, []interface{}) {
	bucket := "test-bucket"
	region := "us-east1"
	prefix := "/"
	credentials := "/org/terraform-test-org/secret/aws-random-tbd"

	flattened := generateFlatTestS3Logging(bucket, region, prefix, credentials)
	s3Logging := buildS3Logging(flattened).S3
	expectedS3Logging := client.S3Logging{
		Bucket:      &bucket,
		Prefix:      &prefix,
		Region:      &region,
		Credentials: &credentials,
	}

	return s3Logging, expectedS3Logging, flattened
}

func generateTestCoralogixLogging() (*client.CoralogixLogging, client.CoralogixLogging, []interface{}) {
	cluster := "coralogix.com"
	credentials := "/org/terraform-test-org/secret/opaque-random-coralogix-tbd"
	app := "{workload}"
	subsystem := "{org}"

	flattened := generateFlatTestCoralogixLogging(cluster, credentials, app, subsystem)
	coralogixLogging := buildCoralogixLogging(flattened).Coralogix
	expectedCoralogixLogging := client.CoralogixLogging{
		Cluster:     &cluster,
		Credentials: &credentials,
		App:         &app,
		Subsystem:   &subsystem,
	}

	return coralogixLogging, expectedCoralogixLogging, flattened
}

func generateTestDatadogLogging() (*client.DatadogLogging, client.DatadogLogging, []interface{}) {
	host := "http-intake.logs.datadoghq.com"
	credentials := "/org/terraform-test-org/secret/opaque-random-datadog-tbd"

	flattened := generateFlatTestDatadogLogging(host, credentials)
	datadogLogging := buildDatadogLogging(flattened).Datadog
	expectedDatadogLogging := client.DatadogLogging{
		Host:        &host,
		Credentials: &credentials,
	}

	return datadogLogging, expectedDatadogLogging, flattened
}

func generateTestLogzioLogging(loggingType string) (*client.LogzioLogging, client.LogzioLogging, []interface{}) {
	listenerHost := "listener.logz.io"
	credentials := "/org/terraform-test-org/secret/opaque-random-datadog-tbd"

	if loggingType == "logzio_logging-different_listener_host" {
		listenerHost = "listener-nl.logz.io"
	}

	flattened := generateFlatTestLogzioLogging(listenerHost, credentials)
	logzio := buildLogzioLogging(flattened).Logzio
	expectedLogzioLogging := client.LogzioLogging{
		ListenerHost: &listenerHost,
		Credentials:  &credentials,
	}

	return logzio, expectedLogzioLogging, flattened
}

func generateTestElasticLogging() (*client.ElasticLogging, client.ElasticLogging, []interface{}) {
	_, expectedAWSLogging, flattenedAWSLogging := generateTestAWSLogging()
	_, expectedElasticCloudLogging, flattenedElasticCloudLogging := generateTestElasticCloudLogging()

	flattened := generateFlatTestElasticLogging(flattenedAWSLogging, flattenedElasticCloudLogging)
	elasticLogging := buildElasticLogging(flattened)
	expectedElasticLogging := client.ElasticLogging{
		AWS:          &expectedAWSLogging,
		ElasticCloud: &expectedElasticCloudLogging,
	}

	return elasticLogging.Elastic, expectedElasticLogging, flattened
}

func generateTestAWSLogging() (*client.AWSLogging, client.AWSLogging, []interface{}) {
	host := "es.amazonaws.com"
	port := 8080
	index := "my-index"
	loggingType := "my-type"
	credentials := "/org/terraform-test-org/secret/opaque-random-elastic-logging-aws-tbd"
	region := "us-east-1"

	flattened := generateFlatTestAWSLogging(host, port, index, loggingType, credentials, region)
	awsLogging := buildAWSLogging(flattened)
	expectedAWSLogging := client.AWSLogging{
		Host:        &host,
		Port:        &port,
		Index:       &index,
		Type:        &loggingType,
		Credentials: &credentials,
		Region:      &region,
	}

	return awsLogging, expectedAWSLogging, flattened
}

func generateTestElasticCloudLogging() (*client.ElasticCloudLogging, client.ElasticCloudLogging, []interface{}) {
	index := "my-index"
	loggingType := "my-type"
	credentials := "/org/terraform-test-org/secret/opaque-random-elastic-logging-elastic-cloud-tbd"
	cloudId := "my-cloud-id"

	flattened := generateFlatTestElasticCloudLogging(index, loggingType, credentials, cloudId)
	elasticCloudLogging := buildElasticCloudLogging(flattened)
	expectedElasticCloudLogging := client.ElasticCloudLogging{
		Index:       &index,
		Type:        &loggingType,
		Credentials: &credentials,
		CloudID:     &cloudId,
	}

	return elasticCloudLogging, expectedElasticCloudLogging, flattened
}

/*** Flatten Functions ***/
func generateFlatTestS3Logging(bucket string, region string, prefix string, credentials string) []interface{} {
	spec := map[string]interface{}{
		"bucket":      bucket,
		"region":      region,
		"prefix":      prefix,
		"credentials": credentials,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestCoralogixLogging(cluster string, credentials string, app string, subsystem string) []interface{} {
	spec := map[string]interface{}{
		"cluster":     cluster,
		"credentials": credentials,
		"app":         app,
		"subsystem":   subsystem,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestDatadogLogging(host string, credentials string) []interface{} {
	spec := map[string]interface{}{
		"host":        host,
		"credentials": credentials,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestLogzioLogging(listenerHost string, credentials string) []interface{} {
	spec := map[string]interface{}{
		"listener_host": listenerHost,
		"credentials":   credentials,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestElasticLogging(awsLogging []interface{}, elasticCloudLogging []interface{}) []interface{} {
	spec := map[string]interface{}{
		"aws":           awsLogging,
		"elastic_cloud": elasticCloudLogging,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestAWSLogging(host string, port int, index string, loggingType string, credentials string, region string) []interface{} {
	spec := map[string]interface{}{
		"host":        host,
		"port":        port,
		"index":       index,
		"type":        loggingType,
		"credentials": credentials,
		"region":      region,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestElasticCloudLogging(index string, loggingType string, credentials string, cloudId string) []interface{} {
	spec := map[string]interface{}{
		"index":       index,
		"type":        loggingType,
		"credentials": credentials,
		"cloud_id":    cloudId,
	}

	return []interface{}{
		spec,
	}
}
