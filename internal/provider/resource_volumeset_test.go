package cpln

import (
	"fmt"
	"strings"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneVolumeSet_basic(t *testing.T) {

	var volumeSet client.VolumeSet

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	gvcName := "gvc-volume-set-" + randomName
	volumeSetRequiredOnlyName := "volume-set-required-only-" + randomName
	volumeSetAllAttributesName := "volume-set-all-attributes-" + randomName
	description := "Volume Set description created using Terraform"

	// Update variables
	descriptionUpdated := "Volume Set description updated using Terraform"

	ep := resource.ExternalProvider{
		Source:            "time",
		VersionConstraint: "0.9.2",
	}

	eps := map[string]resource.ExternalProvider{
		"time": ep,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t, "VOLUME-SET") },
		Providers:         testAccProviders,
		ExternalProviders: eps,
		CheckDestroy:      testAccCheckControlPlaneVolumeSetCheckDestroy,
		Steps: []resource.TestStep{
			{
				// Required Only
				Config: testAccControlPlaneVolumeSet_requiredOnly(gvcName, volumeSetRequiredOnlyName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneVolumeSetExists("cpln_volume_set.new", gvcName, volumeSetRequiredOnlyName, &volumeSet),
					testAccCheckControlPlaneVolumeSetAttributes(&volumeSet, "new"),
					resource.TestCheckResourceAttr("cpln_volume_set.new", "name", volumeSetRequiredOnlyName),
					resource.TestCheckResourceAttr("cpln_volume_set.new", "description", description),
				),
			},
			{
				// Update Required Only
				Config: testAccControlPlaneVolumeSet_requiredOnlyUpdated(gvcName, volumeSetRequiredOnlyName, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneVolumeSetExists("cpln_volume_set.new", gvcName, volumeSetRequiredOnlyName, &volumeSet),
					testAccCheckControlPlaneVolumeSetAttributes(&volumeSet, "update"),
					resource.TestCheckResourceAttr("cpln_volume_set.new", "name", volumeSetRequiredOnlyName),
					resource.TestCheckResourceAttr("cpln_volume_set.new", "description", descriptionUpdated),
				),
			},
			{
				// All Attributes
				Config: testAccControlPlaneVolumeSet_allAttributes(gvcName, volumeSetAllAttributesName, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneVolumeSetExists("cpln_volume_set.new", gvcName, volumeSetAllAttributesName, &volumeSet),
					testAccCheckControlPlaneVolumeSetAttributes(&volumeSet, "all_attributes"),
					resource.TestCheckResourceAttr("cpln_volume_set.new", "name", volumeSetAllAttributesName),
					resource.TestCheckResourceAttr("cpln_volume_set.new", "description", descriptionUpdated),
				),
			},
			{
				// Update All Attributes
				Config: testAccControlPlaneVolumeSet_allAttributesUpdated(gvcName, volumeSetAllAttributesName, descriptionUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckControlPlaneVolumeSetExists("cpln_volume_set.new", gvcName, volumeSetAllAttributesName, &volumeSet),
					testAccCheckControlPlaneVolumeSetAttributes(&volumeSet, "update_all_attributes"),
					resource.TestCheckResourceAttr("cpln_volume_set.new", "name", volumeSetAllAttributesName),
					resource.TestCheckResourceAttr("cpln_volume_set.new", "description", descriptionUpdated),
				),
			},
		},
	})
}

func testAccControlPlaneVolumeSet_requiredOnly(gvcName string, name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_gvc" "new" {

		name        = "%s"
		description = "This is a GVC description"
	  
		locations = ["aws-eu-central-1", "aws-us-west-2"]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
		}
	}

	resource "time_sleep" "wait_30_seconds" {
		depends_on = [cpln_gvc.new]
		destroy_duration = "30s"
    }
	
	resource "cpln_volume_set" "new" {

		depends_on = [time_sleep.wait_30_seconds]

		name 		= "%s"
		description = "%s"

		tags = {
			terraform_generated = "true"
			acceptance_test     = "true"
		}

		gvc 			  = cpln_gvc.new.name
		initial_capacity  = 10
		performance_class = "general-purpose-ssd"
		file_system_type  = "ext4"
	}
	
	`, gvcName, name, description)
}

func testAccControlPlaneVolumeSet_requiredOnlyUpdated(gvcName string, name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_gvc" "new" {

		depends_on = [time_sleep.wait_30_seconds]

		name        = "%s"
		description = "This is a GVC description"
	  
		locations = ["aws-eu-central-1", "aws-us-west-2"]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
		}
	}
	
	resource "time_sleep" "wait_30_seconds" {
		depends_on = [cpln_gvc.new]
		destroy_duration = "30s"
    }
	
	resource "cpln_volume_set" "new" {

		depends_on = [time_sleep.wait_30_seconds]
		
		name 		= "%s"
		description = "%s"

		tags = {
			terraform_generated = "true"
			acceptance_test     = "true"
			update              = "true"
		}

		gvc 			  = cpln_gvc.new.name
		initial_capacity  = 15
		performance_class = "general-purpose-ssd"
		file_system_type  = "ext4"
	}
	
	`, gvcName, name, description)
}

func testAccControlPlaneVolumeSet_allAttributes(gvcName string, name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_gvc" "new" {

		depends_on = [time_sleep.wait_30_seconds]

		name        = "%s"
		description = "This is a GVC description"
	  
		locations = ["aws-eu-central-1", "aws-us-west-2"]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
		}
	}

	resource "time_sleep" "wait_30_seconds" {
		depends_on = [cpln_gvc.new]
		destroy_duration = "30s"
    }
	
	resource "cpln_volume_set" "new" {

		depends_on = [time_sleep.wait_30_seconds]
		
		name 		= "%s"
		description = "%s"

		tags = {
			terraform_generated = "true"
			acceptance_test     = "true"
			update              = "true"
		}

		gvc 			  = cpln_gvc.new.name
		initial_capacity  = 1000
		performance_class = "high-throughput-ssd"
		file_system_type  = "xfs"

		snapshots {
			create_final_snapshot = true
			retention_duration    = "1d"
		}

		autoscaling {
			max_capacity        = 2000
			min_free_percentage = 1
			scaling_factor      = 1.1
		}
	}
	
	`, gvcName, name, description)
}

func testAccControlPlaneVolumeSet_allAttributesUpdated(gvcName string, name string, description string) string {

	return fmt.Sprintf(`

	resource "cpln_gvc" "new" {

		depends_on = [time_sleep.wait_30_seconds]

		name        = "%s"
		description = "This is a GVC description"
	  
		locations = ["aws-eu-central-1", "aws-us-west-2"]
	  
		tags = {
		  terraform_generated = "true"
		  acceptance_test     = "true"
		}
	}

	resource "time_sleep" "wait_30_seconds" {
		depends_on = [cpln_gvc.new]
		destroy_duration = "30s"
    }
	
	resource "cpln_volume_set" "new" {

		depends_on = [time_sleep.wait_30_seconds]

	resource "cpln_volume_set" "new" {
		
		name 		= "%s"
		description = "%s"

		tags = {
			terraform_generated = "true"
			acceptance_test     = "true"
			update              = "true"
		}

		gvc 			  = cpln_gvc.new.name
		initial_capacity  = 1010
		performance_class = "high-throughput-ssd"
		file_system_type  = "xfs"

		snapshots {
			create_final_snapshot = false
			retention_duration    = "2d"
		}

		autoscaling {
			max_capacity        = 2000
			min_free_percentage = 2
			scaling_factor      = 2.2
		}
	}
	
	`, gvcName, name, description)
}

func testAccCheckControlPlaneVolumeSetExists(resourceName string, gvcName string, volumeSetName string, volumeSet *client.VolumeSet) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		TestLogger.Printf("Inside testAccCheckControlPlaneVolumeSetExists. Resources Length: %d", len(s.RootModule().Resources))

		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("Not found: %s", s)
		}

		if rs.Primary.ID != volumeSetName {
			return fmt.Errorf("Volume Set name does not match")
		}

		client := testAccProvider.Meta().(*client.Client)
		_volumeSet, _, err := client.GetVolumeSet(volumeSetName, gvcName)

		if err != nil {
			return err
		}

		if *_volumeSet.Name != volumeSetName {
			return fmt.Errorf("Volume Set name does not match")
		}

		*volumeSet = *_volumeSet

		return nil
	}
}

func testAccCheckControlPlaneVolumeSetAttributes(volumeSet *client.VolumeSet, state string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		tags := *volumeSet.Tags

		if tags["terraform_generated"] != "true" {
			return fmt.Errorf("Volume Set Tags - `terraform_generated` attribute does not match")
		}

		if tags["acceptance_test"] != "true" {
			return fmt.Errorf("Volume Set Tags - `acceptance_test` attribute does not match")
		}

		expectedVolumeSetSpec := generateTestVolumeSetSpec(state)

		// Check Spec
		if diff := deep.Equal(volumeSet.Spec, expectedVolumeSetSpec); diff != nil {
			return fmt.Errorf("Volume Set Spec does not match. Diff: %s", diff)
		}

		return nil
	}
}

func testAccCheckControlPlaneVolumeSetCheckDestroy(s *terraform.State) error {

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error in CheckDestroy for Volume Set: No resources to verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "cpln_gvc" {
			continue
		}

		gvcName := rs.Primary.ID
		gvc, _, _ := c.GetGvc(gvcName)

		if gvc != nil {
			return fmt.Errorf("GVC still exists. Name: %s. Associated Volume Sets might still exist", *gvc.Name)
		}
	}

	return nil
}

/*** Unit Tests ***/
// Build //
func TestControlPlane_BuildVolumeSetSnapshots(t *testing.T) {

	snapshots, expectedSnapshots, _ := generateTestVolumeSetSnapshots("new")

	if diff := deep.Equal(snapshots, expectedSnapshots); diff != nil {
		t.Errorf("Volume Set Snapshots was not built correctly, Diff: %s", diff)
	}
}

func TestControlPlane_BuildVolumeSetScaling(t *testing.T) {

	autoscaling, expectedAutoscaling, _ := generateTestVolumeSetAutoscaling("new")

	if diff := deep.Equal(autoscaling, expectedAutoscaling); diff != nil {
		t.Errorf("Volume Set Autoscaling was not built correctly, Diff: %s", diff)
	}
}

// Flatten //
func TestControlPlane_FlattenVolumeSetSnapshots(t *testing.T) {

	_, expectedSnapshots, expectedFlatten := generateTestVolumeSetSnapshots("new")
	flattenedSnapshots := flattenVolumeSetSnapshots(expectedSnapshots)

	if diff := deep.Equal(expectedFlatten, flattenedSnapshots); diff != nil {
		t.Errorf("Volume Set Snapshots was not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenVolumeSetAutoscaling(t *testing.T) {

	_, expectedAutoscaling, expectedFlatten := generateTestVolumeSetAutoscaling("new")
	flattenedAutoscaling := flattenVolumeSetAutoscaling(expectedAutoscaling)

	if diff := deep.Equal(expectedFlatten, flattenedAutoscaling); diff != nil {
		t.Errorf("Volume Set Autoscaling was not flattened correctly. Diff: %s", diff)
	}
}

/*** Generate ***/
// Build //
func generateTestVolumeSetSpec(state string) *client.VolumeSetSpec {

	isAllAttributes := strings.Contains(state, "attributes")
	initialCapacity := 10
	performanceClass := "general-purpose-ssd"
	fileSystemType := "ext4"

	if isAllAttributes {
		initialCapacity = 1000
		performanceClass = "high-throughput-ssd"
		fileSystemType = "xfs"

	}

	if strings.Contains(state, "update") {
		initialCapacity = 15

		if isAllAttributes {
			initialCapacity = 1010
		}
	}

	spec := client.VolumeSetSpec{
		InitialCapacity:  &initialCapacity,
		PerformanceClass: &performanceClass,
		FileSystemType:   &fileSystemType,
	}

	if strings.Contains(state, "attributes") {
		snapshots, _, _ := generateTestVolumeSetSnapshots(state)
		spec.Snapshots = snapshots
	}

	if strings.Contains(state, "attributes") {
		autoscaling, _, _ := generateTestVolumeSetAutoscaling(state)
		spec.AutoScaling = autoscaling
	}

	return &spec
}

func generateTestVolumeSetSnapshots(state string) (*client.VolumeSetSnapshots, *client.VolumeSetSnapshots, []interface{}) {

	createFinalSnapshot := true
	retentionDuration := "1d"

	if strings.Contains(state, "update") {
		createFinalSnapshot = false
		retentionDuration = "2d"
	}

	flattened := generateFlatTestVolumeSetSnapshots(createFinalSnapshot, retentionDuration)
	snapshots := buildVolumeSetSnapshots(flattened)
	expectedSnapshot := client.VolumeSetSnapshots{
		CreateFinalSnapshot: &createFinalSnapshot,
		RetentionDuration:   &retentionDuration,
	}

	return snapshots, &expectedSnapshot, flattened
}

func generateTestVolumeSetAutoscaling(state string) (*client.VolumeSetScaling, *client.VolumeSetScaling, []interface{}) {

	maxCapacity := 2000
	minFreePercentage := 1
	scalingFactor := 1.1

	if strings.Contains(state, "update") {
		maxCapacity = 2000
		minFreePercentage = 2
		scalingFactor = 2.2
	}

	flattened := generateFlatTestVolumeSetAutoscaling(maxCapacity, minFreePercentage, scalingFactor)
	autoscaling := buildVolumeSetAutoscaling(flattened)
	expectedAutoscaling := client.VolumeSetScaling{
		MaxCapacity:       &maxCapacity,
		MinFreePercentage: &minFreePercentage,
		ScalingFactor:     &scalingFactor,
	}

	return autoscaling, &expectedAutoscaling, flattened
}

// Flatten //
func generateFlatTestVolumeSetSnapshots(createFinalSnapshot bool, retentionDuration string) []interface{} {

	spec := map[string]interface{}{
		"create_final_snapshot": createFinalSnapshot,
		"retention_duration":    retentionDuration,
	}

	return []interface{}{
		spec,
	}
}

func generateFlatTestVolumeSetAutoscaling(maxCapacity int, minFreePercentage int, scalingFactor float64) []interface{} {

	spec := map[string]interface{}{
		"max_capacity":        maxCapacity,
		"min_free_percentage": minFreePercentage,
		"scaling_factor":      scalingFactor,
	}

	return []interface{}{
		spec,
	}
}
