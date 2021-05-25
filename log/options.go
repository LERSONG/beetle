package log

import "go.uber.org/zap"

type Options struct {
	config *zap.Config
}

type Option = func(options *Options)

func Config(config *zap.Config) Option {
	return func(opts *Options) {
		opts.config = config
	}
}
