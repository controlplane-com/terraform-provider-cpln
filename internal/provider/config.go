package cpln

import "os"

// OrgName retrieves the Control Plane organization name from the CPLN_ORG environment variable for use in tests
var OrgName string = os.Getenv("CPLN_ORG")

// GvcScopedKinds lists all resource kinds that are scoped within a GVC.
var GvcScopedKinds = []string{"identity", "workload", "volumeset"}

// IgnoredTagPrefixes contains prefixes for tags that should be excluded from processing.
var IgnoredTagPrefixes = []string{
	"cpln/deployTimestamp",
	"cpln/aws",
	"cpln/azure",
	"cpln/docker",
	"cpln/gcp",
	"cpln/tls",
	"cpln/managedByTerraform",
	"cpln/city",
	"cpln/continent",
	"cpln/country",
	"cpln/state",
}

// GcpRoles contains the GCP roles used during the configuration of the cloud account at GCP.
var GcpRoles = []string{"roles/viewer", "roles/iam.serviceAccountAdmin", "roles/iam.serviceAccountTokenCreator"}

// CustomLocationCloudProviders is a list of custom location cloud proivder names.
var CustomLocationCloudProviders = []string{"byok"}

// AllowedCustomLocationProviders is a list of custom location cloud provider names that are allowed to be used.
var AllowedCustomLocationProviders = []string{"byok"}

// CloudAccountIdentifiers holds the default AWS identifiers
var CloudAccountIdentifiers = []string{"arn:aws:iam::957753459089:user/controlplane-driver", "arn:aws:iam::957753459089:role/controlplane-driver"}
