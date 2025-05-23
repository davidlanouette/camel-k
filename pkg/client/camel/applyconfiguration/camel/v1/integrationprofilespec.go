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

// IntegrationProfileSpecApplyConfiguration represents a declarative configuration of the IntegrationProfileSpec type for use
// with apply.
type IntegrationProfileSpecApplyConfiguration struct {
	Build   *IntegrationProfileBuildSpecApplyConfiguration   `json:"build,omitempty"`
	Traits  *TraitsApplyConfiguration                        `json:"traits,omitempty"`
	Kamelet *IntegrationProfileKameletSpecApplyConfiguration `json:"kamelet,omitempty"`
}

// IntegrationProfileSpecApplyConfiguration constructs a declarative configuration of the IntegrationProfileSpec type for use with
// apply.
func IntegrationProfileSpec() *IntegrationProfileSpecApplyConfiguration {
	return &IntegrationProfileSpecApplyConfiguration{}
}

// WithBuild sets the Build field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Build field is set to the value of the last call.
func (b *IntegrationProfileSpecApplyConfiguration) WithBuild(value *IntegrationProfileBuildSpecApplyConfiguration) *IntegrationProfileSpecApplyConfiguration {
	b.Build = value
	return b
}

// WithTraits sets the Traits field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Traits field is set to the value of the last call.
func (b *IntegrationProfileSpecApplyConfiguration) WithTraits(value *TraitsApplyConfiguration) *IntegrationProfileSpecApplyConfiguration {
	b.Traits = value
	return b
}

// WithKamelet sets the Kamelet field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Kamelet field is set to the value of the last call.
func (b *IntegrationProfileSpecApplyConfiguration) WithKamelet(value *IntegrationProfileKameletSpecApplyConfiguration) *IntegrationProfileSpecApplyConfiguration {
	b.Kamelet = value
	return b
}
