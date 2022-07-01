package pkg

type config struct {
	host string
	port,
	tlsPort int
}

func defaultConfig() *config {
	return &config{
		host:    "0.0.0.0",
		port:    80,
		tlsPort: 443,
	}
}

type Option func(*config)

func WithTLSPort(port int) Option {
	return func(c *config) {
		c.tlsPort = port
	}
}
