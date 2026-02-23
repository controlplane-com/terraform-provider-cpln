package cpln

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// HelmCommonConfig holds all common configuration for building helm CLI arguments.
// Shared between cpln_helm_release (resource) and cpln_helm_template (data source).
type HelmCommonConfig struct {
	Gvc                   types.String
	Repository            types.String
	Version               types.String
	Values                types.List
	Set                   types.Map
	SetString             types.Map
	SetFile               types.Map
	Wait                  types.Bool
	Timeout               types.Int32
	Description           types.String
	Verify                types.Bool
	RepositoryUsername    types.String
	RepositoryPassword    types.String
	RepositoryCaFile      types.String
	RepositoryCertFile    types.String
	RepositoryKeyFile     types.String
	InsecureSkipTLSVerify types.Bool
	RenderSubchartNotes   types.Bool
	Postrender            types.Object
	DependencyUpdate      types.Bool
	MaxHistory            types.Int32
}

// HelmReleaseState represents the internal state carrier for helm release operations.
type HelmReleaseState struct {
	Name      string
	Status    string
	Revision  int
	Manifest  string
	Resources map[string]string
}

// HelmGetAllResponse represents the JSON output from cpln helm get all <name> -o json.
type HelmGetAllResponse struct {
	Name     string         `json:"name"`
	Version  int            `json:"version"`
	Manifest string         `json:"manifest"`
	Gvc      string         `json:"gvc"`
	Info     HelmGetAllInfo `json:"info"`
}

// HelmGetAllInfo represents the info section of the helm get all response.
type HelmGetAllInfo struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

// HelmManifestResource represents a top-level resource in the helm manifest YAML output.
type HelmManifestResource struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
	Gvc  string `yaml:"gvc"`
}

// BuildHelmArgs adds common arguments for helm commands (install, upgrade, template).
func (c *Client) BuildHelmArgs(args []string, cfg HelmCommonConfig) ([]string, []string, error) {
	var tempFilePaths []string

	// Context options
	args = c.AppendCplnContextArgs(args, cfg.Gvc.ValueString())

	// Chart resolution
	if !cfg.Repository.IsNull() && cfg.Repository.ValueString() != "" {
		args = append(args, "--repo", cfg.Repository.ValueString())
	}

	if !cfg.Version.IsNull() && cfg.Version.ValueString() != "" {
		args = append(args, "--version", cfg.Version.ValueString())
	}

	// Values - each entry written to a temp file and passed via --values
	if !cfg.Values.IsNull() && len(cfg.Values.Elements()) > 0 {
		for i, elem := range cfg.Values.Elements() {
			strVal, ok := elem.(types.String)
			if !ok || strVal.IsNull() || strVal.ValueString() == "" {
				continue
			}

			valuesFile, err := os.CreateTemp("", fmt.Sprintf("helm-values-%d-*.yaml", i))
			if err != nil {
				RemoveTempFiles(tempFilePaths)
				return nil, nil, fmt.Errorf("failed to create temp values file: %w", err)
			}

			tempFilePaths = append(tempFilePaths, valuesFile.Name())

			if _, err := valuesFile.WriteString(strVal.ValueString()); err != nil {
				valuesFile.Close()
				RemoveTempFiles(tempFilePaths)
				return nil, nil, fmt.Errorf("failed to write values to temp file: %w", err)
			}

			valuesFile.Close()
			args = append(args, "--values", valuesFile.Name())
		}
	}

	if !cfg.Set.IsNull() {
		for key, value := range cfg.Set.Elements() {
			if strVal, ok := value.(types.String); ok && !strVal.IsNull() {
				args = append(args, "--set", fmt.Sprintf("%s=%s", key, strVal.ValueString()))
			}
		}
	}

	if !cfg.SetString.IsNull() {
		for key, value := range cfg.SetString.Elements() {
			if strVal, ok := value.(types.String); ok && !strVal.IsNull() {
				args = append(args, "--set-string", fmt.Sprintf("%s=%s", key, strVal.ValueString()))
			}
		}
	}

	if !cfg.SetFile.IsNull() {
		for key, value := range cfg.SetFile.Elements() {
			if strVal, ok := value.(types.String); ok && !strVal.IsNull() {
				args = append(args, "--set-file", fmt.Sprintf("%s=%s", key, strVal.ValueString()))
			}
		}
	}

	// Deployment options
	if !cfg.Wait.IsNull() && cfg.Wait.ValueBool() {
		args = append(args, "--wait")

		if !cfg.Timeout.IsNull() {
			args = append(args, "--timeout", fmt.Sprintf("%d", cfg.Timeout.ValueInt32()))
		}
	}

	if !cfg.Description.IsNull() && cfg.Description.ValueString() != "" {
		args = append(args, "--description", cfg.Description.ValueString())
	}

	if !cfg.Verify.IsNull() && cfg.Verify.ValueBool() {
		args = append(args, "--verify")
	}

	if !cfg.DependencyUpdate.IsNull() && cfg.DependencyUpdate.ValueBool() {
		args = append(args, "--dependency-update")
	}

	if !cfg.MaxHistory.IsNull() {
		args = append(args, "--history-limit", fmt.Sprintf("%d", cfg.MaxHistory.ValueInt32()))
	}

	// Repository auth
	if !cfg.RepositoryUsername.IsNull() && cfg.RepositoryUsername.ValueString() != "" {
		args = append(args, "--username", cfg.RepositoryUsername.ValueString())
	}

	if !cfg.RepositoryPassword.IsNull() && cfg.RepositoryPassword.ValueString() != "" {
		args = append(args, "--password", cfg.RepositoryPassword.ValueString())
	}

	// TLS options
	if !cfg.RepositoryCaFile.IsNull() && cfg.RepositoryCaFile.ValueString() != "" {
		args = append(args, "--ca-file", cfg.RepositoryCaFile.ValueString())
	}

	if !cfg.RepositoryCertFile.IsNull() && cfg.RepositoryCertFile.ValueString() != "" {
		args = append(args, "--cert-file", cfg.RepositoryCertFile.ValueString())
	}

	if !cfg.RepositoryKeyFile.IsNull() && cfg.RepositoryKeyFile.ValueString() != "" {
		args = append(args, "--key-file", cfg.RepositoryKeyFile.ValueString())
	}

	if !cfg.InsecureSkipTLSVerify.IsNull() && cfg.InsecureSkipTLSVerify.ValueBool() {
		args = append(args, "--insecure-skip-tls-verify")
	}

	// Rendering options
	if !cfg.RenderSubchartNotes.IsNull() && cfg.RenderSubchartNotes.ValueBool() {
		args = append(args, "--render-subchart-notes")
	}

	if !cfg.Postrender.IsNull() && !cfg.Postrender.IsUnknown() {
		attrs := cfg.Postrender.Attributes()

		if binaryPath, ok := attrs["binary_path"]; ok {
			if strVal, ok := binaryPath.(types.String); ok && !strVal.IsNull() && strVal.ValueString() != "" {
				args = append(args, "--post-renderer", strVal.ValueString())
			}
		}

		if postArgs, ok := attrs["args"]; ok {
			if listVal, ok := postArgs.(types.List); ok && !listVal.IsNull() {
				for _, arg := range listVal.Elements() {
					if strVal, ok := arg.(types.String); ok && !strVal.IsNull() {
						args = append(args, "--post-renderer-args", strVal.ValueString())
					}
				}
			}
		}
	}

	return args, tempFilePaths, nil
}

// AppendCplnAuthArgs appends --org, --token, and --endpoint flags to the args slice.
func (c *Client) AppendCplnAuthArgs(args []string) []string {
	args = append(args, "--org", c.Org)
	args = append(args, "--token", c.Token)

	if c.HostURL != "" && c.HostURL != DefaultClientEndpoint {
		args = append(args, "--endpoint", c.HostURL)
	}

	return args
}

// AppendCplnContextArgs appends auth flags and optionally --gvc to the args slice.
func (c *Client) AppendCplnContextArgs(args []string, gvc string) []string {
	args = c.AppendCplnAuthArgs(args)

	if gvc != "" {
		args = append(args, "--gvc", gvc)
	}

	return args
}

// RemoveTempFiles removes a list of temporary files by their paths.
func RemoveTempFiles(paths []string) {
	for _, p := range paths {
		os.Remove(p)
	}
}
