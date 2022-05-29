package util

import (
	"time"

	sampleV1Alpha1 "github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/apis/samplecontroller/v1alpha1"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewPodSpec(containers []coreV1.Container) coreV1.PodSpec {
	return coreV1.PodSpec{ // nolint:exhaustruct
		Containers:                    containers,
		Volumes:                       nil,
		InitContainers:                nil,
		EphemeralContainers:           nil,
		RestartPolicy:                 "",
		TerminationGracePeriodSeconds: nil,
		ActiveDeadlineSeconds:         nil,
		DNSPolicy:                     "",
		NodeSelector:                  nil,
		ServiceAccountName:            "",
		AutomountServiceAccountToken:  nil,
		NodeName:                      "",
		HostNetwork:                   false,
		HostPID:                       false,
		HostIPC:                       false,
		ShareProcessNamespace:         nil,
		SecurityContext:               nil,
		ImagePullSecrets:              nil,
		Hostname:                      "",
		Subdomain:                     "",
		Affinity:                      nil,
		SchedulerName:                 "",
		Tolerations:                   nil,
		HostAliases:                   nil,
		PriorityClassName:             "",
		Priority:                      nil,
		DNSConfig:                     nil,
		ReadinessGates:                nil,
		RuntimeClassName:              nil,
		EnableServiceLinks:            nil,
		PreemptionPolicy:              nil,
		Overhead:                      nil,
		TopologySpreadConstraints:     nil,
		SetHostnameAsFQDN:             nil,
		OS:                            nil,
	}
}

func NewDeployment(foo *sampleV1Alpha1.Foo) *appsV1.Deployment {
	labels := map[string]string{
		"app":        "nginx",
		"controller": foo.Name,
	}

	return &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{ // nolint:exhaustruct
			Name:      foo.Spec.DeploymentName,
			Namespace: foo.Namespace,
			OwnerReferences: []metaV1.OwnerReference{
				*metaV1.NewControllerRef(foo, sampleV1Alpha1.SchemeGroupVersion.WithKind("Foo")),
			},
			GenerateName:               "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metaV1.Time{Time: time.Time{}},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			Finalizers:                 nil,
			ManagedFields:              nil,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: foo.Spec.Replicas,
			Selector: &metaV1.LabelSelector{
				MatchLabels:      labels,
				MatchExpressions: nil,
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{ // nolint:exhaustruct
					Labels:                     labels,
					Name:                       "",
					GenerateName:               "",
					Namespace:                  "",
					UID:                        "",
					ResourceVersion:            "",
					Generation:                 0,
					CreationTimestamp:          metaV1.Time{Time: time.Time{}},
					DeletionTimestamp:          nil,
					DeletionGracePeriodSeconds: nil,
					Annotations:                nil,
					OwnerReferences:            nil,
					Finalizers:                 nil,
					ManagedFields:              nil,
				},
				Spec: NewPodSpec([]coreV1.Container{
					{Name: "nginx", Image: "nginx:latest"},
				}),
			},
			Strategy: appsV1.DeploymentStrategy{
				Type:          "",
				RollingUpdate: nil,
			},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    nil,
			Paused:                  false,
			ProgressDeadlineSeconds: nil,
		},
		TypeMeta: metaV1.TypeMeta{Kind: "", APIVersion: ""},
		Status: appsV1.DeploymentStatus{
			ObservedGeneration:  0,
			Replicas:            0,
			UpdatedReplicas:     0,
			ReadyReplicas:       0,
			AvailableReplicas:   0,
			UnavailableReplicas: 0,
			Conditions:          nil,
			CollisionCount:      nil,
		},
	}
}

func NewDeletionOptions() metaV1.DeleteOptions {
	return metaV1.DeleteOptions{ // nolint:exhaustruct
		TypeMeta:           metaV1.TypeMeta{Kind: "", APIVersion: ""},
		GracePeriodSeconds: nil,
		Preconditions:      nil,
		PropagationPolicy:  nil,
		DryRun:             nil,
	}
}

func NewCreateOptions() metaV1.CreateOptions {
	return metaV1.CreateOptions{
		TypeMeta:        metaV1.TypeMeta{Kind: "", APIVersion: ""},
		DryRun:          nil,
		FieldManager:    "",
		FieldValidation: "",
	}
}

func NewUpdateOptions() metaV1.UpdateOptions {
	return metaV1.UpdateOptions{
		TypeMeta:        metaV1.TypeMeta{Kind: "", APIVersion: ""},
		DryRun:          nil,
		FieldManager:    "",
		FieldValidation: "",
	}
}

func NewFoo(replicas int32, namespace, fooName, deploymentName string) *sampleV1Alpha1.Foo {
	return &sampleV1Alpha1.Foo{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		ObjectMeta: metaV1.ObjectMeta{ // nolint:exhaustruct
			Name:                       fooName,
			Namespace:                  namespace,
			Labels:                     nil,
			GenerateName:               "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metaV1.Time{Time: time.Time{}},
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ManagedFields:              nil,
		},
		Spec: sampleV1Alpha1.FooSpec{
			DeploymentName: deploymentName,
			Replicas:       &replicas,
		},
		Status: sampleV1Alpha1.FooStatus{
			AvailableReplicas: replicas,
		},
	}
}
