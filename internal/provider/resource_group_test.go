package cpln

import (
	"fmt"
	client "terraform-provider-cpln/internal/provider/client"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccControlPlaneGroup_basic(t *testing.T) {

	randomName := "group-" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, "GROUP") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckControlPlaneGroupCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccControlPlaneGroupWithJSMEPATH(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_group.tf-group", "name", randomName),
					resource.TestCheckResourceAttr("cpln_group.tf-group", "description", "group description "+randomName),
				),
			},
			{
				Config: testAccControlPlaneGroupWithJavaScript(randomName + "-javascript"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cpln_group.tf-group", "name", randomName+"-javascript"),
					resource.TestCheckResourceAttr("cpln_group.tf-group", "description", "group description "+randomName+"-javascript"),
				),
			},
		},
	})
}

func testAccControlPlaneGroupWithJSMEPATH(name string) string {

	return fmt.Sprintf(`
	
	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_service_account" "tf_sa" {

		name = "service-account-${var.random-name}"
		description = "service account description ${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}
	}

	resource "cpln_group" "tf-group" {

		depends_on = [cpln_service_account.tf_sa]

		name = var.random-name
		description = "group description ${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		// user_ids_and_emails = ["unittest@controlplane.com"]

		service_accounts = [cpln_service_account.tf_sa.name]

		member_query {

			fetch = "items"

			spec {
				match = "all"

				terms {
					op = "="
					tag = "firebase/sign_in_provider"
					value = "microsoft.com"
				}
			}
		}

		identity_matcher {
			expression = "groups"
			// language default value is 'jsmepath'
		}
	}
	`, name)
}

func testAccControlPlaneGroupWithJavaScript(name string) string {

	return fmt.Sprintf(`
	
	variable "random-name" {
		type = string
		default = "%s"
	}

	resource "cpln_service_account" "tf_sa" {

		name = "service-account-${var.random-name}"
		description = "service account description ${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}
	}

	resource "cpln_group" "tf-group" {

		depends_on = [cpln_service_account.tf_sa]

		name = var.random-name
		description = "group description ${var.random-name}" 
		
		tags = {
			terraform_generated = "true"
			acceptance_test = "true"
		}

		// user_ids_and_emails = ["unittest@controlplane.com"]

		service_accounts = [cpln_service_account.tf_sa.name]

		member_query {

			fetch = "items"

			spec {
				match = "all"

				terms {
					op = "="
					tag = "firebase/sign_in_provider"
					value = "microsoft.com"
				}
			}
		}

		identity_matcher {
			expression = "if ($.includes('groups')) { const y = $.groups; }"
			language = "javascript"
		}
	}
	`, name)
}

func testAccCheckControlPlaneGroupCheckDestroy(s *terraform.State) error {

	if len(s.RootModule().Resources) == 0 {
		return fmt.Errorf("Error In CheckDestroy. No Resources To Verify")
	}

	c := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "cpln_group" {

			groupName := rs.Primary.ID

			group, _, _ := c.GetGroup(groupName)
			if group != nil {
				return fmt.Errorf("Group still exists. Name: %s", *group.Name)
			}
		}

		if rs.Type == "cpln_service_account" {

			saName := rs.Primary.ID

			sa, _, _ := c.GetGroup(saName)
			if sa != nil {
				return fmt.Errorf("Service Account still exists. Name: %s", *sa.Name)
			}
		}
	}

	return nil
}

/*** Unit Tests ***/
// Build Functions //
func TestControlPlane_BuildGroupMemberLinks(t *testing.T) {

	u, sa := generateFlatTestMemberLinks()

	unitTestGroup := client.Group{}
	buildMemberLinks("testorg", u, sa, &unitTestGroup)

	if diff := deep.Equal(unitTestGroup, generateTestMemberLinks()); diff != nil {
		t.Errorf("Group Member Links was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildGroupQuery(t *testing.T) {

	unitTestGroup := client.Group{}
	unitTestGroup.MemberQuery = BuildQueryHelper("user", generateFlatTestGroupQuery())

	if diff := deep.Equal(&unitTestGroup, generateTestGroupQuery()); diff != nil {
		t.Errorf("Group Query was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildIdentityMatcher_WithJSMEPATH(t *testing.T) {
	identityMatcher, expectedIdentityMatcher, _ := generateTestIdentityMatcher("groups", "jsmepath")

	if diff := deep.Equal(identityMatcher, &expectedIdentityMatcher); diff != nil {
		t.Errorf("Identity Matcher was not built correctly. Diff: %s", diff)
	}
}

func TestControlPlane_BuildIdentityMatcher_WithJavaScript(t *testing.T) {
	identityMatcher, expectedIdentityMatcher, _ := generateTestIdentityMatcher("if ($.includes('groups')) { const y = $.groups; }", "javascript")

	if diff := deep.Equal(identityMatcher, &expectedIdentityMatcher); diff != nil {
		t.Errorf("Identity Matcher was not built correctly. Diff: %s", diff)
	}
}

// Generate Functions //
func generateTestMemberLinks() client.Group {

	testGroup := client.Group{}
	testGroup.MemberLinks = &[]string{
		"/org/testorg/user/username@cpln.io",
		"/org/testorg/user/control_plane_user",
		"/org/testorg/serviceaccount/terraform-service-account",
		"/org/testorg/serviceaccount/test-service-account",
	}

	return testGroup
}

func generateTestGroupQuery() *client.Group {

	testGroup := client.Group{}
	testGroup.MemberQuery = &client.Query{
		Kind:  GetString("user"),
		Fetch: GetString("items"),
	}

	testGroup.MemberQuery.Spec = &client.Spec{
		Match: GetString("all"),
		Terms: &[]client.Term{
			{
				Op:       GetString("="),
				Property: GetString("property"),
				// Rel:      GetString(""),
				// Tag:      GetString(""),
				Value: GetString("property-value"),
			},
			{
				Op: GetString("!="),
				// Property: GetString(""),
				Rel: GetString("rel"),
				// Tag:      GetString(""),
				Value: GetString("rel-value"),
			},
			{
				Op: GetString(">"),
				// Property: GetString(""),
				// Rel:      GetString(""),
				Tag:   GetString("tag"),
				Value: GetString("tag-value"),
			},
		},
	}

	return &testGroup
}

func generateTestIdentityMatcher(expression string, language string) (*client.IdentityMatcher, client.IdentityMatcher, []interface{}) {
	flattened := generateFlatTestIdentityMatcher(expression, language)
	identityMatcher := buildIdentityMatcher(flattened)
	expectedIdentityMatcher := client.IdentityMatcher{
		Expression: &expression,
		Language:   &language,
	}

	return identityMatcher, expectedIdentityMatcher, flattened
}

// Flatten Functions //
func TestControlPlane_FlattenMemberLinks(t *testing.T) {

	userIDs, serviceAccounts, err := flattenMemberLinks("testorg", generateTestMemberLinks().MemberLinks)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	userIDsFlat, serviceAccountsFlat := generateFlatTestMemberLinks()

	if diff := deep.Equal(userIDs, userIDsFlat.(*schema.Set).List()); diff != nil {
		t.Errorf("User IDs were not flattened correctly. Diff: %s", diff)
		return
	}

	if diff := deep.Equal(serviceAccounts, serviceAccountsFlat.(*schema.Set).List()); diff != nil {
		t.Errorf("Service Accounts were not flattened correctly. Diff: %s", diff)
	}
}

func TestControlPlane_FlattenQuery(t *testing.T) {

	query, err := FlattenQueryHelper(generateTestGroupQuery().MemberQuery)

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if diff := deep.Equal(query, generateFlatTestGroupQuery()); diff != nil {
		t.Errorf("Member Query was not flattened correctly. Diff: %s", diff)
	}
}

func generateFlatTestMemberLinks() (interface{}, interface{}) {

	usersFlat := []interface{}{
		"username@cpln.io",
		"control_plane_user",
	}

	serviceAccountsFlat := []interface{}{
		"test-service-account",
		"terraform-service-account",
	}

	stringFunc := schema.HashSchema(StringSchema())

	return schema.NewSet(stringFunc, usersFlat), schema.NewSet(stringFunc, serviceAccountsFlat)
}

func generateFlatTestGroupQuery() []interface{} {

	query := make(map[string]interface{})

	query["fetch"] = "items"

	spec := make(map[string]interface{})
	spec["match"] = "all"

	term01 := make(map[string]interface{})
	term01["op"] = "="
	term01["property"] = "property"
	// term01["rel"] = ""
	// term01["tag"] = ""
	term01["value"] = "property-value"

	term02 := make(map[string]interface{})
	term02["op"] = "!="
	// term02["property"] = ""
	term02["rel"] = "rel"
	// term02["tag"] = ""
	term02["value"] = "rel-value"

	term03 := make(map[string]interface{})
	term03["op"] = ">"
	// term03["property"] = ""
	// term03["rel"] = ""
	term03["tag"] = "tag"
	term03["value"] = "tag-value"

	terms := []interface{}{
		term01,
		term02,
		term03,
	}

	spec["terms"] = terms
	specArray := []interface{}{
		spec,
	}

	query["spec"] = specArray

	return []interface{}{
		query,
	}
}

func generateFlatTestIdentityMatcher(expression string, language string) []interface{} {
	spec := map[string]interface{}{
		"expression": expression,
		"language":   language,
	}

	return []interface{}{
		spec,
	}
}
