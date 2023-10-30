package netpool

import "net"

type Options struct {
	Conns    []net.Conn
	Dialer   Dialer
	Router   Router
	MaxConns uint32
}

type Option interface {
	apply(o *Options)
}

func WithConns(conns ...net.Conn) Option {
	return option(func(o *Options) {
		o.Conns = conns
	})
}

func WithDialer(dialer Dialer) Option {
	return option(func(o *Options) {
		o.Dialer = dialer
	})
}

func WithRouter(router Router) Option {
	return option(func(o *Options) {
		o.Router = router
	})
}

func WithMaxConns(maxConns uint32) Option {
	return option(func(o *Options) {
		o.MaxConns = maxConns
	})
}

type option func(o *Options)

func (f option) apply(o *Options) {
	f(o)
}
