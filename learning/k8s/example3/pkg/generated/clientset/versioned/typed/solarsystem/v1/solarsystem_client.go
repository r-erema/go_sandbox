// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"net/http"

	v1 "github.com/r-erema/go_sendbox/learning/k8s/example3/pkg/apis/solarsystem/v1"
	"github.com/r-erema/go_sendbox/learning/k8s/example3/pkg/generated/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type SolarsystemV1Interface interface {
	RESTClient() rest.Interface
	StarsGetter
}

// SolarsystemV1Client is used to interact with features provided by the solarsystem.k8s.io group.
type SolarsystemV1Client struct {
	restClient rest.Interface
}

func (c *SolarsystemV1Client) Stars(namespace string) StarInterface {
	return newStars(c, namespace)
}

// NewForConfig creates a new SolarsystemV1Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*SolarsystemV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new SolarsystemV1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*SolarsystemV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &SolarsystemV1Client{client}, nil
}

// NewForConfigOrDie creates a new SolarsystemV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *SolarsystemV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new SolarsystemV1Client for the given RESTClient.
func New(c rest.Interface) *SolarsystemV1Client {
	return &SolarsystemV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *SolarsystemV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
