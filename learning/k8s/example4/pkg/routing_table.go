package pkg

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	v1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
)

const backendProtocolAnnotation = "kubernetes-simple-ingress-controller/backend-protocol"

type RoutingTable struct {
	certificatesByHost map[string]map[string]*tls.Certificate
	backendsByHost     map[string][]routingTableBackend
}

func NewRoutingTable(payload *Payload) *RoutingTable {
	table := &RoutingTable{
		certificatesByHost: make(map[string]map[string]*tls.Certificate),
		backendsByHost:     make(map[string][]routingTableBackend),
	}

	table.init(payload)

	return table
}

func (rt *RoutingTable) init(payload *Payload) {
	if payload == nil {
		return
	}

	for _, ingressPayload := range payload.Ingresses {
		for _, rule := range ingressPayload.Ingress.Spec.Rules {
			certificates, ok := rt.certificatesByHost[rule.Host]
			if !ok {
				certificates = make(map[string]*tls.Certificate)
				rt.certificatesByHost[rule.Host] = certificates
			}

			for _, TLS := range ingressPayload.Ingress.Spec.TLS {
				for _, host := range TLS.Hosts {
					if cert, ok := payload.TLSCertificates[TLS.SecretName]; ok {
						certificates[host] = cert
					}
				}

				rt.addBackend(ingressPayload, rule)
			}
		}
	}
}

func (rt *RoutingTable) addBackend(ingressPayload IngressPayload, rule v1.IngressRule) {
	scheme, ok := ingressPayload.Ingress.Annotations[backendProtocolAnnotation]
	if !ok {
		scheme = "http"
	}

	scheme = strings.ToLower(scheme)

	if rule.HTTP == nil {
		if ingressPayload.Ingress.Spec.DefaultBackend != nil {
			backend := ingressPayload.Ingress.Spec.DefaultBackend

			rtb, err := newRoutingTableBackend(
				scheme,
				"",
				backend.Service.Name,
				int(backend.Service.Port.Number),
			)
			if err != nil {
				klog.Errorf(
					"Creation new routing table error: %s. Secret name: %s. Namespace: %s",
					err,
					backend.Service.Name,
					ingressPayload.Ingress.Namespace,
				)

				return
			}

			rt.backendsByHost[rule.Host] = append(rt.backendsByHost[rule.Host], *rtb)
		}
	} else {
		rt.handlePaths(scheme, ingressPayload.Ingress.Namespace, &rule)
	}
}

func (rt *RoutingTable) handlePaths(scheme, namespace string, rule *v1.IngressRule) {
	for _, path := range rule.HTTP.Paths {
		backend := path.Backend

		rtb, err := newRoutingTableBackend(
			scheme,
			path.Path,
			backend.Service.Name,
			int(backend.Service.Port.Number),
		)
		if err != nil {
			klog.Errorf(
				"Creation new routing table for path %s error: %s. Secret name: %s. Namespace: %s",
				path.Path,
				err,
				backend.Service.Name,
				namespace,
			)

			continue
		}

		rt.backendsByHost[rule.Host] = append(rt.backendsByHost[rule.Host], *rtb)
	}
}

func (rt *RoutingTable) matches(sni, certHost string) bool {
	for strings.HasPrefix(certHost, "*.") {
		if idx := strings.IndexByte(sni, '.'); idx >= 0 {
			sni = sni[idx+1:]
		} else {
			return false
		}

		certHost = certHost[2:]
	}

	return sni == certHost
}

func (rt *RoutingTable) GetBackend(host, path string) (*url.URL, error) {
	if idx := strings.IndexByte(host, ':'); idx > 0 {
		host = host[:idx]
	}

	backends := rt.backendsByHost[host]
	for _, backend := range backends {
		if backend.matches(path) {
			return backend.url, nil
		}
	}

	return nil, fmt.Errorf("getting backend error: %w", errBackendNotFound)
}

var (
	errBackendNotFound     = errors.New("backend not found")
	errCertificateNotFound = errors.New("certificate not found")
)

func (rt *RoutingTable) GetCertificate(sni string) (*tls.Certificate, error) {
	if hostCerts, ok := rt.certificatesByHost[sni]; ok {
		for host, certificate := range hostCerts {
			if rt.matches(sni, host) {
				return certificate, nil
			}
		}
	}

	return nil, fmt.Errorf("getting certificate %s error: %w", sni, errCertificateNotFound)
}

type routingTableBackend struct {
	pathRegexp *regexp.Regexp
	url        *url.URL
}

func newRoutingTableBackend(
	scheme, path, serviceName string,
	servicePort int,
) (*routingTableBackend, error) {
	rtb := &routingTableBackend{
		url: &url.URL{
			Scheme: scheme,
			Host:   fmt.Sprintf("%s:%d", serviceName, servicePort),
		},
	}

	var err error

	if path != "" {
		rtb.pathRegexp, err = regexp.Compile(path)
	}

	if err != nil {
		return nil, fmt.Errorf("regexp error: %w", err)
	}

	return rtb, nil
}

func (rtb routingTableBackend) matches(path string) bool {
	if rtb.pathRegexp == nil {
		return true
	}

	return rtb.pathRegexp.MatchString(path)
}
