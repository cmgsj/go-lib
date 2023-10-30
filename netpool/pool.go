package netpool

import (
	"errors"
	"net"
	"sync"
)

type Pool interface {
	Get() (net.Conn, error)
	Count() int
	Close() error
}

func NewPool(opts ...Option) (Pool, error) {
	var o Options
	for _, opt := range opts {
		opt.apply(&o)
	}
	n := uint32(len(o.Conns))
	if n == 0 && o.MaxConns == 0 {
		return nil, errors.New("must have conns or max conns")
	}
	if n > o.MaxConns {
		return nil, errors.New("conns are greater that max conns")
	}
	if n < o.MaxConns && o.Dialer == nil {
		return nil, errors.New("must have dialer when conns are less than max conns")
	}
	if n == 0 {
		o.Conns = make([]net.Conn, 0)
	}
	if o.MaxConns == 0 {
		o.MaxConns = n
	}
	if o.Router == nil {
		o.Router = RoundRobin()
	}
	return &pool{
		conns:  o.Conns,
		dialer: o.Dialer,
		router: o.Router,
		max:    o.MaxConns,
	}, nil
}

type pool struct {
	mu     sync.Mutex
	conns  []net.Conn
	dialer Dialer
	router Router
	max    uint32
	idx    uint32
}

func (p *pool) Get() (net.Conn, error) {
	p.mu.Lock()
	n := uint32(len(p.conns))
	if n < p.max && p.dialer != nil {
		conn, err := p.dialer.Dial()
		if err != nil {
			if n == 0 {
				return nil, err
			}
		} else {
			p.conns = append(p.conns, conn)
			n++
		}
	}
	p.idx = p.router.Next(p.idx, n) % n
	conn := p.conns[p.idx]
	p.mu.Unlock()
	return conn, nil
}

func (p *pool) Count() int {
	p.mu.Lock()
	n := len(p.conns)
	p.mu.Unlock()
	return n
}

func (p *pool) Close() error {
	p.mu.Lock()
	var errs []error
	for _, conn := range p.conns {
		errs = append(errs, conn.Close())
	}
	err := errors.Join(errs...)
	p.mu.Unlock()
	return err
}
