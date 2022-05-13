package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"sort"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func generateFlatTestTargetLinks() []interface{} {

	targetLinksFlat := []interface{}{
		"secret-test-02",
		"secret-test-01",
	}

	return targetLinksFlat
}

func generateTestTargetLinks() *client.Policy {

	testPolicy := client.Policy{}
	testPolicy.TargetLinks = &[]string{
		"/org/testorg/secret/secret-test-02",
		"/org/testorg/secret/secret-test-01",
	}

	return &testPolicy
}

func TestControlPlane_BuildPolicyTargetLinks(t *testing.T) {

	tlf := generateFlatTestTargetLinks()
	stringFunc := schema.HashSchema(StringSchema())
	tlfSet := schema.NewSet(stringFunc, tlf)

	unitTestPolicy := &client.Policy{}
	buildTargetLinks("testorg", "secret", tlfSet, unitTestPolicy)

	if diff := deep.Equal(unitTestPolicy, generateTestTargetLinks()); diff != nil {
		t.Errorf("Policy Target Links was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenPolicyTargetLinks(t *testing.T) {

	ftl := flattenTargetLinks(generateTestTargetLinks().TargetLinks)
	tlf := generateFlatTestTargetLinks()

	if diff := deep.Equal(ftl, tlf); diff != nil {
		t.Errorf("Target links were not flattened correctly. Diff: %s", diff)
		return
	}
}

func generateFlatTestBindings() *schema.Set {

	stringFunc := schema.HashSchema(StringSchema())

	b1 := make(map[string]interface{})

	permSet1 := schema.NewSet(stringFunc, []interface{}{
		"manage",
		"edit",
	})

	b1["permissions"] = permSet1

	pLinkSet1 := schema.NewSet(stringFunc, []interface{}{
		"user/support@controlplane.com",
		"serviceaccount/support",
	})

	b1["principal_links"] = pLinkSet1

	b2 := make(map[string]interface{})

	permSet2 := schema.NewSet(stringFunc, []interface{}{
		"viewer",
	})

	b2["permissions"] = permSet2

	pLinkSet2 := schema.NewSet(stringFunc, []interface{}{
		"group/admins",
		"serviceaccount/tester",
	})

	b2["principal_links"] = pLinkSet2

	sFunc := schema.HashResource(BindingResource())

	return schema.NewSet(sFunc, []interface{}{b1, b2})
}

func generateTestBindings() *client.Policy {

	testPolicy := client.Policy{}

	binding_01 := client.Binding{
		Permissions: &[]string{
			"manage",
			"edit",
		},
		PrincipalLinks: &[]string{
			"/org/testorg/user/support@controlplane.com",
			"/org/testorg/serviceaccount/support",
		},
	}

	binding_02 := client.Binding{
		Permissions: &[]string{
			"viewer",
		},
		PrincipalLinks: &[]string{
			"/org/testorg/group/admins",
			"/org/testorg/serviceaccount/tester",
		},
	}

	testPolicy.Bindings = &[]client.Binding{
		binding_01,
		binding_02,
	}

	return &testPolicy
}

func TestControlPlane_BuildPolicyBindings(t *testing.T) {

	unitTestPolicy := &client.Policy{}
	buildBindings("testorg", generateFlatTestBindings(), unitTestPolicy)

	generatedBindings := generateTestBindings()

	if len(*unitTestPolicy.Bindings) != len(*generatedBindings.Bindings) {
		t.Error("Policy Bindings was not built correctly. Different binding lengths")
	}

	sortInternalBindings(unitTestPolicy.Bindings)
	sortInternalBindings(generatedBindings.Bindings)

	match := false

	for _, b1 := range *unitTestPolicy.Bindings {

		match = false

		for _, b2 := range *generatedBindings.Bindings {

			if diff := deep.Equal(b1, b2); diff == nil {
				match = true
				break
			}
		}

		if !match {
			break
		}
	}

	if !match {
		t.Error("Policy Bindings was not built correctly.")
	}
}

func sortInternalBindings(binding *[]client.Binding) {
	for _, b := range *binding {
		sort.Strings(*b.Permissions)
		sort.Strings(*b.PrincipalLinks)
	}
}

func TestControlPlane_FlattenPolicyBindings(t *testing.T) {

	fb := flattenBindings("testorg", generateTestBindings().Bindings)
	tb := generateFlatTestBindings().List()

	if len(fb) != len(tb) {
		t.Error("Policy Bindings was not built correctly. Different binding lengths")
	}

	sortInternalBindingsInterface(fb)
	sortInternalBindingsSet(tb)

	match := false

	for _, b1 := range fb {

		match = false

		for _, b2 := range tb {

			if diff := deep.Equal(b1, b2); diff == nil {
				match = true
				break
			}
		}

		if !match {
			break
		}
	}

	if !match {
		t.Error("Policy Bindings was not flattened correctly.")
	}
}

func sortInternalBindingsInterface(binding []interface{}) {
	for _, b := range binding {
		b1 := b.(map[string]interface{})
		b1["permissions"] = sortInterfaceStrings(b1["permissions"].([]interface{}))
		b1["principal_links"] = sortInterfaceStrings(b1["principal_links"].([]interface{}))
	}
}

func sortInternalBindingsSet(binding []interface{}) {
	for _, b := range binding {
		b1 := b.(map[string]interface{})
		b1["permissions"] = sortInterfaceStrings(b1["permissions"].(*schema.Set).List())
		b1["principal_links"] = sortInterfaceStrings(b1["principal_links"].(*schema.Set).List())
	}
}

func sortInterfaceStrings(input []interface{}) []interface{} {

	s := make([]string, len(input))
	for i, v := range input {
		s[i] = v.(string)
	}

	sort.Strings(s)

	o := []interface{}{}

	for _, p := range s {
		o = append(o, p)
	}

	return o
}

func TestAccControlPlanePolicy_basic(t *testing.T) {

	randomName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "POLICY") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlanePolicyCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlanePolicy(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_policy.terraform_policy", "name", "policy-"+randomName),
					resource.TestCheckResourceAttr("cpln_policy.terraform_policy", "description", "Policy description for policy-"+randomName),
				),
			},
		},
	})
}

func testAccControlPlanePolicy(name string) string {

	return fmt.Sprintf(`
	
	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_gvc" "terraform_gvc" {
	
		name        = "gvc-${var.random-name}"	
		description = "GVC description for gvc-${var.random-name}"

		locations = ["aws-eu-central-1"]

		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}
	}

	resource "cpln_identity" "terraform_identity" {

  		gvc = cpln_gvc.terraform_gvc.name

		name        = "identity-${var.random-name}"	
		description = "Identity description for identity-${var.random-name}"
 
		tags = {
		  terraform_generated = "true"
		  acceptance_test = "true"
		}
	}


	resource "cpln_policy" "terraform_policy" {

		name = "policy-${var.random-name}"
		description = "Policy description for policy-${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		target_kind = "secret"
		target = "all"

	}

	resource "cpln_service_account" "tf_sa" {

		name = "service-account-${var.random-name}"
		description = "service account description ${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}
	}

	resource "cpln_policy" "terraform_policy_01" {

		name = "policy-01-${var.random-name}"
		description = "Policy description for policy-01-${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}
	
		target_kind = "serviceaccount"
		target_links = [cpln_service_account.tf_sa.name]
		// target = "all"

		target_query {
		
			spec {
				# match is either "all", "any", or "none"
				match = "all"

				terms {
					op = "="
					tag = "firebase/sign_in_provider"
					value = "microsoft.com"
				}
			}
		}

		binding {
			permissions = ["manage", "edit"]
			principal_links = ["user/support@controlplane.com", "group/viewers", "serviceaccount/service-account-${var.random-name}","gvc/${cpln_gvc.terraform_gvc.name}/identity/${cpln_identity.terraform_identity.name}"]
		}
	}

	`, name)
}

func testAccCheckControlPlanePolicyCheckDestroy(s *terraform.State) error {

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy For Policy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "cpln_policy" {
			policy, _, _ := c.GetPolicy(rs.Primary.ID)
			if policy != nil {
				return fmt.Errorf("Policy still exists. Name: %s", *policy.Name)
			}
		}
	}

	return nil
}
