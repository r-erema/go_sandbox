// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	internalinterfaces "github.com/r-erema/go_sendbox/learning/k8s/example3/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Stars returns a StarInformer.
	Stars() StarInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Stars returns a StarInformer.
func (v *version) Stars() StarInformer {
	return &starInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}