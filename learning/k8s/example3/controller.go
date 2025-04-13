package example3

import (
	"context"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type StarReconciler struct {
	mu              sync.Mutex
	reconciledTimes uint32
	client.Client
	Scheme *runtime.Scheme
}

func (sr *StarReconciler) Reconcile(
	_ context.Context,
	_ reconcile.Request,
) (reconcile.Result, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	sr.reconciledTimes++

	return ctrl.Result{
		Requeue:      false,
		RequeueAfter: 0,
	}, nil
}

func (sr *StarReconciler) ReconciledTimes() uint32 {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	return sr.reconciledTimes
}
