package resourceprovider

import (
	"github.com/spf13/cobra"
)

type ResourceProviderAllowlistOptions struct {
	EnableAllowlist bool
}

func AddResourceProviderCliFlags(cmd *cobra.Command, options *ResourceProviderOptions) {
}
