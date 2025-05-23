/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// CamelSchemeScopeApplyConfiguration represents a declarative configuration of the CamelSchemeScope type for use
// with apply.
type CamelSchemeScopeApplyConfiguration struct {
	Dependencies []CamelArtifactDependencyApplyConfiguration `json:"dependencies,omitempty"`
}

// CamelSchemeScopeApplyConfiguration constructs a declarative configuration of the CamelSchemeScope type for use with
// apply.
func CamelSchemeScope() *CamelSchemeScopeApplyConfiguration {
	return &CamelSchemeScopeApplyConfiguration{}
}

// WithDependencies adds the given value to the Dependencies field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Dependencies field.
func (b *CamelSchemeScopeApplyConfiguration) WithDependencies(values ...*CamelArtifactDependencyApplyConfiguration) *CamelSchemeScopeApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithDependencies")
		}
		b.Dependencies = append(b.Dependencies, *values[i])
	}
	return b
}
