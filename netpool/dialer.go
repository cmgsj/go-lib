package netpool

import (
	"net"
	"time"
)

type Dialer interface {
	Dial() (net.Conn, error)
}

func Dial(network, address string) Dialer {
	return dialer(func() (net.Conn, error) {
		return net.Dial(network, address)
	})
}

func DialTimeout(network, address string, timeout time.Duration) Dialer {
	return dialer(func() (net.Conn, error) {
		return net.DialTimeout(network, address, timeout)
	})
}

func DialIP(network string, laddr, raddr *net.IPAddr) Dialer {
	return dialer(func() (net.Conn, error) {
		return net.DialIP(network, laddr, raddr)
	})
}

func DialTCP(network string, laddr, raddr *net.TCPAddr) Dialer {
	return dialer(func() (net.Conn, error) {
		return net.DialTCP(network, laddr, raddr)
	})
}

func DialUDP(network string, laddr, raddr *net.UDPAddr) Dialer {
	return dialer(func() (net.Conn, error) {
		return net.DialUDP(network, laddr, raddr)
	})
}

func DialUnix(network string, laddr, raddr *net.UnixAddr) Dialer {
	return dialer(func() (net.Conn, error) {
		return net.DialUnix(network, laddr, raddr)
	})
}

type dialer func() (net.Conn, error)

func (f dialer) Dial() (net.Conn, error) {
	return f()
}
