package cpln

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetBase - Set Base Helper
func SetBase(d *schema.ResourceData, base client.Base) error {

	if err := d.Set("cpln_id", base.ID); err != nil {
		return err
	}

	if err := d.Set("name", base.Name); err != nil {
		return err
	}

	if base.Description != nil {
		if err := d.Set("description", DescriptionHelper(*base.Name, *base.Description)); err != nil {
			return err
		}
	}

	if err := d.Set("tags", GetTags(base.Tags)); err != nil {
		return err
	}

	return nil
}

func FormatTypeToString(v interface{}) string {
	switch v := v.(type) {
	case float64, float32:
		return fmt.Sprintf("%.0f", v)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

func GetTags(tags *map[string]interface{}) map[string]interface{} {

	stringTypes := make(map[string]interface{})

	if tags != nil {
		for k, v := range *tags {

			// Remove certain server side generated tags
			if strings.HasPrefix(k, "cpln/deployTimestamp") || strings.HasPrefix(k, "cpln/aws") ||
				strings.HasPrefix(k, "cpln/azure") || strings.HasPrefix(k, "cpln/docker") ||
				strings.HasPrefix(k, "cpln/gcp") || strings.HasPrefix(k, "cpln/tls") ||
				strings.HasPrefix(k, "cpln/managedByTerraform") || strings.HasPrefix(k, "cpln/city") ||
				strings.HasPrefix(k, "cpln/continent") || strings.HasPrefix(k, "cpln/country") ||
				strings.HasPrefix(k, "cpln/state") {

				continue
			}

			stringTypes[k] = FormatTypeToString(v)
		}
	}

	return stringTypes
}

func GetTagChanges(d *schema.ResourceData) *map[string]interface{} {

	old, new := d.GetChange("tags")

	oldMap := map[string]interface{}{}

	for key, value := range old.(map[string]interface{}) {
		oldMap[key] = value
	}

	newMap := new.(map[string]interface{})

	for k := range oldMap {
		oldMap[k] = newMap[k]
	}

	for k := range newMap {
		oldMap[k] = newMap[k]
	}

	return &oldMap
}

func GetGVCEnvChanges(d *schema.ResourceData) *[]client.NameValue {
	_, new := d.GetChange("env")

	envArr := []client.NameValue{}
	for k, v := range new.(map[string]interface{}) {
		if v != nil {
			keyString := strings.Clone(k)
			valueString := v.(string)
			localEnvObj := client.NameValue{
				Name:  &keyString,
				Value: &valueString,
			}
			envArr = append(envArr, localEnvObj)
		}
	}

	return &envArr
}

func GetString(s interface{}) *string {

	if s == nil {
		return nil
	}

	output := s.(string)

	if strings.TrimSpace(output) == "" {
		return nil
	}

	return &output
}

func GetDescriptionString(s interface{}, name string) *string {

	if s == nil {
		return &name
	}

	output := s.(string)

	if strings.TrimSpace(output) == "" {
		return &name
	}

	return &output
}

func DescriptionHelper(name, description string) *string {

	descTrim := strings.TrimSpace(description)

	if descTrim == "" {
		descTrim = strings.TrimSpace(name)
	}

	return &descTrim
}

func DiffSuppressDescription(k, old, new string, d *schema.ResourceData) bool {

	if new == "" && old == d.State().ID {
		return true
	}

	return old == new
}

func GetInt(s interface{}) *int {
	if s == nil {
		return nil
	}

	output := s.(int)
	return &output
}

func GetFloat64(s interface{}) *float64 {

	if s == nil {
		return nil
	}

	output := s.(float64)
	return &output
}

func GetPortInt(s interface{}) *int {
	if s == nil {
		return nil
	}

	output := s.(int)

	if output == 0 {
		return nil
	} else {
		return &output
	}
}

func GetBool(s interface{}) *bool {
	if s == nil {
		return nil
	}

	output := s.(bool)
	return &output
}

func GetStringMap(s interface{}) *map[string]interface{} {

	if s == nil {
		return &map[string]interface{}{}
	}

	output := s.(map[string]interface{})
	return &output
}

func GetInterface(s interface{}) *interface{} {

	if s == nil {
		return nil
	}

	return &s
}

// MapSortHelper - Map Sort Helper
func MapSortHelper(i interface{}) ([]string, map[string]interface{}) {

	m := i.(map[string]interface{})

	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys, m
}

func NameValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	re := regexp.MustCompile(`^[a-z]([-a-z0-9])*[a-z0-9]$`)

	if !re.MatchString(v) {
		errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, v))
	}

	if len(v) > 63 {
		errs = append(errs, fmt.Errorf("%q length is > 63, got length: %d", key, len(v)))
	}

	return
}

func DescriptionValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	trim := strings.Trim(v, " ")

	if v != trim {
		errs = append(errs, fmt.Errorf("%q contains whitespace at the beginning or end, got: '%s'", key, v))
	}

	if len(v) > 250 {
		errs = append(errs, fmt.Errorf("%q length is > 250, got length: %d", key, len(v)))
	}

	return
}

func DescriptionDomainValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	trim := strings.Trim(v, " ")

	if v != trim {
		errs = append(errs, fmt.Errorf("%q contains whitespace at the beginning or end, got: '%s'", key, v))
	}

	// if len(v) > 250 {
	// 	errs = append(errs, fmt.Errorf("%q length is > 250, got length: %d", key, len(v)))
	// }

	return
}

func TagValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(map[string]interface{})

	if len(v) > 50 {
		errs = append(errs, fmt.Errorf("%q cannot have > 50 tags per object, got length: %d", key, len(v)))
	}

	return
}

func LinkValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	re := regexp.MustCompile(`(\/org\/[^/]+\/.*)|(\/\/.+)`)

	if !re.MatchString(v) {
		errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, v))
	}

	return
}

func KindValidator(val interface{}, key string) (warns []string, errs []error) {

	kind := val.(string)

	kinds := []string{
		"org",
		"cloudaccount",
		"policy",
		"user",
		"group",
		"resource",
		"task",
		"permissions",
		"serviceaccount",
		"secret",
		"location",
		"gvc",
		"workload",
		"quota",
		"identity",
		"deployment",
		"event",
		"domain",
		"image",
		"resourcepolicy",
		"agent",
		"accessreport",
		"policymembership",
		"auditctx",
	}

	for _, v := range kinds {
		if v == kind {
			return
		}
	}

	errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, kind))

	return
}

func AwsAccessKeyValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	re := regexp.MustCompile(`^AKIA.*`)

	if !re.MatchString(v) {
		errs = append(errs, fmt.Errorf("%q is invalid. must start with 'AKIA', got: %s", key, v))
		return
	}

	return
}

func EncodingValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	if v != "plain" && v != "base64" {
		errs = append(errs, fmt.Errorf("%q must be set to 'plain' or 'base64', got: %s", key, v))
	}

	return
}

func EmptyValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	if strings.TrimSpace(v) == "" {
		errs = append(errs, fmt.Errorf("%q must be must not be empty, got: %s", key, v))
	}

	return
}

func AwsRoleArnValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	re := regexp.MustCompile(`^arn:.*`)

	if !re.MatchString(v) {
		errs = append(errs, fmt.Errorf("%q is invalid. must start with 'arn:', got: %s", key, v))
		return
	}

	return
}

func PortValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(int)

	if v < 80 || v > 65535 {
		errs = append(errs, fmt.Errorf("%q must be between 80 and 65535 inclusive, got: %d", key, v))
	}

	return
}

func ObservabilityValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(int)

	if v < 0 {
		errs = append(errs, fmt.Errorf("%q must be >= 0, got: %d", key, v))
	}

	return
}

func CpuMemoryValidator(val interface{}, key string) (warns []string, errs []error) {

	v := val.(string)

	re := regexp.MustCompile(`^([+-]?[0-9.]+)([eEinumkKMGTP]*[-+]?[0-9]*)$`)

	if !re.MatchString(v) {
		errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, v))
		return
	}

	if len(v) > 20 {
		errs = append(errs, fmt.Errorf("%q cannot be greater than 20 characters, got: %d", key, len(v)))
	}

	return
}

func ThresholdValidator(val interface{}, key string) (warns []string, errs []error) {
	v := val.(int)
	if v < 1 || v > 20 {
		errs = append(errs, fmt.Errorf("%q must be between 1 and 20 inclusive, got: %d", key, v))
	}

	return
}

func QuerySchemaResource() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// "kind": {
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// },
			// "context": {
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// },
			"fetch": {
				Type:        schema.TypeString,
				Description: "Type of fetch. Specify either: `links` or `items`. Default: `items`.",
				Optional:    true,
				Default:     "items",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

					v := val.(string)

					if v != "items" && v != "links" {
						errs = append(errs, fmt.Errorf("%q must be either items or links. got: %s", key, v))
					}

					return
				},
			},
			"spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match": {
							Type:        schema.TypeString,
							Description: "Type of match. Available values: `all`, `any`, `none`. Default: `all`.",
							Optional:    true,
							Default:     "all",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

								v := val.(string)

								if v != "all" && v != "any" && v != "none" {
									errs = append(errs, fmt.Errorf("%q must be either 'all', 'any', 'none'. got: %s", key, v))
								}

								return
							},
						},
						"terms": {
							Type:        schema.TypeList,
							Description: "Terms can only contain one of the following attributes: `property`, `rel`, `tag`.",
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"op": {
										Type:        schema.TypeString,
										Description: "Type of query operation. Available values: `=`, `>`, `>=`, `<`, `<=`, `!=`, `exists`, `!exists`. Default: `=`.",
										Optional:    true,
										Default:     "=",
										ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {

											v := val.(string)

											if v != "=" && v != ">" && v != ">=" && v != "<" && v != "<=" && v != "!=" && v != "~" && v != "exists" && v != "!exists" {
												errs = append(errs, fmt.Errorf("%q must be either '=', '>', '>=', '<', '<=', '!=', '~', 'exists', '!exists'. got: %s", key, v))
											}

											return
										},
									},
									"property": {
										Type:        schema.TypeString,
										Description: "Property to use for query evaluation.",
										Optional:    true,
									},
									"rel": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tag": {
										Type:        schema.TypeString,
										Description: "Tag key to use for query evaluation.",
										Optional:    true,
									},
									"value": {
										Type:        schema.TypeString,
										Description: "Testing value for query evaluation.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func BuildQueryHelper(kind string, query interface{}) (queryObj *client.Query) {

	if query != nil {

		queryArray := query.([]interface{})

		if len(queryArray) > 0 {

			memberQueryMap := queryArray[0].(map[string]interface{})

			queryObj = &client.Query{}
			queryObj.Kind = GetString(kind)
			queryObj.Fetch = GetString(memberQueryMap["fetch"])

			specArray := memberQueryMap["spec"].([]interface{})

			if len(specArray) > 0 {

				specMap := specArray[0].(map[string]interface{})

				queryObj.Spec = &client.Spec{}
				queryObj.Spec.Match = GetString(specMap["match"])

				termsArray := specMap["terms"].([]interface{})

				if len(termsArray) > 0 {

					ct := []client.Term{}

					for _, t := range termsArray {

						term := t.(map[string]interface{})

						newTerm := client.Term{}

						newTerm.Op = GetString(term["op"])
						newTerm.Property = GetString(term["property"])
						newTerm.Rel = GetString(term["rel"])
						newTerm.Tag = GetString(term["tag"])
						newTerm.Value = GetString(term["value"])

						ct = append(ct, newTerm)
					}

					if len(ct) > 0 {
						queryObj.Spec.Terms = &ct
					}
				}

				return queryObj
			}
		}
	}

	return nil
}

func FlattenQueryHelper(query *client.Query) ([]interface{}, error) {

	if query == nil {
		return nil, nil
	}

	mq := make(map[string]interface{})

	// mq["kind"] = query.Kind

	if query.Fetch != nil {
		mq["fetch"] = *query.Fetch
	}

	if query.Spec != nil {

		spec := make(map[string]interface{})

		if query.Spec.Match != nil {
			spec["match"] = *query.Spec.Match
		}

		if query.Spec.Terms != nil && len(*query.Spec.Terms) > 0 {

			terms := []interface{}{}

			for _, term := range *query.Spec.Terms {

				t := make(map[string]interface{})

				if term.Op != nil {
					t["op"] = *term.Op
				}

				if term.Property != nil {
					t["property"] = *term.Property
				}

				if term.Rel != nil {
					t["rel"] = *term.Rel
				}

				if term.Tag != nil {
					t["tag"] = *term.Tag
				}

				if term.Value != nil {
					t["value"] = *term.Value
				}

				terms = append(terms, t)
			}

			spec["terms"] = terms
		}

		specArray := []interface{}{
			spec,
		}

		mq["spec"] = specArray
	}

	mqList := []interface{}{
		mq,
	}

	return mqList, nil
}

func BuildStringTypeSet(set interface{}) *[]string {

	if set == nil {
		return nil
	}

	output := []string{}

	for _, value := range set.(*schema.Set).List() {
		output = append(output, value.(string))
	}

	return &output
}

func FlattenStringTypeSet(set *[]string) []interface{} {

	if set == nil {
		return nil
	}

	output := []interface{}{}

	for _, value := range *set {
		output = append(output, value)
	}

	return output
}

func ConvertStringSliceToSet(slice []string) *schema.Set {

	var expanderSet []interface{}

	for _, value := range slice {
		expanderSet = append(expanderSet, value)
	}

	return schema.NewSet(schema.HashSchema(StringSchema()), expanderSet)
}

func ResourceExistsHelper() diag.Diagnostics {

	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "Resource already exists",
		Detail:   "Run 'terraform import' then run 'terraform apply' again.",
	})

	return diags
}

func SetSelfLink(links *[]client.Link, d *schema.ResourceData) error {

	if err := d.Set("self_link", GetSelfLink(links)); err != nil {
		return err
	}

	return nil
}

func GetSelfLink(links *[]client.Link) string {

	selfLink := ""

	if links != nil && len(*links) > 0 {
		for _, ls := range *links {
			if ls.Rel == "self" {
				selfLink = ls.Href
				break
			}
		}
	}

	return selfLink
}

func StringSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeString,
	}
}

func IntSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeInt,
	}
}

func WorkloadTypeValidator(val interface{}, key string) (warns []string, errs []error) {

	workloadType := val.(string)

	workloadTypes := []string{
		"serverless",
		"standard",
		"cron",
		"stateful",
	}

	for _, v := range workloadTypes {
		if v == workloadType {
			return
		}
	}

	errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, workloadType))

	return
}

func PortProtocolValidator(val interface{}, key string) (warns []string, errs []error) {

	portProtocol := val.(string)

	portProtocols := []string{
		"http",
		"http2",
		"grpc",
		"tcp",
	}

	for _, v := range portProtocols {
		if v == portProtocol {
			return
		}
	}

	errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, portProtocol))

	return
}

func GetNameFromSelfLink(selfLink string) string {
	parts := strings.Split(selfLink, "/")
	return parts[len(parts)-1]
}

func ExtractNumberAndCharactersFromString(value string) (int, string) {

	numberRegex := regexp.MustCompile("[0-9]+")
	charactersRegex := regexp.MustCompile("[A-Za-z]+")

	number, _ := strconv.Atoi(numberRegex.FindString(value))
	characters := charactersRegex.FindString(value)

	return number, characters
}

func MinValidator(min interface{}) schema.SchemaValidateFunc {
	return func(val interface{}, key string) (warns []string, errs []error) {

		// Convert values to float64 for comparison
		valFloat := toFloat64(val)
		minFloat := toFloat64(min)

		// Perform minimum value validation
		if valFloat < minFloat {
			errs = append(errs, fmt.Errorf("%s must be at least %v, got: %v", key, min, val))
		}

		return
	}
}

func RangeValidator(min, max interface{}) schema.SchemaValidateFunc {
	return func(val interface{}, key string) (warns []string, errs []error) {

		// Convert values to float64 for comparison
		valFloat := toFloat64(val)
		minFloat := toFloat64(min)
		maxFloat := toFloat64(max)

		// Perform range validation
		if valFloat < minFloat {
			errs = append(errs, fmt.Errorf("%s must be at least %v, got: %v", key, min, val))
		}
		if valFloat > maxFloat {
			errs = append(errs, fmt.Errorf("%s must be at most %v, got: %v", key, max, val))
		}

		return
	}
}

func RegexValidator(pattern string) schema.SchemaValidateFunc {
	return func(val interface{}, key string) (warns []string, errs []error) {

		v := val.(string)

		// Compile the regex pattern
		re, err := regexp.Compile(pattern)

		if err != nil {
			errs = append(errs, fmt.Errorf("invalid regex pattern for %s: %s", key, err))
			return
		}

		// Validate the value against the pattern
		if !re.MatchString(v) {
			errs = append(errs, fmt.Errorf("%s must match the pattern %s, got: %s", key, pattern, v))
		}

		return
	}
}

func InSliceValidator(validValues []string) schema.SchemaValidateFunc {
	return func(val interface{}, key string) (warns []string, errs []error) {
		v := val.(string)

		// Check if the value is in the validValues slice
		valid := false

		for _, allowed := range validValues {
			if v == allowed {
				valid = true
				break
			}
		}

		if !valid {
			errs = append(errs, fmt.Errorf("%q must be one of %v, got: %s", key, validValues, v))
		}

		return
	}
}

func CompositeValidator(base schema.SchemaValidateFunc, additional ...schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(val interface{}, key string) ([]string, []error) {
		var allWarnings []string
		var allErrors []error

		// Run base validation
		baseWarnings, baseErrors := base(val, key)
		allWarnings = append(allWarnings, baseWarnings...)
		allErrors = append(allErrors, baseErrors...)

		// Run additional validations
		for _, validator := range additional {
			warnings, errs := validator(val, key)
			allWarnings = append(allWarnings, warnings...)
			allErrors = append(allErrors, errs...)
		}

		return allWarnings, allErrors
	}
}

func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		return 0.0
	}
}

func toInterfaceSlice[T any](src []T) []interface{} {
	dst := make([]interface{}, len(src))

	for i, v := range src {
		dst[i] = v
	}

	return dst
}

func contains(opts []string, opt string) bool {
	for _, o := range opts {
		if o == opt {
			return true
		}
	}

	return false
}
