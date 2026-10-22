package cpln

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*** Unit Tests ***/

// TestBuildMapJSON verifies which string values are sent to the API as decoded JSON and which stay strings.
func TestBuildMapJSON(t *testing.T) {
	// Define the table of cases
	cases := []struct {
		name  string
		input string
		want  interface{}
	}{
		{name: "integer", input: "2500", want: float64(2500)},
		{name: "fraction", input: "0.5", want: 0.5},
		{name: "boolean", input: "true", want: true},
		{name: "list", input: `["__name__"]`, want: []interface{}{"__name__"}},
		{name: "object", input: `{"values":["a"]}`, want: map[string]interface{}{"values": []interface{}{"a"}}},
		{name: "plain string", input: "go_.*", want: "go_.*"},
		{name: "JSON string stays quoted", input: `"quoted"`, want: `"quoted"`},
		{name: "null stays a string", input: "null", want: "null"},
		{name: "leading space stays a string", input: " 2500", want: " 2500"},
		{name: "non canonical number stays a string", input: "1.0", want: "1.0"},
		{name: "unsorted keys stay a string", input: `{"b":1,"a":2}`, want: `{"b":1,"a":2}`},
	}

	// Iterate over each case
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build the map from a single entry
			var diags diag.Diagnostics
			output := BuildMapJSON(context.Background(), &diags, types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue(tc.input)}))

			// Fail on diagnostics or a nil result
			if diags.HasError() || output == nil {
				t.Fatalf("unexpected diagnostics or nil output: %v", diags)
			}

			// Compare the decoded value
			if got := (*output)["key"]; !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
		})
	}
}

// TestMapJSONRoundTrip verifies that every value comes back unchanged after passing through the API as JSON.
func TestMapJSONRoundTrip(t *testing.T) {
	// Define the values a configuration can hold
	inputs := []string{"2500", "0.5", "false", `["__name__"]`, `{"values":["a"]}`, `{"a":"<b>"}`, "go_.*", "", " 2500", "1.0", "null", `"quoted"`, `{"b":1,"a":2}`}

	// Iterate over each value
	for _, input := range inputs {
		// Build the API map from a single entry
		var diags diag.Diagnostics
		built := BuildMapJSON(context.Background(), &diags, types.MapValueMust(types.StringType, map[string]attr.Value{"key": types.StringValue(input)}))

		// Fail on diagnostics
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics for %q: %v", input, diags)
		}

		// Send the map through JSON the way the API client does
		encoded, err := json.Marshal(built)
		if err != nil {
			t.Fatalf("marshal %q: %v", input, err)
		}

		// Decode the API response into a fresh map
		var received map[string]interface{}
		if err := json.Unmarshal(encoded, &received); err != nil {
			t.Fatalf("unmarshal %q: %v", input, err)
		}

		// Flatten the received map back into state
		flattened := FlattenMapJSON(&received).Elements()["key"]

		// Compare against the original configuration value
		if !flattened.Equal(types.StringValue(input)) {
			t.Errorf("round trip of %q produced %s", input, flattened)
		}
	}
}
