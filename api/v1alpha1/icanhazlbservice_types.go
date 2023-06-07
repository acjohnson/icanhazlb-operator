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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IcanhazlbServiceStatus defines the observed state of IcanhazlbService
type IcanhazlbServiceStatus struct {
	Conditions []IcanhazlbServiceCondition `json:"conditions,omitempty"`
}

// IcanhazlbServiceCondition defines the status condition for the IcanhazlbService
type IcanhazlbServiceCondition struct {
	Type               string                 `json:"type,omitempty"`
	Status             corev1.ConditionStatus `json:"status,omitempty"`
	LastTransitionTime metav1.Time            `json:"lastTransitionTime,omitempty"`
	Reason             string                 `json:"reason,omitempty"`
	Message            string                 `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// IcanhazlbService is the Schema for the icanhazlbservices API
type IcanhazlbService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IcanhazlbServiceSpec   `json:"spec,omitempty"`
	Status IcanhazlbServiceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IcanhazlbServiceList contains a list of IcanhazlbService
type IcanhazlbServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IcanhazlbService `json:"items"`
}

// IcanhazlbServiceSpec defines the desired state for IcanhazlbService
type IcanhazlbServiceSpec struct {
	EndpointSlices EndpointSliceSpec `json:"endpointSlices,omitempty"`
	Services       ServiceSpec       `json:"services,omitempty"`
	Ingresses      IngressSpec       `json:"ingresses,omitempty"`
}

// EndpointSliceSpec defines the specification for an EndpointSlice
type EndpointSliceSpec struct {

	// Name specifies the name of the EndpointSlice
	Name string `json:"name,omitempty"`
	// ProviderSpecific stores provider specific config
	// +optional

	Labels map[string]string `json:"labels,omitempty"`
	// Ports specifies the list of ports exposed by the service.
	// +optional
	Ports []discoveryv1.EndpointPort `json:"ports,omitempty"`

	// AddressType specifies the type of addresses carried by the service.
	// Default is "IPv4".
	// +optional
	AddressType discoveryv1.AddressType `json:"addressType,omitempty"`

	// Endpoints contains information about individual endpoints that
	// comprise the service.
	// +optional
	Endpoints []discoveryv1.Endpoint `json:"endpoints,omitempty"`

	// Conditions represents the latest available observations of the
	// endpoint slice's current state.
	// +optional
	Conditions []discoveryv1.EndpointConditions `json:"conditions,omitempty"`
}

// ServiceSpec defines the specification for the Service
type ServiceSpec struct {
	// Name specifies the name of the Service
	Name                  string                                   `json:"name,omitempty"`
	Labels                map[string]string                        `json:"labels,omitempty"`
	ExternalTrafficPolicy corev1.ServiceExternalTrafficPolicyType  `json:"externalTrafficPolicy,omitempty"`
	InternalTrafficPolicy *corev1.ServiceInternalTrafficPolicyType `json:"internalTrafficPolicy,omitempty"`
	IPFamilies            []corev1.IPFamily                        `json:"ipFamilies,omitempty"`
	Ports                 []corev1.ServicePort                     `json:"ports,omitempty"`
	Selector              map[string]string                        `json:"selector,omitempty"`
	SessionAffinity       corev1.ServiceAffinity                   `json:"sessionAffinity,omitempty"`
	Type                  corev1.ServiceType                       `json:"type,omitempty"`
}

// IngressSpec defines the specifications for ingresses
type IngressSpec struct {
	// Name specifies the name of the Ingress
	Name             string                     `json:"name,omitempty"`
	Annotations      map[string]string          `json:"annotations,omitempty"`
	IngressClassName *string                    `json:"ingressClassName,omitempty"`
	Rules            []networkingv1.IngressRule `json:"rules,omitempty"`
}

func init() {
	SchemeBuilder.Register(&IcanhazlbService{}, &IcanhazlbServiceList{})
}
