//go:build integration
// +build integration

// To enable compilation of this file in Goland, go to "Settings -> Go -> Vendoring & Build Tags -> Custom Tags" and add "integration"

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

package common

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"

	camelv1 "github.com/apache/camel-k/v2/pkg/apis/camel/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"

	. "github.com/apache/camel-k/v2/e2e/support"
)

func TestServiceTrait(t *testing.T) {
	t.Parallel()
	WithNewTestNamespace(t, func(ctx context.Context, g *WithT, ns string) {
		t.Run("Integration Binding", func(t *testing.T) {
			// Not supported when running as a Knative Service as
			// the Knative operator creates an external service with the same name of the Integration
			g.Expect(KamelRun(t, ctx, ns, "hello.yaml", "-t", "knative-service.enabled=false").Execute()).To(Succeed())
			g.Eventually(IntegrationConditionStatus(t, ctx, ns, "hello",
				camelv1.IntegrationConditionReady), TestTimeoutLong).Should(Equal(v1.ConditionTrue))
			g.Eventually(IntegrationPodPhase(t, ctx, ns, "hello")).Should(Equal(corev1.PodRunning))

			ExpectExecSucceed(t, g, Kubectl("apply", "-f", "it-binding.yaml", "-n", ns))
			g.Eventually(IntegrationConditionStatus(t, ctx, ns, "timer-to-hello-it",
				camelv1.IntegrationConditionReady), TestTimeoutLong).Should(Equal(v1.ConditionTrue))
			g.Eventually(IntegrationPodPhase(t, ctx, ns, "timer-to-hello-it")).Should(Equal(corev1.PodRunning))

			g.Eventually(IntegrationLogs(t, ctx, ns, "timer-to-hello-it")).Should(
				ContainSubstring("Hello from Camel"))
		})
	})
}
