package cpln

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
