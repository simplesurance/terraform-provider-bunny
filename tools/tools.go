//go:build tools
// +build tools

package tools

import (
	// tfplugindocs is required as dependency for go:generate tfplugindocs
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
