package v1alpha1

import (
	"time"

	"github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/apis/samplecontroller"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: samplecontroller.GroupName, Version: "v1alpha1"} //nolint:gochecknoglobals

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

/*func SetupSchemeBuilder() runtime.SchemeBuilder {
	return runtime.NewSchemeBuilder(func(scheme *runtime.Scheme) error {
		scheme.AddKnownTypes(SchemeGroupVersion, &Foo{}, &FooList{})
		v1.AddToGroupVersion(scheme, SchemeGroupVersion)

		return nil
	})
}*/

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, &Foo{
		TypeMeta: v1.TypeMeta{Kind: "", APIVersion: ""},
		ObjectMeta: v1.ObjectMeta{
			Name:                       "",
			Namespace:                  "",
			Labels:                     nil,
			GenerateName:               "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          v1.Time{Time: time.Time{}},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ManagedFields:              nil,
		},
		Spec: FooSpec{
			DeploymentName: "",
			Replicas:       nil,
		},
		Status: FooStatus{
			AvailableReplicas: 0,
		},
	}, &FooList{
		TypeMeta: v1.TypeMeta{Kind: "", APIVersion: ""},
		ListMeta: v1.ListMeta{
			SelfLink:           "",
			ResourceVersion:    "",
			Continue:           "",
			RemainingItemCount: nil,
		},
		Items: nil,
	})
	v1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes) //nolint:gochecknoglobals
	AddToScheme   = SchemeBuilder.AddToScheme               //nolint:gochecknoglobals
)
