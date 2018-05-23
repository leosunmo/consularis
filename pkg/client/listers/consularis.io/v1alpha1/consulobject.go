/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/leosunmo/consularis/pkg/apis/consularis.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ConsulObjectLister helps list ConsulObjects.
type ConsulObjectLister interface {
	// List lists all ConsulObjects in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.ConsulObject, err error)
	// ConsulObjects returns an object that can list and get ConsulObjects.
	ConsulObjects(namespace string) ConsulObjectNamespaceLister
	ConsulObjectListerExpansion
}

// consulObjectLister implements the ConsulObjectLister interface.
type consulObjectLister struct {
	indexer cache.Indexer
}

// NewConsulObjectLister returns a new ConsulObjectLister.
func NewConsulObjectLister(indexer cache.Indexer) ConsulObjectLister {
	return &consulObjectLister{indexer: indexer}
}

// List lists all ConsulObjects in the indexer.
func (s *consulObjectLister) List(selector labels.Selector) (ret []*v1alpha1.ConsulObject, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ConsulObject))
	})
	return ret, err
}

// ConsulObjects returns an object that can list and get ConsulObjects.
func (s *consulObjectLister) ConsulObjects(namespace string) ConsulObjectNamespaceLister {
	return consulObjectNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ConsulObjectNamespaceLister helps list and get ConsulObjects.
type ConsulObjectNamespaceLister interface {
	// List lists all ConsulObjects in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.ConsulObject, err error)
	// Get retrieves the ConsulObject from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.ConsulObject, error)
	ConsulObjectNamespaceListerExpansion
}

// consulObjectNamespaceLister implements the ConsulObjectNamespaceLister
// interface.
type consulObjectNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ConsulObjects in the indexer for a given namespace.
func (s consulObjectNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ConsulObject, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ConsulObject))
	})
	return ret, err
}

// Get retrieves the ConsulObject from the indexer for a given namespace and name.
func (s consulObjectNamespaceLister) Get(name string) (*v1alpha1.ConsulObject, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("consulobject"), name)
	}
	return obj.(*v1alpha1.ConsulObject), nil
}
