package example1

import (
	"context"
	"errors"
	"fmt"
	"time"

	sampleV1Alpha1 "github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/apis/samplecontroller/v1alpha1"
	clientSet "github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/generated/clientset/versioned"
	"github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/generated/clientset/versioned/scheme"
	informers "github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/generated/informers/externalversions/samplecontroller/v1alpha1"
	listers "github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/generated/listers/samplecontroller/v1alpha1"
	"github.com/r-erema/go_sendbox/learning/k8s/example1/pkg/util"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsInformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	typedCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appsListers "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const (
	controllerAgentName   = "sample-controller"
	logLevel              = 4
	ErrResourceExists     = "ErrResourceExists"
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	MessageResourceSynced = "Foo synced successfully"
	SuccessSynced         = "Synced"
)

var (
	errWaitForCachesSync     = errors.New("failed to wait for caches to sync")
	errResourceAlreadyExists = errors.New("resource already exists and is not managed by Foo")
	errWorkQueue             = errors.New("bad object for work queue")
	errBadResourceName       = errors.New("bad resource name")
	errInvalidObjectType     = errors.New("invalid object type")
)

type Controller struct {
	kubeClientset   kubernetes.Interface
	sampleClientset clientSet.Interface

	deploymentLister appsListers.DeploymentLister
	deploymentSynced cache.InformerSynced
	fooLister        listers.FooLister
	fooSynced        cache.InformerSynced

	//nolint
	// todo: get rid of nolint
	workQueue workqueue.RateLimitingInterface //nolint

	recorder record.EventRecorder
}

func NewController(
	kubeClientSet kubernetes.Interface,
	sampleClientSet clientSet.Interface,
	deploymentInformer appsInformers.DeploymentInformer,
	fooInformer informers.FooInformer,
) *Controller {
	utilRuntime.Must(sampleV1Alpha1.AddToScheme(scheme.Scheme))
	klog.V(logLevel).Info("Creating event broadcaster")

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedCoreV1.EventSinkImpl{Interface: kubeClientSet.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, coreV1.EventSource{Component: controllerAgentName, Host: ""})

	controller := &Controller{
		kubeClientset:    kubeClientSet,
		sampleClientset:  sampleClientSet,
		deploymentLister: deploymentInformer.Lister(),
		deploymentSynced: deploymentInformer.Informer().HasSynced,
		fooLister:        fooInformer.Lister(),
		fooSynced:        fooInformer.Informer().HasSynced,

		//nolint
		// todo: get rid of nolint
		workQueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Foos"), //nolint
		recorder:  recorder,
	}

	klog.Info("Setting up event handlers")

	_, _ = fooInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueFoo,
		UpdateFunc: func(_, newObj interface{}) {
			controller.enqueueFoo(newObj)
		},
		DeleteFunc: nil,
	})

	_, _ = deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(oldObj, newObj interface{}) {
			var (
				newDeployment, oldDeployment *appsV1.Deployment
				ok                           bool
			)

			if newDeployment, ok = newObj.(*appsV1.Deployment); !ok {
				klog.Errorf("couldn't assert type for the new deployment object")

				return
			}

			if oldDeployment, ok = oldObj.(*appsV1.Deployment); !ok {
				klog.Errorf("couldn't assert type for the old deployment object")

				return
			}

			if newDeployment.ResourceVersion == oldDeployment.ResourceVersion {
				return
			}

			controller.handleObject(newObj)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilRuntime.HandleCrash()
	defer c.workQueue.ShutDown()

	klog.Info("Starting Foo controller")

	klog.Info("Waiting for informer caches to sync")

	if !cache.WaitForCacheSync(stopCh, c.deploymentSynced, c.fooSynced) {
		return errWaitForCachesSync
	}

	klog.Info("Starting workers")

	for range workers {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workQueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workQueue.Done(obj)

		var (
			key string
			ok  bool
		)

		if key, ok = obj.(string); !ok {
			c.workQueue.Forget(obj)
			utilRuntime.HandleError(fmt.Errorf("%w, expected string in work queue but got %#v", errWorkQueue, obj))

			return nil
		}

		if err := c.syncHandler(key); err != nil {
			c.workQueue.AddRateLimited(key)

			return fmt.Errorf("error syncing '%s': %w, requeuing", key, err)
		}

		c.workQueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)

		return nil
	}(obj)
	if err != nil {
		utilRuntime.HandleError(err)

		return true
	}

	return true
}

func (c *Controller) handleObject(obj interface{}) {
	var (
		object metaV1.Object
		ok     bool
	)

	if object, ok = obj.(metaV1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			c.handleObject(fmt.Errorf("%w, %s", errInvalidObjectType, "error decoding object"))

			return
		}

		object, ok = tombstone.Obj.(metaV1.Object)
		if !ok {
			c.handleObject(fmt.Errorf("%w, %s", errInvalidObjectType, "error decoding object tombstone"))

			return
		}

		klog.V(logLevel).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}

	klog.V(logLevel).Infof("Processing object: %s", object.GetName())

	if ownerRef := metaV1.GetControllerOf(object); ownerRef != nil {
		if ownerRef.Kind != "Foo" {
			return
		}

		foo, err := c.fooLister.Foos(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(logLevel).Infof("ignoring orphaned object '%s/%s' of foo '%s'", object.GetNamespace(), object.GetName(), ownerRef.Name)

			return
		}

		c.enqueueFoo(foo)

		return
	}
}

func (c *Controller) enqueueFoo(obj interface{}) {
	var (
		key string
		err error
	)

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilRuntime.HandleError(err)

		return
	}

	c.workQueue.Add(key)
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilRuntime.HandleError(fmt.Errorf("invalid resource key: %s, error: %w", key, err))

		return nil
	}

	foo, err := c.fooLister.Foos(namespace).Get(name)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			utilRuntime.HandleError(fmt.Errorf("%w, foo '%s' in work queue no longer exists", errWorkQueue, key))

			return nil
		}

		return fmt.Errorf("getting resource error: %w", err)
	}

	deploymentName := foo.Spec.DeploymentName
	if deploymentName == "" {
		utilRuntime.HandleError(fmt.Errorf("%w, %s: deployment name must be specified", errBadResourceName, key))

		return nil
	}

	deployment, err := c.deploymentLister.Deployments(foo.Namespace).Get(deploymentName)
	if k8sErrors.IsNotFound(err) {
		deployment, err = c.kubeClientset.AppsV1().Deployments(foo.Namespace).Create(
			context.TODO(),
			util.NewDeployment(foo),
			util.NewCreateOptions(),
		)
	}

	if err != nil {
		return fmt.Errorf("getting deployment error: %w", err)
	}

	if !metaV1.IsControlledBy(deployment, foo) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(foo, coreV1.EventTypeWarning, ErrResourceExists, msg)

		return fmt.Errorf("%w: %s", errResourceAlreadyExists, msg)
	}

	deployment, err = c.handleDeploymentUpdate(foo, deployment, name)
	if err != nil {
		return fmt.Errorf("deployment update error: %w", err)
	}

	err = c.updateFooStatus(foo, deployment)
	if err != nil {
		return fmt.Errorf("update status error: %w", err)
	}

	c.recorder.Event(foo, coreV1.EventTypeNormal, SuccessSynced, MessageResourceSynced)

	return nil
}

func (c *Controller) handleDeploymentUpdate(
	foo *sampleV1Alpha1.Foo,
	deployment *appsV1.Deployment,
	name string,
) (*appsV1.Deployment, error) {
	var err error

	if foo.Spec.Replicas != nil && *foo.Spec.Replicas != *deployment.Spec.Replicas {
		klog.V(logLevel).Infof("Foo %s replicas: %d, deployment replicas: %d", name, *foo.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeClientset.AppsV1().Deployments(foo.Namespace).Update(
			context.TODO(),
			util.NewDeployment(foo),
			util.NewUpdateOptions(),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("deployment update error: %w", err)
	}

	return deployment, nil
}

func (c *Controller) updateFooStatus(foo *sampleV1Alpha1.Foo, deployment *appsV1.Deployment) error {
	fooCopy := foo.DeepCopy()
	fooCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas

	_, err := c.sampleClientset.SamplecontrollerV1alpha1().Foos(foo.Namespace).UpdateStatus(
		context.TODO(),
		fooCopy,
		util.NewUpdateOptions(),
	)
	if err != nil {
		return fmt.Errorf("updating status error: %w", err)
	}

	return nil
}
