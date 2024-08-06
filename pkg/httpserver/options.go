package httpserver

import "net"

type Option func(server *Server)

func Port(port string) Option {
	return func(s *Server) {
		s.addr = net.JoinHostPort("", port)
	}
}
