package registry

import (
	"context"
	"crypto/tls"
	"time"
)

// Registry Options
type Options struct {
	ServiceName string
	ServiceAddr string
	Addrs       []string
	Metadata    map[string]string
	Timeout     time.Duration
	Secure      bool
	TLSConfig   *tls.Config
	Context     context.Context
}

type RegisterOptions struct {
	RegisterTTL      time.Duration
	RegisterInterval time.Duration
}

type Option func(*Options)

type RegisterOption func(*RegisterOptions)

func ServiceName(name string) Option {
	return func(o *Options) {
		o.ServiceName = name
	}
}

func ServiceAddr(name string) Option {
	return func(o *Options) {
		o.ServiceAddr = name
	}
}

func Addrs(addrs []string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

func Metadata(m map[string]string) Option {
	return func(o *Options) {
		o.Metadata = m
	}
}

func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// Secure communication with the registry
func Secure(b bool) Option {
	return func(o *Options) {
		o.Secure = b
	}
}

// Specify TLS Config
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func RegisterTTL(ttl time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.RegisterTTL = ttl
	}
}

func RegisterInterval(interval time.Duration) RegisterOption {
	return func(o *RegisterOptions) {
		o.RegisterInterval = interval
	}
}
