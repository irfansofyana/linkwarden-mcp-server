package observability

import (
	"github.com/irfansofyana/linkwarden-mcp-server/pkg/log"
)

// Observability...
type Observability struct {
	Logger log.Logger
}

// Option...
type Option func(*Observability)

// New...
func New(opts ...Option) *Observability {
	obs := Observability{}
	for _, opt := range opts {
		opt(&obs)
	}
	return &obs
}

// WithLogging...
func WithLogging(s log.Logger) Option {
	return func(o *Observability) {
		o.Logger = s
	}
}
