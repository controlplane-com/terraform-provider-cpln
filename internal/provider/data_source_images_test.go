package cpln

import (
	"fmt"
	"testing"

	"github.com/go-test/deep"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

/*** Acceptance Test ***/

// TestAccControlPlaneDataSourceImages_basic performs an acceptance test for the data source.
func TestAccControlPlaneDataSourceImages_basic(t *testing.T) {
	// Initialize the test
	dataSourceTest := NewImagesDataSourceTest()

	// Run the acceptance test case for the data source, covering create, read, update, and import functionalities
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "DATA_SOURCE_IMAGES") },
		ProtoV6ProviderFactories: GetProviderServer(),
		Steps:                    dataSourceTest.Steps,
	})
}

/*** Data Source Test ***/

// ImagesDataSourceTest defines the necessary functionality to test the data source.
type ImagesDataSourceTest struct {
	Steps []resource.TestStep
}

// NewImagesDataSourceTest creates a ImagesDataSourceTest with initialized test cases.
func NewImagesDataSourceTest() ImagesDataSourceTest {
	// Create a data source test instance
	dataSourceTest := ImagesDataSourceTest{}

	// Initialize the test steps slice
	steps := []resource.TestStep{}

	// Fill the steps slice
	steps = append(steps, dataSourceTest.NewAllImagesScenario()...)
	steps = append(steps, dataSourceTest.NewSpecificImagesScenario()...)

	// Set the cases for the data source test
	dataSourceTest.Steps = steps

	// Return the data source test
	return dataSourceTest
}

// Test Scenarios //

// NewAllImagesScenario creates a test case with initial and updated configurations.
func (idst *ImagesDataSourceTest) NewAllImagesScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "all-images"

	// Build test steps
	_, initialStep := idst.BuildAllImagesTestStep(dataSourceName)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// NewSpecificImagesScenario creates a test case with initial and updated configurations.
func (idst *ImagesDataSourceTest) NewSpecificImagesScenario() []resource.TestStep {
	// Define necessary variables
	dataSourceName := "specific-images"

	// Build test steps
	_, initialStep := idst.BuildSpecificImagesTestStep(dataSourceName)

	// Return the complete test steps
	return []resource.TestStep{
		// Read
		initialStep,
	}
}

// Test Cases //

// BuildAllImagesTestStep returns a test case for the data source.
func (idst *ImagesDataSourceTest) BuildAllImagesTestStep(dataSourceName string) (ImagesDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := ImagesDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "image",
			ResourceName:    dataSourceName,
			ResourceAddress: fmt.Sprintf("data.cpln_images.%s", dataSourceName),
		},
		Query: client.Query{
			Kind: StringPointer("image"),
			Spec: &client.QuerySpec{
				Match: StringPointer("all"),
			},
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: idst.AllImagesHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.Attributes(19),
			resource.TestCheckResourceAttrSet(c.ResourceAddress, "images.#"),
		),
	}
}

// BuildSpecificImagesTestStep returns a test case for the data source.
func (idst *ImagesDataSourceTest) BuildSpecificImagesTestStep(dataSourceName string) (ImagesDataSourceTestCase, resource.TestStep) {
	// Create the test case with metadata and descriptions
	c := ImagesDataSourceTestCase{
		ProviderTestCase: ProviderTestCase{
			Kind:            "image",
			ResourceName:    dataSourceName,
			ResourceAddress: fmt.Sprintf("data.cpln_images.%s", dataSourceName),
		},
		Query: client.Query{
			Kind: StringPointer("image"),
			Spec: &client.QuerySpec{
				Match: StringPointer("all"),
				Terms: &[]client.QueryTerm{
					{
						Op:       StringPointer("="),
						Property: StringPointer("repository"),
						Value:    StringPointer("call-internal-service-3000"),
					},
				},
			},
		},
	}

	// Initialize and return the inital test step
	return c, resource.TestStep{
		Config: idst.SpecificImagesHcl(c),
		Check: resource.ComposeAggregateTestCheckFunc(
			c.Exists(),
			c.Attributes(5),
			resource.TestCheckResourceAttrSet(c.ResourceAddress, "images.#"),
		),
	}
}

// Configs //

// AllImagesHcl returns a data source HCL.
func (idst *ImagesDataSourceTest) AllImagesHcl(c ImagesDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_images" "%s" {}
`, c.ResourceName)
}

// SpecificImagesHcl returns a data source HCL.
func (idst *ImagesDataSourceTest) SpecificImagesHcl(c ImagesDataSourceTestCase) string {
	return fmt.Sprintf(`
data "cpln_images" "%s" {
  query {
    fetch = "items"
    spec {
      match = "all"
      terms {
        op       = "="
        property = "repository"
        value    = "call-internal-service-3000"
      }
    }
  }
}
`, c.ResourceName)
}

/*** Data Source Test Case ***/

// ImagesDataSourceTestCase defines a specific data source test case.
type ImagesDataSourceTestCase struct {
	ProviderTestCase
	Query  client.Query
	Images *client.ImagesQueryResult
}

// Exists verifies that a specified resource exist within the Terraform state and in the data service.
func (idstc *ImagesDataSourceTestCase) Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[idstc.ResourceAddress]

		// Return an error if the resource is not found in state
		if !ok {
			return fmt.Errorf("Can't find images data source: %s", idstc.ResourceAddress)
		}

		// Ensure the Terraform state has set the resource ID
		if rs.Primary.ID == "" {
			return fmt.Errorf("Images data source ID not set")
		}

		// Execute the external API call to fetch images matching the query
		_images, err := TestProvider.client.GetImagesQuery(idstc.Query)

		// Propagate any errors from the external API call
		if err != nil {
			return err
		}

		// Update the provided images pointer with the fetched result
		idstc.Images = _images

		// Indicate successful existence check
		return nil
	}
}

// Attributes asserts that the number of images in the data source matches the expected count.
func (idstc *ImagesDataSourceTestCase) Attributes(expectedAmount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Determine the count of retrieved images
		amount := len(idstc.Images.Items)

		// Error if no images were retrieved
		if amount == 0 {
			return fmt.Errorf("%s has no images", idstc.ResourceAddress)
		}

		// Compare actual count with expected count
		if diff := deep.Equal(amount, expectedAmount); diff != nil {
			return fmt.Errorf("%s images amount does not match. Diff: %s", idstc.ResourceAddress, diff)
		}

		// All checks passed
		return nil
	}
}
