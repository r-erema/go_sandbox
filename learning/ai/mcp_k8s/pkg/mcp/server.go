package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/r-erema/go_sendbox/learning/ai/mcp_k8s/internal/config"
	"github.com/r-erema/go_sendbox/learning/ai/mcp_k8s/internal/logging"
	"k8s.io/client-go/kubernetes"
)

type Server struct {
	config    *config.Config
	clientset *kubernetes.Interface
	logger    *logging.Logger
	mcpServer *server.MCPServer
	formatter *ResourceFormatter
}
