package pkg

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
)

type Server struct {
	cfg          *config
	routingTable atomic.Value

	ready *Event
}

func NewServer(options ...Option) *Server {
	cfg := defaultConfig()
	for _, option := range options {
		option(cfg)
	}

	server := &Server{ //nolint:exhaustruct
		cfg:   cfg,
		ready: NewEvent(),
	}

	server.routingTable.Store(NewRoutingTable(nil))

	return server
}

var errTypeCasting = errors.New("type casting error")

func (s *Server) Run(ctx context.Context) error {
	s.ready.Wait(ctx)

	srv := http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf("%s:%d", s.cfg.host, s.cfg.tlsPort),
		Handler:           s,
		ReadHeaderTimeout: time.Second,
	}
	srv.TLSConfig = &tls.Config{ //nolint:exhaustruct
		MinVersion: tls.VersionTLS12,
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if table, ok := s.routingTable.Load().(*RoutingTable); ok {
				return table.GetCertificate(info.ServerName)
			}

			return nil, fmt.Errorf("getting certifaicate error: %w", errTypeCasting)
		},
	}

	err := srv.ListenAndServeTLS("", "")
	if err != nil {
		return fmt.Errorf("starting server error: %w", err)
	}

	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	backendURL, err := s.routingTable.Load().(*RoutingTable).GetBackend(request.Host, request.URL.Path)
	if err != nil {
		http.Error(writer, "getting backend error", http.StatusInternalServerError)

		return
	}

	p := httputil.NewSingleHostReverseProxy(backendURL)
	p.Transport = &http2.Transport{ //nolint:exhaustruct
		AllowHTTP: true,
	}

	p.ServeHTTP(writer, request)
}

func (s *Server) Update(payload *Payload) {
	s.routingTable.Store(NewRoutingTable(payload))
	s.ready.Set()
}
