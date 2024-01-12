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

	var testLogging []client.Logging

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "ORG_LOGGING") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneOrgCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneOrgS3(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
				),
			},
			{
				Config: testAccControlPlaneOrgCoralogix(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
				),
			},
			{
				Config: testAccControlPlaneOrgDatadog(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
				),
			},
			{
				Config: testAccControlPlaneOrgLogzio(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
				),
			},
			{
				Config: testAccControlPlaneOrgLogzioWithDifferentListenerHost(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
				),
			},
			{
				Config: testAccControlPlaneOrgElasticAWS(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
				),
			},
			{
				Config: testAccControlPlaneOrgElasticCloud(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
				),
			},
			{
				Config: testAccControlPlaneOrgElasticGeneric(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
				),
			},
			{
				Config: testAccControlPlaneOrgThreeUniqueLoggings(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
				),
			},
			{
				Config: testAccControlPlaneOrgTwoUniqueLoggings(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneLoggingExists("cpln_org_logging.tf-logging", &testLogging),
					testAccCheckControlPlaneLoggingAttributes(&testLogging),
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

	resource "cpln_secret" "opaque-1" {

        name = "opaque-random-datadog-tbd-1"
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

		datadog_logging {
			host = "http-intake.logs.datadoghq.com"

			// Opaque Secret Only
			credentials = cpln_secret.opaque-1.self_link  
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

func testAccControlPlaneOrgElasticGeneric() string {

	TestLogger.Printf("Inside testAccControlPlaneOrgElasticGeneric")

	return `

	resource "cpln_secret" "opaque" {

        name = "opaque-random-elastic-logging-generic-tbd"
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
			generic {
				host  = "example.com"
				port  = 9200
				path  = "/var/log/elasticsearch/"
				index = "my-index"
				type  = "my-type"
				credentials = cpln_secret.opaque.self_link
			}
        }
    }
    `
}

func testAccControlPlaneOrgThreeUniqueLoggings() string {
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

	resource "cpln_secret" "opaque-coralogix" {

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

	resource "cpln_secret" "opaque-datadog" {

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

	resource "cpln_secret" "opaque-datadog-1" {

		name = "opaque-random-datadog-tbd-1"
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

		s3_logging {
			bucket = "test-bucket"
			region = "us-east1"
			prefix = "/"

			// AWS Secret Only
			credentials = cpln_secret.aws.self_link
		}


		coralogix_logging {
			cluster = "coralogix.com"

			// Opaque Secret Only
			credentials = cpln_secret.opaque-coralogix.self_link

			app = "{workload}"
			subsystem = "{org}"
		}


		datadog_logging {
			host = "http-intake.logs.datadoghq.com"

			// Opaque Secret Only
			credentials = cpln_secret.opaque-datadog.self_link  
		}

		datadog_logging {
			host = "http-intake.logs.datadoghq.com"

			// Opaque Secret Only
			credentials = cpln_secret.opaque-datadog-1.self_link  
		}
	}

	`
}

func testAccControlPlaneOrgTwoUniqueLoggings() string {
	return `

	resource "cpln_secret" "opaque-coralogix" {

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

	resource "cpln_secret" "opaque-datadog" {

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

	resource "cpln_secret" "opaque-datadog-1" {

		name = "opaque-random-datadog-tbd-1"
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

		coralogix_logging {
			cluster = "coralogix.com"

			// Opaque Secret Only
			credentials = cpln_secret.opaque-coralogix.self_link

			app = "{workload}"
			subsystem = "{org}"
		}


		datadog_logging {
			host = "http-intake.logs.datadoghq.com"

			// Opaque Secret Only
			credentials = cpln_secret.opaque-datadog.self_link  
		}

		datadog_logging {
			host = "http-intake.logs.datadoghq.com"

			// Opaque Secret Only
			credentials = cpln_secret.opaque-datadog-1.self_link  
		}
	}

	`
}

func testAccCheckControlPlaneLoggingExists(resourceName string, loggings *[]client.Logging) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		TestLogger.Printf("Inside testAccCheckControlPlaneLoggingExists. Resources Length: %d", len(s.RootModule().Resources))

		_, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Resource not found: %s", s)
		}

		c := testAccProvider.Meta().(*client.Client)
		org, _, err := c.GetOrg()

		if err != nil {
			return err
		}

		if org.Spec == nil {
			return fmt.Errorf("Org spec is nil")
		}

		_loggings := []client.Logging{}

		if org.Spec.Logging != nil {
			_loggings = append(_loggings, *org.Spec.Logging)
		}

		if org.Spec.ExtraLogging != nil && len(*org.Spec.ExtraLogging) != 0 {
			_loggings = append(_loggings, *org.Spec.ExtraLogging...)
		}

		*loggings = _loggings

		return nil
	}
}

func testAccCheckControlPlaneLoggingAttributes(loggings *[]client.Logging) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if loggings == nil || len(*loggings) == 0 {
			return fmt.Errorf("Logging Attributes: Loggings were nil or empty")
		}

		var s3TestValue, coralogixTestValue, dataDogTestValue,
			logzioTestvalue, elasticTestValue []client.Logging

		for _, logging := range *loggings {

			var loggingType string

			// Determine logging type
			if logging.S3 != nil {
				loggingType = "s3_logging"
			} else if logging.Coralogix != nil {
				loggingType = "coralogix_logging"
			} else if logging.Datadog != nil {
				loggingType = "datadog_logging"
			} else if logging.Logzio != nil {
				loggingType = "logzio_logging"

				if *logging.Logzio.ListenerHost == "listener-nl.logz.io" {
					loggingType = "logzio_logging-different_listener_host"
				}
			} else if logging.Elastic != nil {
				if logging.Elastic.AWS != nil {
					loggingType = "elastic_logging-aws"
				} else if logging.Elastic.ElasticCloud != nil {
					loggingType = "elastic_logging-elastic_cloud"
				} else if logging.Elastic.Generic != nil {
					loggingType = "elastic_logging-generic"
				}
			} else {
				return fmt.Errorf("Logging Attributes: We were not able to determine logging type")
			}

			switch loggingType {
			case "s3_logging":

				temp := client.Logging{
					S3: logging.S3,
				}
				s3TestValue = append(s3TestValue, temp)

			case "coralogix_logging":

				temp := client.Logging{
					Coralogix: logging.Coralogix,
				}

				coralogixTestValue = append(coralogixTestValue, temp)

			case "datadog_logging":

				temp := client.Logging{
					Datadog: logging.Datadog,
				}

				dataDogTestValue = append(dataDogTestValue, temp)

			case "logzio_logging":

				temp := client.Logging{
					Logzio: logging.Logzio,
				}

				logzioTestvalue = append(logzioTestvalue, temp)

			case "elastic_logging-aws":
			case "elastic_logging-elastic_cloud":
			case "elastic_logging-generic":

				temp := client.Logging{
					Elastic: logging.Elastic,
				}

				elasticTestValue = append(elasticTestValue, temp)

			default:
				return nil
			}
		}

		for _, logging := range loggingNames {

			var expectedValue interface{}
			var toTestValue interface{}

			switch logging {
			case "s3_logging":
				expectedValue, _, _ = generateTestS3Logging()
				toTestValue = s3TestValue

			case "coralogix_logging":
				expectedValue, _, _ = generateTestCoralogixLogging()
				toTestValue = coralogixTestValue

			case "datadog_logging":
				expectedValue, _, _ = generateTestDatadogLogging2()
				toTestValue = dataDogTestValue

			case "logzio_logging":
				expectedValue, _, _ = generateTestLogzioLogging("")
				toTestValue = logzioTestvalue

			case "elastic_logging":
				var loggingType string

				if elasticTestValue != nil {
					if elasticTestValue[0].Elastic.AWS != nil {
						loggingType = "elastic_logging-aws"
					} else if elasticTestValue[0].Elastic.ElasticCloud != nil {
						loggingType = "elastic_logging-elastic_cloud"
					} else if elasticTestValue[0].Elastic.Generic != nil {
						loggingType = "elastic_logging-generic"
					} else {
						return fmt.Errorf("Logging Attributes: Unable to get logging type from elastic logging")
					}
				}

				expectedValue, _, _ = generateTestElasticLogging(loggingType)
				toTestValue = elasticTestValue

			default:
				return nil
			}

			if toTestValue != nil && len(toTestValue.([]client.Logging)) > 0 {

				if diff := deep.Equal(expectedValue, toTestValue); diff != nil {
					return fmt.Errorf("%s attributes do not match. Diff: %s", logging, diff)
				}
			}
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

		if org.Spec.Logging != nil || (org.Spec.ExtraLogging != nil && len(*org.Spec.ExtraLogging) != 0) {
			return fmt.Errorf("Org Spec Logging still exists. Org Name: %s", *org.Name)
		}
	}

	return nil
}

/*** Unit Tests ***/
// Build Functions //

func TestControlPlane_BuildS3Logging(t *testing.T) {
	s3Logging, expectedS3Logging, _ := generateTestS3Logging()
	if diff := deep.Equal(s3Logging, expectedS3Logging); diff != nil {
		t.Errorf("Coralogix Logging was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildCoralogixLogging(t *testing.T) {
	coralogixLogging, expectedCoralogixLogging, _ := generateTestCoralogixLogging()
	if diff := deep.Equal(coralogixLogging, expectedCoralogixLogging); diff != nil {
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
	if diff := deep.Equal(logzioLogging, expectedLogzioLogging); diff != nil {
		t.Errorf("Logzio Logging was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildLogzioLogging_DifferentListenerHost(t *testing.T) {
	logzio, expectedLogzioLogging, _ := generateTestLogzioLogging("logzio_logging-different_listener_host")
	if diff := deep.Equal(logzio, expectedLogzioLogging); diff != nil {
		t.Errorf("Logzio Logging with different listener host was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildElasticLogging(t *testing.T) {
	elasticLogging, expectedElasticLogging, _ := generateTestElasticLogging("")
	if diff := deep.Equal(elasticLogging, expectedElasticLogging); diff != nil {
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

func TestControlPlane_BuildGenericLogging(t *testing.T) {
	genericLogging, expectedGenericLogging, _ := generateTestGenericLogging()
	if diff := deep.Equal(genericLogging, &expectedGenericLogging); diff != nil {
		t.Errorf("Elastic Generic Logging was not built correctly. Diff: %s", diff)
	}
}

/*** Generate Functions ***/

func generateTestS3Logging() ([]client.Logging, []client.Logging, []interface{}) {

	bucket := "test-bucket"
	region := "us-east1"
	prefix := "/"
	credentials := "/org/terraform-test-org/secret/aws-random-tbd"

	flattened := generateFlatTestS3Logging(bucket, region, prefix, credentials)
	s3Logging := buildS3Logging(flattened)

	expectedS3 := client.S3Logging{
		Bucket:      &bucket,
		Prefix:      &prefix,
		Region:      &region,
		Credentials: &credentials,
	}

	expectedS3Logging := client.Logging{
		S3: &expectedS3,
	}

	output := []client.Logging{
		expectedS3Logging,
	}

	return s3Logging, output, flattened
}

func generateTestCoralogixLogging() ([]client.Logging, []client.Logging, []interface{}) {

	cluster := "coralogix.com"
	credentials := "/org/terraform-test-org/secret/opaque-random-coralogix-tbd"
	app := "{workload}"
	subsystem := "{org}"

	flattened := generateFlatTestCoralogixLogging(cluster, credentials, app, subsystem)
	coralogixLogging := buildCoralogixLogging(flattened)

	expectedCoralogix := client.CoralogixLogging{
		Cluster:     &cluster,
		Credentials: &credentials,
		App:         &app,
		Subsystem:   &subsystem,
	}

	expectedS3Logging := client.Logging{
		Coralogix: &expectedCoralogix,
	}

	output := []client.Logging{
		expectedS3Logging,
	}

	return coralogixLogging, output, flattened
}

func generateTestDatadogLogging() (*client.DatadogLogging, client.DatadogLogging, []interface{}) {
	host := "http-intake.logs.datadoghq.com"
	credentials := "/org/terraform-test-org/secret/opaque-random-datadog-tbd"

	flattened := generateFlatTestDatadogLogging(host, credentials)
	datadogLogging := buildDatadogLogging(flattened)[0].Datadog
	expectedDatadogLogging := client.DatadogLogging{
		Host:        &host,
		Credentials: &credentials,
	}

	return datadogLogging, expectedDatadogLogging, flattened
}

func generateTestDatadogLogging2() ([]client.Logging, []client.Logging, []interface{}) {

	host := "http-intake.logs.datadoghq.com"
	credentials := "/org/terraform-test-org/secret/opaque-random-datadog-tbd"

	host2 := "http-intake.logs.datadoghq.com"
	credentials2 := "/org/terraform-test-org/secret/opaque-random-datadog-tbd-1"

	flattened := generateFlatTestDatadogLogging2(host, credentials, host2, credentials2)
	datadogLogging := buildDatadogLogging(flattened)

	expectedDatadogLogging1 := client.DatadogLogging{
		Host:        &host,
		Credentials: &credentials,
	}

	expectedDatadogLogging2 := client.DatadogLogging{
		Host:        &host2,
		Credentials: &credentials2,
	}

	log1 := client.Logging{
		Datadog: &expectedDatadogLogging1,
	}

	log2 := client.Logging{
		Datadog: &expectedDatadogLogging2,
	}

	output := []client.Logging{
		log1,
		log2,
	}

	return datadogLogging, output, flattened
}

func generateTestLogzioLogging(loggingType string) ([]client.Logging, []client.Logging, []interface{}) {

	listenerHost := "listener.logz.io"
	credentials := "/org/terraform-test-org/secret/opaque-random-datadog-tbd"

	if loggingType == "logzio_logging-different_listener_host" {
		listenerHost = "listener-nl.logz.io"
	}

	flattened := generateFlatTestLogzioLogging(listenerHost, credentials)
	logzio := buildLogzioLogging(flattened)

	expectedLogzioLogging := client.LogzioLogging{
		ListenerHost: &listenerHost,
		Credentials:  &credentials,
	}

	expectedS3Logging := client.Logging{
		Logzio: &expectedLogzioLogging,
	}

	output := []client.Logging{
		expectedS3Logging,
	}

	return logzio, output, flattened
}

func generateTestElasticLogging(loggingType string) ([]client.Logging, []client.Logging, []interface{}) {
	var expectedAWSLogging client.AWSLogging
	var expectedElasticCloudLogging client.ElasticCloudLogging
	var expectedGenericLogging client.GenericLogging

	var flattenedAWSLogging []interface{}
	var flattenedElasticCloudLogging []interface{}
	var flattenedGenericLogging []interface{}

	switch loggingType {
	case "elastic_logging-aws":
		_, expectedAWSLogging, flattenedAWSLogging = generateTestAWSLogging()
	case "elastic_logging-elastic_cloud":
		_, expectedElasticCloudLogging, flattenedElasticCloudLogging = generateTestElasticCloudLogging()
	case "elastic_logging-generic":
		_, expectedGenericLogging, flattenedGenericLogging = generateTestGenericLogging()
	default:
		_, expectedAWSLogging, flattenedAWSLogging = generateTestAWSLogging()
		_, expectedElasticCloudLogging, flattenedElasticCloudLogging = generateTestElasticCloudLogging()
		_, expectedGenericLogging, flattenedGenericLogging = generateTestGenericLogging()
	}

	flattened := generateFlatTestElasticLogging(flattenedAWSLogging, flattenedElasticCloudLogging, flattenedGenericLogging)
	elasticLogging := buildElasticLogging(flattened)
	expectedElasticLogging := client.ElasticLogging{
		AWS:          &expectedAWSLogging,
		ElasticCloud: &expectedElasticCloudLogging,
		Generic:      &expectedGenericLogging,
	}

	expectedLogging := client.Logging{
		Elastic: &expectedElasticLogging,
	}

	output := []client.Logging{
		expectedLogging,
	}

	return elasticLogging, output, flattened
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

func generateTestGenericLogging() (*client.GenericLogging, client.GenericLogging, []interface{}) {
	host := "example.com"
	port := 9200
	path := "/var/log/elasticsearch/"
	index := "my-index"
	loggingType := "my-type"
	credentials := "/org/terraform-test-org/secret/opaque-random-elastic-logging-generic-tbd"

	flattened := generateFlatTestGenericLogging(host, port, path, index, loggingType, credentials)
	genericLogging := buildGenericLogging(flattened)
	expectedGenericLogging := client.GenericLogging{
		Host:        &host,
		Port:        &port,
		Path:        &path,
		Index:       &index,
		Type:        &loggingType,
		Credentials: &credentials,
	}

	return genericLogging, expectedGenericLogging, flattened
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

func generateFlatTestDatadogLogging2(host, credential, host2, credential2 string) []interface{} {
	spec := map[string]interface{}{
		"host":        host,
		"credentials": credential,
	}

	spec2 := map[string]interface{}{
		"host":        host2,
		"credentials": credential2,
	}

	return []interface{}{
		spec,
		spec2,
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

func generateFlatTestElasticLogging(awsLogging []interface{}, elasticCloudLogging []interface{}, genericLogging []interface{}) []interface{} {
	spec := make(map[string]interface{})

	if awsLogging != nil {
		spec["aws"] = awsLogging
	}

	if elasticCloudLogging != nil {
		spec["elastic_cloud"] = elasticCloudLogging
	}

	if genericLogging != nil {
		spec["generic"] = genericLogging
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

func generateFlatTestGenericLogging(host string, port int, path string, index string, loggingType string, credentials string) []interface{} {
	spec := map[string]interface{}{
		"host":        host,
		"port":        port,
		"path":        path,
		"index":       index,
		"type":        loggingType,
		"credentials": credentials,
	}

	return []interface{}{
		spec,
	}
}
