package k8s

import "time"

type PodInfo struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Status    string            `yaml:"status"`
	Phase     string            `yaml:"phase"`
	Node      string            `yaml:"node"`
	Labels    map[string]string `yaml:"labels"`
	CreatedAt time.Time         `yaml:"createdAt"`
	Restarts  int               `yaml:"restarts"`
}

type ServiceInfo struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace"`
	Type      string            `yaml:"type"`
	ClusterIP string            `yaml:"clusterIP"`
	Ports     []ServicePort     `yaml:"ports"`
	Labels    map[string]string `yaml:"labels"`
	CreatedAt time.Time         `yaml:"createdAt"`
}

type ServicePort struct {
	Name       string `yaml:"name"`
	Port       int32  `yaml:"port"`
	TargetPort string `yaml:"port"`
	Protocol   string `yaml:"protocol"`
}

type DeploymentInfo struct {
	Name            string            `yaml:"name"`
	Namespace       string            `yaml:"namespace"`
	TotalReplicas   int32             `yaml:"totalReplicas"`
	ReadyReplicas   int32             `yaml:"readyReplicas"`
	UpdatedReplicas int32             `yaml:"updatedReplicas"`
	Labels          map[string]string `yaml:"labels"`
	CreatedAt       time.Time         `yaml:"createdAt"`
	Strategy        string            `yaml:"strategy"`
}

type NamespaceInfo struct {
	Name      string            `yaml:"name"`
	Status    string            `yaml:"status"`
	Labels    map[string]string `yaml:"labels"`
	CreatedAt time.Time         `yaml:"createdAt"`
}
