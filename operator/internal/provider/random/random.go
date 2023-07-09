package random

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	argusiov1alpha1 "github.com/ContainerSolutions/argus/operator/api/v1alpha1"
	provider "github.com/ContainerSolutions/argus/operator/internal/provider/schema"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var clients = map[string]provider.AttestationClient{}

var mu sync.Mutex

type Client struct {
	reroll  time.Duration
	last    time.Time
	current argusiov1alpha1.AttestationResultType
}

func (c *Client) Attest() (argusiov1alpha1.AttestationResult, error) {
	a := time.Since(c.last)
	if a > c.reroll {
		c.last = time.Now()
		c.current = argusiov1alpha1.AttestationResultTypePass
		rdn := int(rand.Uint64() % 100)
		if rdn == 0 {
			c.current = argusiov1alpha1.AttestationResultTypeFail
		}
	}
	next := c.reroll - time.Since(c.last)
	res := argusiov1alpha1.AttestationResult{
		Result: c.current,
		Logs:   "random",
		RunAt:  v1.Now(),
		Reason: fmt.Sprintf("next random result in %v", next),
	}
	return res, nil
}

func (c *Client) Close() error {
	return nil
}

type Provider struct{}

func (p *Provider) New(name string, spec *argusiov1alpha1.AttestationProviderSpec) (provider.AttestationClient, error) {
	if client, exists := clients[name]; exists {
		return client, nil
	}
	r, ok := spec.ProviderConfig["regenerate"]
	if !ok {
		reroll, err := time.ParseDuration("15m")
		if err != nil {
			return nil, err
		}
		return &Client{reroll: reroll}, nil
	}
	reroll, err := time.ParseDuration(r)
	if err != nil {
		return nil, err
	}
	client := &Client{reroll: reroll}
	mu.Lock()
	clients[name] = client
	defer mu.Unlock()
	return client, nil
}

func init() {
	provider.Register(&Provider{}, "random")
}
