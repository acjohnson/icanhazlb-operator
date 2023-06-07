# icanhazlb-operator
This is the the backend logic for `icanhazlb` as it uses a unified object representing
a `icanhazlb` service.

## Description
This operator creates EndpointSlice, Service and Ingress objects in Kubernetes.


This includes an http(s) frontend, TCP load balancer and IP backend configuration.
When an `IcanhazlbService` custom resource is created all three resources will
be created inside Kubernetes brining the backend service online immediatly.

## License

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

