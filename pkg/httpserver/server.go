package httpserver

import (
	"context"
	"github.com/valyala/fasthttp"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
	defaultAddr            = ":8080"
)

type Server struct {
	server          *fasthttp.Server
	notify          chan error
	shutdownTimeout time.Duration
	addr            string
}

func NewServer(h fasthttp.RequestHandler, opts ...Option) *Server {
	httpServer := &fasthttp.Server{
		Handler:      h,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
		addr:            defaultAddr,
	}

	for _, option := range opts {
		option(s)
	}

	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe(s.addr)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.ShutdownWithContext(ctx)
}
