package cpln

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDataSourceImage_basic performs an acceptance test for the data source.
func TestAccControlPlaneDataSourceImage_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewImageDataSourceTest()

	// Run the acceptance test case for the data source, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DATA_SOURCE_IMAGE") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// ImageDataSourceTest defines the necessary functionality to test the data source.
type ImageDataSourceTest struct {
	Steps []resource.TestStep
}

// NewImageDataSourceTest creates a ImageDataSourceTest with initialized test cases.
func NewImageDataSourceTest() ImageDataSourceTest {
	// Create a data source test instance
	dataSourceTest := ImageDataSourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, dataSourceTest.NewImageWithTagScenario()...)
	steps = append(steps, dataSourceTest.NewImageNameOnlyScenario()...)

	// Set the cases for the data source test
	dataSourceTest.Steps = steps

	// Return the data source test
	return dataSourceTest
}

// Test Scenarios //

// NewImageWithTagScenario creates a test case with initial and updated configurations.
func (idst *ImageDataSourceTest) NewImageWithTagScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "specific-image"
	name := "cpln_doc_demo:7"

	// Build test steps
	_, initialStep := idst.BuildImageWithTagTestStep(dataSourceName, name)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// NewImageNameOnlyScenario creates a test case with initial and updated configurations.
func (idst *ImageDataSourceTest) NewImageNameOnlyScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "latest-image"
	nameOnly := "call-internal-service-3000"
	latestImageName := "call-internal-service-3000:6"

	// Build test steps
	_, initialStep := idst.BuildImageNameOnlyTestStep(dataSourceName, nameOnly, latestImageName)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// Test Cases //

// BuildImageWithTagTestStep returns a test case for the data source.
func (idst *ImageDataSourceTest) BuildImageWithTagTestStep(dataSourceName string, name string) (ImageDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := ImageDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "image",
			ResourceName:    dataSourceName,
			Name:            name,
			Description:     name,
			ResourceAddress: fmt.Sprintf("data.cpln_image.%s", dataSourceName),
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: idst.ImageWithTagHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", c.Name),
			c.TestCheckResourceAttr("name", c.Name),
			c.TestCheckResourceAttr("tags.%", "0"),
			c.TestCheckResourceAttr("self_link", c.GetSelfLink()),
			c.TestCheckResourceAttr("tag", "7"),
			c.TestCheckResourceAttr("repository", "cpln_doc_demo"),
			c.TestCheckResourceAttr("digest", "sha256:f989803c03061af7498903457d929ff6ab3bfaa440c2b0fce88e7a1b708942cc"),
			c.TestCheckNestedBlocks("manifest", []map[string]interface{}{
				{
					"config": []map[string]interface{}{
						{
							"size":       "9115",
							"digest":     "sha256:f989803c03061af7498903457d929ff6ab3bfaa440c2b0fce88e7a1b708942cc",
							"media_type": "application/vnd.docker.container.image.v1+json",
						},
					},
					"layers": []map[string]interface{}{
						{
							"digest":     "sha256:7919f5b7d60254cafc73c0d097b8ccffb72e0b6472957ece4dd5b378c5ca7cc1",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "45377037",
						},
						{
							"digest":     "sha256:0e107167dcc5392ce7b34cb6af6bcfe1cf99f76f9e632d6c321526666558ab50",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "10752152",
						},
						{
							"digest":     "sha256:66a456bba435b99e4c17dd5da957e63bef2c43ae5291b055ce82bd26d9153259",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "4340590",
						},
						{
							"digest":     "sha256:5435318a0426be8944e8142c90c8501d9f27e23cbd27e1094311e292dcfbc926",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "50110915",
						},
						{
							"digest":     "sha256:8494dd3284650baec0f038af6bec4db4ddbcf520c80f88d690d5e0327e3f4d80",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "214315680",
						},
						{
							"digest":     "sha256:3b01939c65060e98dcf4b23f8fd6b94ec5638cc32c35301aeed7f422c7d543da",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "4160",
						},
						{
							"digest":     "sha256:5caada2abfdf81185a9204fcdabc7719b1b647ff8fd4399ee481ef91c17fa78e",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "33436702",
						},
						{
							"digest":     "sha256:97297040d40f1047a1cb233b9d31772a3b1bbb76310df4ea55a6b47f194ef621",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "2366773",
						},
						{
							"digest":     "sha256:538838247bdb13ea900d492a5cd9a0b3900464a9296111328cf85fa4d7a30414",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "294",
						},
						{
							"digest":     "sha256:38c9d16fed1bfcc281203742259e52c2216f0fbe60a8cde6edd7b0ade59b0c00",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "127",
						},
						{
							"digest":     "sha256:ffde8b83c412be620342086591809a58be400410901591302c7a7364a217e389",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "605093",
						},
						{
							"digest":     "sha256:70cab3cefbc6f79eb11a08dadb04f7f262766b06eacbc8a6b5dc98a4597c66e4",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "36724745",
						},
						{
							"digest":     "sha256:598f2e47dadd51c61e4da70c4562e6825fdf3f11f98c8a09159b923bb997e5c6",
							"media_type": "application/vnd.docker.image.rootfs.diff.tar.gzip",
							"size":       "31330372",
						},
					},
					"media_type":     "application/vnd.docker.distribution.manifest.v2+json",
					"schema_version": "2",
				},
			}),
		),
	}
}

// BuildImageWithTagTestStep returns a test case for the data source.
func (idst *ImageDataSourceTest) BuildImageNameOnlyTestStep(dataSourceName string, name string, latestImageName string) (ImageDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := ImageDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "image",
			ResourceName:    dataSourceName,
			Name:            latestImageName,
			Description:     latestImageName,
			ResourceAddress: fmt.Sprintf("data.cpln_image.%s", dataSourceName),
		},
		ImageName: name,
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: idst.ImageNameOnlyHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.TestCheckResourceAttr("id", latestImageName),
			c.TestCheckResourceAttr("name", latestImageName),
			c.TestCheckResourceAttr("tags.%", "0"),
			c.TestCheckResourceAttr("self_link", c.GetSelfLink()),
			c.TestCheckResourceAttr("tag", "6"),
			c.TestCheckResourceAttr("repository", "call-internal-service-3000"),
			c.TestCheckResourceAttr("digest", "sha256:bcfe045e21c71864c556ae5a2a1241321828c89de2c54ce4f752468a3f451cc9"),
			c.TestCheckNestedBlocks("manifest", []map[string]interface{}{
				{
					"config": []map[string]interface{}{
						{
							"size":       "8053",
							"digest":     "sha256:bcfe045e21c71864c556ae5a2a1241321828c89de2c54ce4f752468a3f451cc9",
							"media_type": "application/vnd.docker.container.image.v1+json",
						},
					},
					// Skip layers check
					"media_type":     "application/vnd.docker.distribution.manifest.v2+json",
					"schema_version": "2",
				},
			}),
		),
	}
}

// Configs //

// ImageWithTagHcl returns a data source HCL.
func (idst *ImageDataSourceTest) ImageWithTagHcl(c ImageDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_image" "%s" {
  name = "%s"
}
`, c.ResourceName, c.Name)
}

// ImageWithTagHcl returns a data source HCL.
func (idst *ImageDataSourceTest) ImageNameOnlyHcl(c ImageDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_image" "%s" {
  name = "%s"
}
`, c.ResourceName, c.ImageName)
}

/*** Data Source Test Case ***/

// ImageDataSourceTestCase defines a specific data source test case.
type ImageDataSourceTestCase struct {
	ProviderTestCase
	ImageName string
}
