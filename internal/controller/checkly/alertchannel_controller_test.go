/*
Copyright 2022.

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

// Example code used for influence: https://github.com/Azure/azure-databricks-operator/blob/0f722a710fea06b86ecdccd9455336ca712bf775/controllers/dcluster_controller_test.go

package checkly

import (
	"context"
	"time"

	"github.com/checkly/checkly-go-sdk"
	checklyv1alpha1 "github.com/checkly/checkly-operator/api/checkly/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("ApiCheck Controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
		acKey := types.NamespacedName{
			Name: "test-alert-channel",
		}
		f := &checklyv1alpha1.AlertChannel{}
		k8sClient.Get(context.Background(), acKey, f)
		k8sClient.Delete(context.Background(), f)
	})

	// Add Tests for OpenAPI validation (or additonal CRD features) specified in
	// your API definition.
	// Avoid adding tests for vanilla CRUD operations because they would
	// test Kubernetes API server, which isn't the goal here.
	Context("AlertChannels", func() {
		It("Full reconciliation", func() {

			acKey := types.NamespacedName{
				Name: "test-alert-channel",
			}

			secretKey := types.NamespacedName{
				Name:      "test-secret",
				Namespace: "default",
			}

			secretData := map[string][]byte{
				"TEST": []byte("test"),
			}

			alertChannel := &checklyv1alpha1.AlertChannel{
				ObjectMeta: metav1.ObjectMeta{
					Name: acKey.Name,
				},
				Spec: checklyv1alpha1.AlertChannelSpec{
					SendFailure: false,
					Email: checkly.AlertChannelEmail{
						Address: "foo@bar.baz",
					},
				},
			}

			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      secretKey.Name,
					Namespace: secretKey.Namespace,
				},
				Data: secretData,
			}

			// Create
			Expect(k8sClient.Create(context.Background(), alertChannel)).Should(Succeed())
			Expect(k8sClient.Create(context.Background(), secret)).Should(Succeed())

			// Status.ID should be present
			By("Expecting AlertChannel ID")
			Eventually(func() bool {
				f := &checklyv1alpha1.AlertChannel{}
				err := k8sClient.Get(context.Background(), acKey, f)
				if f.Status.ID == 3 && err == nil {
					return true
				} else {
					return false
				}
			}, timeout, interval).Should(BeTrue())

			// Finalizer should be present
			By("Expecting finalizer")
			Eventually(func() bool {
				f := &checklyv1alpha1.AlertChannel{}
				err := k8sClient.Get(context.Background(), acKey, f)
				if len(f.Finalizers) == 1 && err == nil {
					return true
				} else {
					return false
				}
			}, timeout, interval).Should(BeTrue())

			// Update
			By("Expecting field update")
			Eventually(func() bool {
				f := &checklyv1alpha1.AlertChannel{}
				err := k8sClient.Get(context.Background(), acKey, f)
				if err != nil {
					return false
				}

				f.Spec.Email = checkly.AlertChannelEmail{}
				f.Spec.SendFailure = true
				f.Spec.OpsGenie = checklyv1alpha1.AlertChannelOpsGenie{
					APISecret: corev1.ObjectReference{
						Namespace: secretKey.Namespace,
						Name:      secretKey.Name,
						FieldPath: "TEST",
					},
					Priority: "999",
					Region:   "US",
				}
				err = k8sClient.Update(context.Background(), f)
				if err != nil {
					return false
				}

				u := &checklyv1alpha1.AlertChannel{}
				err = k8sClient.Get(context.Background(), acKey, u)
				if err != nil {
					return false
				}

				if u.Spec.SendFailure != true {
					return false
				}

				if u.Spec.Email != (checkly.AlertChannelEmail{}) {
					return false
				}

				if u.Spec.OpsGenie.Priority != "999" {
					return false
				}

				return true
			}, timeout, interval).Should(BeTrue())

			// Delete AlertChannel
			By("Expecting to delete alertchannel successfully")
			Eventually(func() error {
				f := &checklyv1alpha1.AlertChannel{}
				k8sClient.Get(context.Background(), acKey, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())

			By("Expecting delete to finish")
			Eventually(func() error {
				f := &checklyv1alpha1.AlertChannel{}
				return k8sClient.Get(context.Background(), acKey, f)
			}, timeout, interval).ShouldNot(Succeed())

			// Delete secret
			By("Expecting to delete secret successfully")
			Eventually(func() error {
				f := &corev1.Secret{}
				k8sClient.Get(context.Background(), secretKey, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())
		})
		// return
	})
})
