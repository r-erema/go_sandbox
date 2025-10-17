package mcp

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/json"
)

type ResourceFormatter struct{}

func NewResourceFormatter() *ResourceFormatter {
	return &ResourceFormatter{}
}

func (f ResourceFormatter) FormatPodForAI(podData []byte) (string, error) {

	var pod map[string]interface{}

	if err := json.Unmarshal(podData, &pod); err != nil {
		return "", fmt.Errorf("failed to unmarshal pod data: %v", err)
	}

	summary := &strings.Builder{}
	summary.WriteString("# Pod Summary\n\n")

	summary.WriteString(fmt.Sprintf("**Name**: %s\n", pod["name"]))
	summary.WriteString(fmt.Sprintf("**Namespace**: %s\n", pod["namespace"]))
	summary.WriteString(fmt.Sprintf("**Status**: %s\n", pod["status"]))
	summary.WriteString(fmt.Sprintf("**Node**: %s\n", pod["node"]))

	if restarts, ok := pod["restarts"].(float64); ok && restarts > 0 {
		summary.WriteString(fmt.Sprintf("**⚠️ Restarts**: %.0f\n", int(restarts)))
	}

	if createdAt, ok := pod["createdAt"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			age := time.Since(t)

			summary.WriteString(fmt.Sprintf("**Age**: %s\n", age))
		}
	}

	summary.WriteString(fmt.Sprintf("\n## Containers\n\n"))

	if containers, ok := pod["containers"].([]interface{}); ok {
		for i := range containers {
			if container, ok := containers[i].(map[string]interface{}); ok {
				name := container["name"].(string)
				image := container["image"].(string)
				ready := container["ready"].(bool)
				state := container["state"].(string)

				status := "🟢 Ready"
				if !ready {
					status = "🔴 Not Ready"
				}

				summary.WriteString(fmt.Sprintf("- **%s**: %s\n", name, status))
				summary.WriteString(fmt.Sprintf("- Image: `%s`\n", image))
				summary.WriteString(fmt.Sprintf("- State: %s\n", state))

				if restarts, ok := container["restarts"].(float64); ok && restarts > 0 {
					summary.WriteString(fmt.Sprintf("- Restarts**: %.0f\n", int(restarts)))
				}
			}
		}
	}

	if conditions, ok := pod["conditions"].([]interface{}); ok && len(conditions) > 0 {
		summary.WriteString(fmt.Sprintf("\n## Conditions\n\n"))
		for i := range conditions {
			summary.WriteString(fmt.Sprintf("- %s\n", conditions[i]))
		}
	}

	if labels, ok := pod["labels"].(map[string]interface{}); ok && len(labels) > 0 {
		summary.WriteString(fmt.Sprintf("\n## Labels\n\n"))
		for k, v := range labels {
			summary.WriteString(fmt.Sprintf("- `%s`: `%s`\n", k, v))
		}
	}

	summary.WriteString(fmt.Sprintf("\n---\n"))
	summary.WriteString(fmt.Sprintf("*Use this information to understand the pod's current state and troubleshoot any issues.*"))

	return summary.String(), nil
}
