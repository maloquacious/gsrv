package gsrv

import (
	"context"
	"net"
)

type Options []Option
type Option func(*Server) error

func WithContext(ctx context.Context) Option {
	return func(s *Server) error {
		s.admin.ctx = ctx
		return nil
	}
}

func WithHost(host string) Option {
	return func(s *Server) error {
		s.host = host
		s.Addr = net.JoinHostPort(s.host, s.port)
		return nil
	}
}

func WithPort(port string) Option {
	return func(s *Server) error {
		s.port = port
		s.Addr = net.JoinHostPort(s.host, s.port)
		return nil
	}
}

func WithShutdownKey(key string) Option {
	return func(s *Server) error {
		s.admin.keys.shutdown = key
		return nil
	}
}
