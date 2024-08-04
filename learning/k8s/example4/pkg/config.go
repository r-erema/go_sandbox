package pkg

type config struct {
	host string
	port,
	tlsPort int
}

const (
	port    = 80
	tlsPort = 80
)

func defaultConfig() *config {
	return &config{
		host:    "0.0.0.0",
		port:    port,
		tlsPort: tlsPort,
	}
}

type Option func(*config)

func WithTLSPort(port int) Option {
	return func(c *config) {
		c.tlsPort = port
	}
}
