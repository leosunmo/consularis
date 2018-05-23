package utils

import (
	clientset "github.com/leosunmo/consularis/pkg/client/clientset/versioned"
	"github.com/leosunmo/consularis/pkg/config"
	log "github.com/sirupsen/logrus"
	apiclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClient returns a k8s clientset to the request from inside of cluster
func GetClient() kubernetes.Interface {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.WithError(err).Fatal("Can not get kubernetes config")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.WithError(err).Fatal("Can not create kubernetes client")
	}

	return clientset
}

// GetAPIExtensionClient return a K8s API extension client for accessing CRD
func GetAPIExtensionClient(conf *config.Config) *apiclientset.Clientset {
	cfg, err := clientcmd.BuildConfigFromFlags(conf.MasterURL, conf.Kubeconfig)
	if err != nil {
		log.WithError(err).Fatal("Error building kubeconfig")
	}
	apiextensionsClient, err := apiclientset.NewForConfig(cfg)
	if err != nil {
		log.WithError(err).Fatal("Error building apiExtension clientset")
	}
	return apiextensionsClient
}

// GetCRDClient return a K8s API ConsulObject client
func GetCRDClient(conf *config.Config) *clientset.Clientset {
	cfg, err := clientcmd.BuildConfigFromFlags(conf.MasterURL, conf.Kubeconfig)
	if err != nil {
		log.WithError(err).Fatal("Error building kubeconfig")
	}

	ConsulObjectClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		log.WithError(err).Fatal("Error building ConsulObject clientset")
	}
	return ConsulObjectClient
}
