package provider

import (
	"fmt"

	_ "github.com/ContainerSolutions/argus/operator/internal/provider/checkov"
	_ "github.com/ContainerSolutions/argus/operator/internal/provider/fake"
	_ "github.com/ContainerSolutions/argus/operator/internal/provider/file"
	_ "github.com/ContainerSolutions/argus/operator/internal/provider/random"
	"github.com/ContainerSolutions/argus/operator/internal/provider/schema"
)

func GetProvider(providerName string) (schema.Provider, error) {
	provider, exists := schema.Builder[providerName]
	if !exists {
		return nil, fmt.Errorf("provider '%v' doesn't exist", providerName)
	}
	return provider, nil
}
