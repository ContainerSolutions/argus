package schema

import (
	"fmt"
	"sync"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
)

var Builder map[string]Provider
var buildlock sync.RWMutex

type Provider interface {
	New(name string, spec *argusiov1alpha1.AttestationProviderSpec) (AttestationClient, error)
}

type AttestationClient interface {
	Attest() (argusiov1alpha1.AttestationResult, error)
	Close() error
}

func init() {
	Builder = make(map[string]Provider)
}

func Register(s Provider, providerName string) {
	buildlock.Lock()
	defer buildlock.Unlock()
	_, exists := Builder[providerName]
	if exists {
		panic(fmt.Sprintf("provider %q already registered", providerName))
	}

	Builder[providerName] = s
}

func ForceRegister(s Provider, providerName string) {
	buildlock.Lock()
	defer buildlock.Unlock()
	Builder[providerName] = s
}
