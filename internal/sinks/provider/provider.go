package provider

import "github.com/hashicorp/nomad/api"

type Provider interface {
	// Name returns the name of the Provider.
	Name() string
	// Push pushes a batch of event to upstream. The implementation varies across providers.
	Push([]api.Event)
}
