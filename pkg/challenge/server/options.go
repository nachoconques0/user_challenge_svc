package server

import "context"

// Pass the address where we want the server to initialize
func WithAddress(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}

type Options struct {
	// Address where transport will be exposed
	Address string
	Context context.Context
}

type Option func(o *Options)
