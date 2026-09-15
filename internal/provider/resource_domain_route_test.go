package cpln

import (
	"testing"

	client "github.com/controlplane-com/terraform-provider-cpln/internal/provider/client"
)

/*** Unit Tests ***/

// The acceptance coverage for cpln_domain_route lives alongside the domain it belongs to, in resource_domain_test.go.

// TestDomainRouteIdentityKey exercises every relevant decision branch of the route identity key.
func TestDomainRouteIdentityKey(t *testing.T) {
	// Helper to construct a *string from a literal
	stringPtr := func(s string) *string { return &s }

	// Define the table of cases
	cases := []struct {
		name  string
		route client.DomainRoute
		want  string
	}{
		{
			name:  "route without a path has no identity",
			route: client.DomainRoute{HostPrefix: stringPtr("app.")},
			want:  "",
		},
		{
			name:  "prefix without a host",
			route: client.DomainRoute{Prefix: stringPtr("/")},
			want:  "prefix:/;",
		},
		{
			name:  "regex without a host",
			route: client.DomainRoute{Regex: stringPtr("^/user/.*$")},
			want:  "regex:^/user/.*$;",
		},
		{
			name:  "prefix narrowed to a host prefix",
			route: client.DomainRoute{Prefix: stringPtr("/"), HostPrefix: stringPtr("app.")},
			want:  "prefix:/;hostPrefix:app.",
		},
		{
			name:  "prefix narrowed to a host regex",
			route: client.DomainRoute{Prefix: stringPtr("/"), HostRegex: stringPtr("^app.*$")},
			want:  "prefix:/;hostRegex:^app.*$",
		},
		{
			name:  "regex narrowed to a host prefix",
			route: client.DomainRoute{Regex: stringPtr("^/api$"), HostPrefix: stringPtr("app.")},
			want:  "regex:^/api$;hostPrefix:app.",
		},
		{
			name:  "attributes outside the identity are ignored",
			route: client.DomainRoute{Prefix: stringPtr("/"), WorkloadLink: stringPtr("//gvc/g/workload/w"), Replica: IntPointer(2)},
			want:  "prefix:/;",
		},
	}

	// Run each case
	for _, tc := range cases {
		// Run the case as a subtest
		t.Run(tc.name, func(t *testing.T) {
			// Invoke the method under test
			got := tc.route.IdentityKey()

			// Verify the returned key matches the expected key
			if got != tc.want {
				t.Fatalf("IdentityKey = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestDomainRouteIdentityKeyDistinguishesHosts verifies that routes sharing a path on different hosts never collide.
func TestDomainRouteIdentityKeyDistinguishesHosts(t *testing.T) {
	// Helper to construct a *string from a literal
	stringPtr := func(s string) *string { return &s }

	// Build the route set a wildcard domain carries, where every route shares the same path
	routes := []client.DomainRoute{
		{Prefix: stringPtr("/"), HostPrefix: stringPtr("anntaylor.")},
		{Prefix: stringPtr("/"), HostPrefix: stringPtr("simplybe.")},
		{Prefix: stringPtr("/"), HostPrefix: stringPtr("loft.")},
		{Prefix: stringPtr("/"), HostRegex: stringPtr("^shop.*$")},
		{Prefix: stringPtr("/")},
	}

	// Collect the identity key of every route
	seen := map[string]int{}

	for index, route := range routes {
		// Fail when two routes in the set produce the same key
		if previous, exists := seen[route.IdentityKey()]; exists {
			t.Fatalf("routes %d and %d share the identity key %q", previous, index, route.IdentityKey())
		}

		seen[route.IdentityKey()] = index
	}
}

// TestBuildDomainRouteID exercises every relevant decision branch of the route identifier builder.
func TestBuildDomainRouteID(t *testing.T) {
	// Define the domain link used across cases
	const domainLink = "/org/my-org/domain/example.com"

	// Helper to construct a *string from a literal
	stringPtr := func(s string) *string { return &s }

	// Define the table of cases
	cases := []struct {
		name  string
		port  int32
		route client.DomainRoute
		want  string
	}{
		{
			name:  "prefix without a host keeps the identifier earlier versions issued",
			port:  443,
			route: client.DomainRoute{Prefix: stringPtr("/")},
			want:  "/org/my-org/domain/example.com_443_/",
		},
		{
			name:  "regex without a host keeps the identifier earlier versions issued",
			port:  80,
			route: client.DomainRoute{Regex: stringPtr("^/user/.*$")},
			want:  "/org/my-org/domain/example.com_80_^/user/.*$",
		},
		{
			name:  "host prefix is carried in the identifier",
			port:  443,
			route: client.DomainRoute{Prefix: stringPtr("/"), HostPrefix: stringPtr("app.")},
			want:  "/org/my-org/domain/example.com_443_/_app.",
		},
		{
			name:  "host regex is carried in the identifier",
			port:  443,
			route: client.DomainRoute{Prefix: stringPtr("/"), HostRegex: stringPtr("^app.*$")},
			want:  "/org/my-org/domain/example.com_443_/_^app.*$",
		},
	}

	// Run each case
	for _, tc := range cases {
		// Run the case as a subtest
		t.Run(tc.name, func(t *testing.T) {
			// Invoke the helper
			got := buildDomainRouteID(domainLink, tc.port, tc.route)

			// Verify the returned identifier matches the expected identifier
			if got != tc.want {
				t.Fatalf("buildDomainRouteID = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestBuildImportIdentifier exercises the import identifier the provider hands back in its conflict diagnostics.
func TestBuildImportIdentifier(t *testing.T) {
	// Define the domain link used across cases
	const domainLink = "/org/my-org/domain/example.com"

	// Helper to construct a *string from a literal
	stringPtr := func(s string) *string { return &s }

	// Define the table of cases
	cases := []struct {
		name  string
		route client.DomainRoute
		want  string
	}{
		{
			name:  "a route that matches every host ends in a bare host segment",
			route: client.DomainRoute{Prefix: stringPtr("/")},
			want:  "/org/my-org/domain/example.com:443:/:",
		},
		{
			name:  "a route narrowed to a host prefix names it",
			route: client.DomainRoute{Prefix: stringPtr("/"), HostPrefix: stringPtr("store.")},
			want:  "/org/my-org/domain/example.com:443:/:store.",
		},
		{
			name:  "a regex route names the regex",
			route: client.DomainRoute{Regex: stringPtr("^/api$"), HostRegex: stringPtr("^app.*$")},
			want:  "/org/my-org/domain/example.com:443:^/api$:^app.*$",
		},
	}

	// Run each case
	for _, tc := range cases {
		// Run the case as a subtest
		t.Run(tc.name, func(t *testing.T) {
			// Invoke the helper
			got := buildImportIdentifier(domainLink, 443, tc.route)

			// Verify the returned identifier matches the expected identifier
			if got != tc.want {
				t.Fatalf("buildImportIdentifier = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestDomainRouteIdentifier exercises the human readable description used in route error messages.
func TestDomainRouteIdentifier(t *testing.T) {
	// Helper to construct a *string from a literal
	stringPtr := func(s string) *string { return &s }

	// Define the table of cases
	cases := []struct {
		name  string
		route client.DomainRoute
		want  string
	}{
		{
			name:  "prefix only",
			route: client.DomainRoute{Prefix: stringPtr("/api")},
			want:  "with prefix '/api'",
		},
		{
			name:  "regex only",
			route: client.DomainRoute{Regex: stringPtr("^/api$")},
			want:  "with regex '^/api$'",
		},
		{
			name:  "prefix narrowed to a host prefix",
			route: client.DomainRoute{Prefix: stringPtr("/"), HostPrefix: stringPtr("app.")},
			want:  "with prefix '/' and host prefix 'app.'",
		},
		{
			name:  "prefix narrowed to a host regex",
			route: client.DomainRoute{Prefix: stringPtr("/"), HostRegex: stringPtr("^app.*$")},
			want:  "with prefix '/' and host regex '^app.*$'",
		},
	}

	// Run each case
	for _, tc := range cases {
		// Run the case as a subtest
		t.Run(tc.name, func(t *testing.T) {
			// Invoke the method under test
			got := tc.route.Identifier()

			// Verify the returned description matches the expected description
			if got != tc.want {
				t.Fatalf("Identifier = %q, want %q", got, tc.want)
			}
		})
	}
}
