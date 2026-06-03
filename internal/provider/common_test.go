package cpln

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*** Unit Tests ***/

// TestFlattenLinkString exercises every relevant decision branch of the single-string link normalizer.
func TestFlattenLinkString(t *testing.T) {
	// Define the org name used across cases
	const org = "my-org"

	// Helper to construct a *string from a literal
	stringPtr := func(s string) *string { return &s }

	// Define the table of cases
	cases := []struct {
		name  string
		state types.String
		input *string
		want  types.String
	}{
		{
			name:  "nil input returns null",
			state: types.StringValue("//gvc/g/workload/w"),
			input: nil,
			want:  types.StringNull(),
		},
		{
			name:  "null state defers to API form",
			state: types.StringNull(),
			input: stringPtr("/org/my-org/gvc/g/workload/w"),
			want:  types.StringValue("/org/my-org/gvc/g/workload/w"),
		},
		{
			name:  "unknown state defers to API form",
			state: types.StringUnknown(),
			input: stringPtr("/org/my-org/gvc/g/workload/w"),
			want:  types.StringValue("/org/my-org/gvc/g/workload/w"),
		},
		{
			name:  "prior short form is preserved",
			state: types.StringValue("//gvc/g/workload/w"),
			input: stringPtr("/org/my-org/gvc/g/workload/w"),
			want:  types.StringValue("//gvc/g/workload/w"),
		},
		{
			name:  "prior long form is preserved",
			state: types.StringValue("/org/my-org/gvc/g/workload/w"),
			input: stringPtr("/org/my-org/gvc/g/workload/w"),
			want:  types.StringValue("/org/my-org/gvc/g/workload/w"),
		},
		{
			name:  "API value outside /org/<org>/ stays as-is",
			state: types.StringValue("//gvc/g/workload/w"),
			input: stringPtr("/org/other-org/gvc/g/workload/w"),
			want:  types.StringValue("/org/other-org/gvc/g/workload/w"),
		},
		{
			name:  "prior short form for different resource still resolves to API form",
			state: types.StringValue("//gvc/other/workload/other"),
			input: stringPtr("/org/my-org/gvc/g/workload/w"),
			want:  types.StringValue("/org/my-org/gvc/g/workload/w"),
		},
	}

	// Run each case
	for _, tc := range cases {
		// Run the case as a subtest
		t.Run(tc.name, func(t *testing.T) {
			// Invoke the helper
			got := FlattenLinkString(tc.state, tc.input, org)

			// Verify the returned value matches the expected value
			if !got.Equal(tc.want) {
				t.Fatalf("FlattenLinkString = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestFlattenLinkSet exercises every relevant decision branch of the set-of-links normalizer.
func TestFlattenLinkSet(t *testing.T) {
	// Define the org name used across cases
	const org = "my-org"

	// Helper to construct a *[]string from a slice literal
	stringSlicePtr := func(s []string) *[]string { return &s }

	// Helper to construct a types.Set of strings from literal values
	stringSet := func(values ...string) types.Set {
		// Allocate the elements slice
		elements := make([]attr.Value, 0, len(values))

		// Append each value as a types.String
		for _, v := range values {
			elements = append(elements, types.StringValue(v))
		}

		// Build the set and ignore diagnostics (test values are well-formed)
		s, _ := types.SetValue(types.StringType, elements)
		return s
	}

	// Define the table of cases
	cases := []struct {
		name  string
		state types.Set
		input *[]string
		want  types.Set
	}{
		{
			name:  "nil input returns null",
			state: stringSet("//gvc/g/workload/w"),
			input: nil,
			want:  types.SetNull(types.StringType),
		},
		{
			name:  "null state defers to API form for every entry",
			state: types.SetNull(types.StringType),
			input: stringSlicePtr([]string{"/org/my-org/gvc/g/workload/w"}),
			want:  stringSet("/org/my-org/gvc/g/workload/w"),
		},
		{
			name:  "prior short form is preserved when matching API entry is present",
			state: stringSet("//gvc/g/workload/w"),
			input: stringSlicePtr([]string{"/org/my-org/gvc/g/workload/w"}),
			want:  stringSet("//gvc/g/workload/w"),
		},
		{
			name:  "mixed prior forms preserve per-item choice",
			state: stringSet("//gvc/g/workload/a", "/org/my-org/gvc/g/workload/b"),
			input: stringSlicePtr([]string{"/org/my-org/gvc/g/workload/a", "/org/my-org/gvc/g/workload/b"}),
			want:  stringSet("//gvc/g/workload/a", "/org/my-org/gvc/g/workload/b"),
		},
		{
			name:  "API entry outside /org/<org>/ stays as-is",
			state: stringSet("//gvc/g/workload/w"),
			input: stringSlicePtr([]string{"/org/other-org/gvc/g/workload/w"}),
			want:  stringSet("/org/other-org/gvc/g/workload/w"),
		},
		{
			name:  "API entry without prior short form defaults to API form",
			state: stringSet("//gvc/g/workload/other"),
			input: stringSlicePtr([]string{"/org/my-org/gvc/g/workload/w"}),
			want:  stringSet("/org/my-org/gvc/g/workload/w"),
		},
		{
			name:  "empty API list yields empty set",
			state: stringSet("//gvc/g/workload/w"),
			input: stringSlicePtr([]string{}),
			want:  stringSet(),
		},
	}

	// Run each case
	for _, tc := range cases {
		// Run the case as a subtest
		t.Run(tc.name, func(t *testing.T) {
			// Allocate a fresh diagnostics container for this case
			diags := diag.Diagnostics{}

			// Invoke the helper
			got := FlattenLinkSet(&diags, tc.state, tc.input, org)

			// Verify that the helper did not raise any diagnostics
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			// Verify the returned value matches the expected value
			if !got.Equal(tc.want) {
				t.Fatalf("FlattenLinkSet = %v, want %v", got, tc.want)
			}
		})
	}
}
