//go:build tools

// This file is not used other than to ensure tfplugindocs is included in go.mod.

package tools

import (
       // Documentation generation
       _ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
