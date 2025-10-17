package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	clientset kubernetes.Interface
	logger    *logrus.Logger
}

func NewClient(configPath string, logger *logrus.Logger) (*Client, error) {
	cfg, err := buildConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("building kubeconfig err: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating k8s client err: %w", err)
	}

	return &Client{
		clientset: clientset,
		logger:    logger,
	}, nil
}

func buildConfig(cfgPath string) (*rest.Config, error) {
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}

	if cfgPath == "" {
		if home := homedir.HomeDir(); home != "" {
			cfgPath = filepath.Join(home, ".kube", "config")
		}
	}

	return clientcmd.BuildConfigFromFlags("", cfgPath)
}

func (c *Client) HealthCheck(ctx context.Context) error {
	if _, err := c.clientset.Discovery().ServerVersion(); err != nil {
		return fmt.Errorf("k8s cluster is not reachable: %w", err)
	}

	return nil
}

func (c *Client) ListPods(ctx context.Context, namespace string) ([]PodInfo, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list pods err: %w", err)
	}

	podInfos := make([]PodInfo, len(pods.Items))
	for i := range pods.Items {
		podInfos[i] = PodInfo{
			Name:      pods.Items[i].GetName(),
			Namespace: pods.Items[i].GetNamespace(),
			Status:    string(pods.Items[i].Status.Phase),
			Phase:     string(pods.Items[i].Status.Phase),
			Node:      pods.Items[i].Spec.NodeName,
			Labels:    pods.Items[i].GetLabels(),
			CreatedAt: pods.Items[i].GetCreationTimestamp().Time,
			Restarts:  getTotalRestarts(&pods.Items[i]),
		}
	}

	return podInfos, nil
}

func (c *Client) ListServices(ctx context.Context, namespace string) ([]ServiceInfo, error) {
	services, err := c.clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list services err: %w", err)
	}

	svcInfos := make([]ServiceInfo, len(services.Items))
	for i := range services.Items {
		ports := make([]ServicePort, len(services.Items[i].Spec.Ports))
		for j := range services.Items[i].Spec.Ports {
			ports[j] = ServicePort{
				Name:       services.Items[i].Spec.Ports[j].Name,
				Port:       services.Items[i].Spec.Ports[j].Port,
				TargetPort: services.Items[i].Spec.Ports[j].TargetPort.String(),
				Protocol:   string(services.Items[i].Spec.Ports[j].Protocol),
			}
		}

		svcInfos[i] = ServiceInfo{
			Name:      services.Items[i].GetName(),
			Namespace: services.Items[i].GetNamespace(),
			Type:      string(services.Items[i].Spec.Type),
			ClusterIP: services.Items[i].Spec.ClusterIP,
			Ports:     ports,
			Labels:    services.Items[i].GetLabels(),
			CreatedAt: services.Items[i].GetCreationTimestamp().Time,
		}
	}

	return svcInfos, nil
}

func (c *Client) ListDeployments(ctx context.Context, namespace string) ([]DeploymentInfo, error) {
	deployments, err := c.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list deployments err: %w", err)
	}

	deploymentInfos := make([]DeploymentInfo, len(deployments.Items))

	for i := range deployments.Items {
		deploymentInfos[i] = DeploymentInfo{
			Name:            deployments.Items[i].GetName(),
			Namespace:       deployments.Items[i].GetNamespace(),
			TotalReplicas:   *deployments.Items[i].Spec.Replicas,
			ReadyReplicas:   deployments.Items[i].Status.ReadyReplicas,
			UpdatedReplicas: deployments.Items[i].Status.UpdatedReplicas,
			Labels:          deployments.Items[i].GetLabels(),
			CreatedAt:       deployments.Items[i].GetCreationTimestamp().Time,
			Strategy:        string(deployments.Items[i].Spec.Strategy.Type),
		}
	}

	return deploymentInfos, nil
}

func getTotalRestarts(_ *v1.Pod) int {
	return 0
}
