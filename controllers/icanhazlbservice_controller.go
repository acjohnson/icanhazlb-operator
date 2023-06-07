/*
Copyright 2023.

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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	servicev1alpha1 "github.com/acjohnson/icanhazlb-operator/api/v1alpha1"
)

// IcanhazlbServiceReconciler reconciles a IcanhazlbService object
type IcanhazlbServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=service.icanhazlb.com,resources=icanhazlbservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=service.icanhazlb.com,resources=icanhazlbservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=service.icanhazlb.com,resources=icanhazlbservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the IcanhazlbService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *IcanhazlbServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the IcanhazlbService object
	icanhazlbService := &servicev1alpha1.IcanhazlbService{}
	if err := r.Get(ctx, req.NamespacedName, icanhazlbService); err != nil {
		logger.Error(err, "Failed to retrieve IcanhazlbService")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//////////////////// BEGIN EndpointSlice ////////////////////

	// Construct the EndpointSliceObj based on the IcanhazlbService resource
	endpointSliceObj, err := r.newEndpointSlice(icanhazlbService)
	if err != nil {
		logger.Error(err, "Failed to create EndpointSlice")
		return ctrl.Result{}, err
	}

	// Check if the EndpointSlice already exists
	foundEndpointSlice := &discoveryv1.EndpointSlice{}
	err = r.Get(ctx, types.NamespacedName{Name: endpointSliceObj.Name, Namespace: endpointSliceObj.Namespace}, foundEndpointSlice)
	if err != nil && client.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to retrieve EndpointSlice")
		return ctrl.Result{}, err
	}

	//Set Owner Reference for the EndpointSlice
	if err := ctrl.SetControllerReference(icanhazlbService, endpointSliceObj, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on EndpointSlice")
		return ctrl.Result{}, err
	}

	// Create or update the EndpointSlice
	if err := r.updateEndpointSlice(ctx, foundEndpointSlice, endpointSliceObj); err != nil {
		logger.Error(err, "Failed to create or update EndpointSlice")
		return ctrl.Result{}, err
	}

	//////////////////// END EndpointSlice ////////////////////
	//////////////////// BEGIN ServiceObj ////////////////////

	// Construct the ServiceObj based on the IcanhazlbService resource
	serviceObj, err := r.newService(icanhazlbService)
	if err != nil {
		logger.Error(err, "Failed to create ServiceObj")
		return ctrl.Result{}, err
	}

	// Check if the Service already exists
	foundService := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: serviceObj.Name, Namespace: serviceObj.Namespace}, foundService)
	if err != nil && client.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to retrieve Service")
		return ctrl.Result{}, err
	}

	//Set Owner Reference for the Service
	if err := ctrl.SetControllerReference(icanhazlbService, serviceObj, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on Service")
		return ctrl.Result{}, err
	}

	// Create or update the Service
	if err := r.updateService(ctx, foundService, serviceObj); err != nil {
		logger.Error(err, "Failed to create or update Service")
		return ctrl.Result{}, err
	}

	//////////////////// END ServiceObj ////////////////////
	//////////////////// BEGIN IngressObj ////////////////////

	// Construct the IngressObj based on the IcanhazlbService resource
	ingressObj, err := r.newIngress(icanhazlbService)
	if err != nil {
		logger.Error(err, "Failed to create IngressObj")
		return ctrl.Result{}, err
	}

	// Check if the Ingress already exists
	foundIngress := &networkingv1.Ingress{}
	err = r.Get(ctx, types.NamespacedName{Name: ingressObj.Name, Namespace: ingressObj.Namespace}, foundIngress)
	if err != nil && client.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to retrieve Ingress")
		return ctrl.Result{}, err
	}

	// Set Owner Reference for the Ingress
	if err := ctrl.SetControllerReference(icanhazlbService, ingressObj, r.Scheme); err != nil {
		logger.Error(err, "Failed to set owner reference on Ingress")
		return ctrl.Result{}, err
	}

	// Create or update the Ingress
	if err := r.updateIngress(ctx, foundIngress, ingressObj); err != nil {
		logger.Error(err, "Failed to create or update Ingress")
		return ctrl.Result{}, err
	}
	//////////////////// END IngressObj ////////////////////

	logger.Info("Reconciliation completed")
	return ctrl.Result{}, nil
}

func (r *IcanhazlbServiceReconciler) newEndpointSlice(service *servicev1alpha1.IcanhazlbService) (*discoveryv1.EndpointSlice, error) {
	serviceOut, err := json.Marshal(service)
	if err != nil {
		panic(err)
	}

	fmt.Println("server variable newEndpointSlice:" + string(serviceOut))

	// Construct the EndpointSlice object based on the IcanhazlbService
	endpointSliceObj := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Spec.EndpointSlices.Name,
			Labels:    service.Spec.EndpointSlices.Labels,
			Namespace: service.Namespace,
		},
		Ports:       service.Spec.EndpointSlices.Ports,
		AddressType: service.Spec.EndpointSlices.AddressType,
		Endpoints:   service.Spec.EndpointSlices.Endpoints,
	}

	return endpointSliceObj, nil
}

func (r *IcanhazlbServiceReconciler) updateEndpointSlice(ctx context.Context, found *discoveryv1.EndpointSlice, desired *discoveryv1.EndpointSlice) error {
	if found.Endpoints == nil || len(found.Endpoints) == 0 {
		// Create the EndpointSlice if it doesn't exist
		err := r.Create(ctx, desired)
		if err != nil {
			return err
		}
	} else {
		// Update the existing EndpointSlice with the desired state
		found.Ports = desired.Ports
		found.AddressType = desired.AddressType
		found.Endpoints = desired.Endpoints

		err := r.Update(ctx, found)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *IcanhazlbServiceReconciler) newService(service *servicev1alpha1.IcanhazlbService) (*corev1.Service, error) {
	serviceObj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: service.Namespace,
			Name:      service.Spec.Services.Name,
			Labels:    service.Spec.Services.Labels,
		},
		Spec: corev1.ServiceSpec{
			ExternalTrafficPolicy: service.Spec.Services.ExternalTrafficPolicy,
			InternalTrafficPolicy: service.Spec.Services.InternalTrafficPolicy,
			IPFamilies:            service.Spec.Services.IPFamilies,
			Ports:                 []corev1.ServicePort{},
			Selector:              service.Spec.Services.Selector,
			SessionAffinity:       service.Spec.Services.SessionAffinity,
			Type:                  service.Spec.Services.Type,
		},
	}

	// Add ports
	for _, portSpec := range service.Spec.Services.Ports {
		port := corev1.ServicePort{
			Name:       portSpec.Name,
			Port:       portSpec.Port,
			TargetPort: portSpec.TargetPort,
			Protocol:   portSpec.Protocol,
		}
		serviceObj.Spec.Ports = append(serviceObj.Spec.Ports, port)
	}

	return serviceObj, nil
}

func (r *IcanhazlbServiceReconciler) updateService(ctx context.Context, found *corev1.Service, desired *corev1.Service) error {
	if found.ObjectMeta.CreationTimestamp.IsZero() {
		// Service doesn't exist, create a new one
		err := r.Create(ctx, desired)
		if err != nil {
			return fmt.Errorf("failed to create Service: %w", err)
		}
	} else {
		// Service already exists, perform an update
		found.Spec = desired.Spec
		err := r.Update(ctx, found)
		if err != nil {
			return fmt.Errorf("failed to update Service: %w", err)
		}
	}
	return nil
}

func (r *IcanhazlbServiceReconciler) newIngress(service *servicev1alpha1.IcanhazlbService) (*networkingv1.Ingress, error) {
	ingressObj := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   service.Namespace,
			Name:        service.Spec.Ingresses.Name,
			Annotations: service.Spec.Ingresses.Annotations,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: service.Spec.Ingresses.IngressClassName,
			Rules:            []networkingv1.IngressRule{},
		},
	}

	// Add rules
	for _, ruleSpec := range service.Spec.Ingresses.Rules {
		httpPath := networkingv1.HTTPIngressPath{
			Path:     ruleSpec.HTTP.Paths[0].Path,
			PathType: ruleSpec.HTTP.Paths[0].PathType,
			Backend: networkingv1.IngressBackend{
				Service: &networkingv1.IngressServiceBackend{
					Name: ruleSpec.HTTP.Paths[0].Backend.Service.Name,
					Port: networkingv1.ServiceBackendPort{
						Number: ruleSpec.HTTP.Paths[0].Backend.Service.Port.Number,
					},
				},
			},
		}

		httpRule := networkingv1.IngressRule{
			Host:             ruleSpec.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{HTTP: &networkingv1.HTTPIngressRuleValue{Paths: []networkingv1.HTTPIngressPath{httpPath}}},
		}

		ingressObj.Spec.Rules = append(ingressObj.Spec.Rules, httpRule)
	}

	return ingressObj, nil
}

func (r *IcanhazlbServiceReconciler) updateIngress(ctx context.Context, found *networkingv1.Ingress, desired *networkingv1.Ingress) error {
	if found.ObjectMeta.CreationTimestamp.IsZero() {
		// Ingress doesn't exist, create a new one
		err := r.Create(ctx, desired)
		if err != nil {
			return fmt.Errorf("failed to create Ingress: %w", err)
		}
	} else {
		// Ingress already exists, perform an update
		found.Spec = desired.Spec
		err := r.Update(ctx, found)
		if err != nil {
			return fmt.Errorf("failed to update Ingress: %w", err)
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IcanhazlbServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&servicev1alpha1.IcanhazlbService{}).
		Complete(r)
}
