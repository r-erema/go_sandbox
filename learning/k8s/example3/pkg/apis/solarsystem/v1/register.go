package v1

import (
	"github.com/r-erema/go_sendbox/learning/k8s/example3/pkg/apis/solarsystem"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: solarsystem.GroupName, Version: "v1"} //nolint:gochecknoglobals

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&Star{},
		&StarList{},
	)
	v1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes) //nolint:gochecknoglobals
	AddToScheme   = SchemeBuilder.AddToScheme               //nolint:gochecknoglobals
)
