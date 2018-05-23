package controller

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	consularisv1alpha1 "github.com/leosunmo/consularis/pkg/apis/consularis.io/v1alpha1"
	clientset "github.com/leosunmo/consularis/pkg/client/clientset/versioned"
	informers "github.com/leosunmo/consularis/pkg/client/informers/externalversions"
	listers "github.com/leosunmo/consularis/pkg/client/listers/consularis.io/v1alpha1"
	"github.com/leosunmo/consularis/pkg/config"
	consul "github.com/leosunmo/consularis/pkg/consul"
	"github.com/leosunmo/consularis/pkg/signals"
	"github.com/leosunmo/consularis/pkg/utils"
	log "github.com/sirupsen/logrus"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "consulobjects-controller"
const (
	CRDPlural    string = "consulobjects"
	CRDSingular  string = "consulobject"
	CRDShortname string = "co"
	CRDGroup     string = "consularis.io"
	CRDVersion   string = "v1alpha1"
	FullCRDName  string = CRDPlural + "." + CRDGroup
)

// Controller is the controller implementation for Employee resources
type Controller struct {
	// sampleclientset is a clientset for our own API group
	fixturesclientset clientset.Interface

	fixturesLister listers.ConsulObjectLister

	fixturesSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface

	conf *config.Config
}

// Run runs the event loop processing
func Run(conf *config.Config) {

	// Set up CRD in-case it doesn't exist already
	initCRDorDie(conf)
	//namespace := conf.Namespace
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	fixtureClient := utils.GetCRDClient(conf)

	fixturesInformerFactory := informers.NewSharedInformerFactory(fixtureClient, time.Second*30)

	controller := NewController(fixtureClient, fixturesInformerFactory, conf)

	go fixturesInformerFactory.Start(stopCh)

	if err := controller.Run(2, stopCh); err != nil {
		log.WithError(err).Fatal("Error running controller")
	}
}

func initCRDorDie(conf *config.Config) {
	apiExtensionClient := utils.GetAPIExtensionClient(conf)
	_, err := CreateCRD(apiExtensionClient)
	if err != nil {
		log.WithError(err).Fatal("Error creating ConsulObjects CRD")
	}
}

// CreateCRD creates the CRD resource, ignore error if it already exists
func CreateCRD(clientset apiextensionsclientset.Interface) (*apiextv1beta1.CustomResourceDefinition, error) {
	crd := &apiextv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: FullCRDName},
		Spec: apiextv1beta1.CustomResourceDefinitionSpec{
			Group:   CRDGroup,
			Version: CRDVersion,
			Scope:   apiextv1beta1.NamespaceScoped,
			Names: apiextv1beta1.CustomResourceDefinitionNames{
				Singular:   CRDSingular,
				Plural:     CRDPlural,
				ShortNames: strings.Split(CRDShortname, ","),
				Kind:       reflect.TypeOf(consularisv1alpha1.ConsulObject{}).Name(),
			},
		},
	}

	createdCRD, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return createdCRD, nil
	}
	log.Info("Created Custom Resource Definition ", FullCRDName)
	time.Sleep(2 * time.Second)
	return nil, err

	// Note the original apiextensions example adds logic to wait for creation and exception handling
}

// NewController returns a new sample controller
func NewController(
	fixturesclientset clientset.Interface,
	fixturesInformerFactory informers.SharedInformerFactory,
	conf *config.Config) *Controller {

	// obtain references to shared index informers for the ConsulObjects type.
	fixturesInformer := fixturesInformerFactory.Consularis().V1alpha1().ConsulObjects()

	controller := &Controller{
		fixturesclientset: fixturesclientset,
		fixturesLister:    fixturesInformer.Lister(),
		fixturesSynced:    fixturesInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "consulobjects"),
		conf:              conf,
	}

	log.Info("Setting up event handlers")
	// Set up an event handler for when ConsulObjects resources change
	fixturesInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueFixture,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueFixture(new)
		},
		DeleteFunc: controller.enqueueFixture,
	})
	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	log.Info("Starting Archduke controller")

	// Wait for the caches to be synced before starting workers
	log.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.fixturesSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	log.Info("Starting workers")
	// Launch two workers to process Employee resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	log.Info("Started workers")
	<-stopCh
	log.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var ok bool
		var object *consularisv1alpha1.ConsulObject
		// We expect consulObjects to come off the workqueue.
		// We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.

		if object, ok = obj.(*consularisv1alpha1.ConsulObject); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected ConsulObject in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Fixture resource to be synced.
		if err := c.syncHandler(obj); err != nil {
			return fmt.Errorf("error syncing '%s': %s", object.GetName(), err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		log.WithFields(log.Fields{
			"namespace": object.GetNamespace(),
			"name":      object.GetName(),
		}).Info("Successfully synced consulobject")
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Employee resource
// with the current status of the resource.
func (c *Controller) syncHandler(obj interface{}) error {
	var object *consularisv1alpha1.ConsulObject
	var ok bool
	if object, ok = obj.(*consularisv1alpha1.ConsulObject); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return nil
		}
		object, ok = tombstone.Obj.(*consularisv1alpha1.ConsulObject)
		if !ok {
			runtime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return nil
		}
		log.WithFields(log.Fields{
			"namespace": object.GetNamespace(),
			"name":      object.GetName(),
		}).Info("Recovered deleted object from tombstone")
	}

	log.WithFields(log.Fields{
		"namespace": object.GetNamespace(),
		"name":      object.GetName(),
	}).Info("Processing object")
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return err
	}
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	// Create the consul client after we're sure we've got the object
	consulClient, err := consul.NewClient(c.conf)
	// Get the ConsulObjects resource with this namespace/name from API server
	fixture, err := c.fixturesLister.ConsulObjects(namespace).Get(name)
	if err != nil {
		// The fixture resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			err := consul.KVDelete(consulClient, object)
			if err != nil {
				runtime.HandleError(fmt.Errorf("Failed to delete KV: %s", err))
			}

			return nil
		}
		return err
	}

	// Update consul KV with new values. Either new KV or updating existing KV
	err = consul.KVUpdate(consulClient, fixture)
	if err != nil {
		runtime.HandleError(fmt.Errorf("Failed to update KV: %s", err))
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the Consul Fixture resource to reflect
	// the processed/synced state.
	//err = c.updateFixtureStatus(fixture)
	//if err != nil {
	//	return err
	//}

	return nil
}

// enqueueFixture takes a ConsulObject resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than ConsulObject.
func (c *Controller) enqueueFixture(obj interface{}) {
	c.workqueue.AddRateLimited(obj)
}
