package pkg

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listersCoreV1 "k8s.io/client-go/listers/core/v1"
	ingressListerV1 "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type Payload struct {
	Ingresses       []IngressPayload
	TLSCertificates map[string]*tls.Certificate
}

type IngressPayload struct {
	Ingress      *v1.Ingress
	ServicePorts map[string]map[string]int
}

type Watcher struct {
	client   kubernetes.Interface
	onChange func(*Payload)
}

func NewWatcher(client kubernetes.Interface, onChange func(*Payload)) *Watcher {
	return &Watcher{client: client, onChange: onChange}
}

func (w Watcher) Run(ctx context.Context) error {
	factory := informers.NewSharedInformerFactoryWithOptions(w.client, time.Minute)
	secretsLister := factory.Core().V1().Secrets().Lister()
	servicesLister := factory.Core().V1().Services().Lister()
	ingressesLister := factory.Networking().V1().Ingresses().Lister()

	handler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(_ interface{}) {
			onChange(w.onChange, ingressesLister, servicesLister, secretsLister)
		},
		UpdateFunc: func(_, _ interface{}) {
			onChange(w.onChange, ingressesLister, servicesLister, secretsLister)
		},
		DeleteFunc: func(_ interface{}) {
			onChange(w.onChange, ingressesLister, servicesLister, secretsLister)
		},
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go func() {
		informer := factory.Core().V1().Secrets().Informer()
		informer.AddEventHandler(handler)
		informer.Run(ctx.Done())
		waitGroup.Done()
	}()

	waitGroup.Add(1)

	go func() {
		informer := factory.Networking().V1().Ingresses().Informer()
		informer.AddEventHandler(handler)
		informer.Run(ctx.Done())
		waitGroup.Done()
	}()

	waitGroup.Add(1)

	go func() {
		informer := factory.Core().V1().Services().Informer()
		informer.AddEventHandler(handler)
		informer.Run(ctx.Done())
		waitGroup.Done()
	}()

	waitGroup.Wait()

	return nil
}

func addBackend(servicesLister listersCoreV1.ServiceLister, ingressPayload *IngressPayload, backend v1.IngressBackend) {
	svc, err := servicesLister.Services(ingressPayload.Ingress.Namespace).Get(backend.Service.Name)
	if err != nil {
		klog.Errorf(
			"Adding backend error: %s. Service name: %s. Namespace: %s",
			err,
			backend.Service.Name,
			ingressPayload.Ingress.Namespace,
		)

		return
	}

	ports := make(map[string]int)
	for _, port := range svc.Spec.Ports {
		ports[port.Name] = int(port.Port)
	}

	ingressPayload.ServicePorts[svc.Name] = ports
}

func onChange(
	watcherOnChange func(*Payload),
	ingressesLister ingressListerV1.IngressLister,
	servicesLister listersCoreV1.ServiceLister,
	secretsLister listersCoreV1.SecretLister,
) {
	payload := &Payload{
		TLSCertificates: make(map[string]*tls.Certificate),
	}

	ingresses, err := ingressesLister.List(labels.Everything())
	if err != nil {
		klog.Errorf("Getting ingresses list error: %s", err)

		return
	}

	for _, ingress := range ingresses {
		ingressPayload := IngressPayload{
			Ingress:      ingress,
			ServicePorts: make(map[string]map[string]int),
		}

		payload.Ingresses = append(payload.Ingresses, ingressPayload)

		if ingress.Spec.DefaultBackend != nil {
			addBackend(servicesLister, &ingressPayload, *ingress.Spec.DefaultBackend)
		}

		for _, rule := range ingress.Spec.Rules {
			if rule.HTTP != nil {
				continue
			}

			for _, path := range rule.HTTP.Paths {
				addBackend(servicesLister, &ingressPayload, path.Backend)
			}
		}

		handleIngressSpec(ingress, payload, ingressPayload, secretsLister)
	}

	watcherOnChange(payload)
}

func handleIngressSpec(ingress *v1.Ingress, payload *Payload, ingressPayload IngressPayload, secretsLister listersCoreV1.SecretLister) {
	for _, rec := range ingress.Spec.TLS {
		if rec.SecretName == "" {
			continue
		} else {
			secret, err := secretsLister.Secrets(ingress.Namespace).Get(rec.SecretName)
			if err != nil {
				klog.Errorf(
					"Getting secrets error: %s. Secret name: %s. Namespace: %s",
					err,
					rec.SecretName,
					ingressPayload.Ingress.Namespace,
				)

				continue
			}

			klog.Infof("Secret `%s` has been found, namespace: %s", rec.SecretName, ingressPayload.Ingress.Namespace)

			cert, err := tls.X509KeyPair(secret.Data["tls.crt"], secret.Data["tls.key"])
			if err != nil {
				klog.Errorf(
					"Generating X509 key pair error: %s. Secret name: %s. Namespace: %s",
					err,
					rec.SecretName,
					ingressPayload.Ingress.Namespace,
				)

				continue
			}

			payload.TLSCertificates[rec.SecretName] = &cert
		}
	}
}
