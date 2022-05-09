package networking

import (
	"context"
	"fmt"
	"time"

	checklyv1alpha1 "github.com/imgarena/checkly-operator/apis/checkly/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Ingress Controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	BeforeEach(func() {
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("Ingress", func() {

		// Test happy path
		It("full reconciliation", func() {

			testHost := "foo.bar"
			testPath := "baz"
			testGroup := "ingress-group"
			testSuccessCode := "200"

			key := types.NamespacedName{
				Name:      "test-ingress",
				Namespace: "default",
			}

			annotation := make(map[string]string)
			annotation["checkly.imgarena.com/enabled"] = "true"
			annotation["checkly.imgarena.com/path"] = testPath
			annotation["checkly.imgarena.com/success"] = testSuccessCode
			annotation["checkly.imgarena.com/group"] = testGroup

			rules := make([]networkingv1.IngressRule, 0)
			rules = append(rules, networkingv1.IngressRule{
				Host: testHost,
			})

			ingress := &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:        key.Name,
					Namespace:   key.Namespace,
					Annotations: annotation,
				},
				Spec: networkingv1.IngressSpec{
					Rules: rules,
					DefaultBackend: &networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: "test-service",
							Port: networkingv1.ServiceBackendPort{
								Number: 7777,
							},
						},
					},
				},
			}

			// Create
			Expect(k8sClient.Create(context.Background(), ingress)).Should(Succeed())

			By("Expecting submitted")
			Eventually(func() bool {
				f := &networkingv1.Ingress{}
				err := k8sClient.Get(context.Background(), key, f)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("Expecting ApiCheck and OwnerReference to exist")
			Eventually(func() bool {
				f := &checklyv1alpha1.ApiCheck{}
				err := k8sClient.Get(context.Background(), key, f)
				if err != nil {
					return false
				}

				if len(f.OwnerReferences) != 1 {
					return false
				}

				Expect(f.Spec.Endpoint == fmt.Sprintf("https://%s%s", testHost, testPath)).To(BeTrue())
				Expect(f.Spec.Group).To(Equal(testGroup))
				Expect(f.Spec.Success).To(Equal(testSuccessCode))

				for _, o := range f.OwnerReferences {
					if o.Name != key.Name {
						return false
					}
				}

				return true
			}, timeout, interval).Should(BeTrue())

			// Update
			updatePath := "baaz"
			updateHost := "foo.update"
			annotation["checkly.imgarena.com/path"] = updatePath
			annotation["checkly.imgarena.com/endpoint"] = updateHost
			annotation["checkly.imgarena.com/success"] = ""
			ingress = &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:        key.Name,
					Namespace:   key.Namespace,
					Annotations: annotation,
				},
				Spec: networkingv1.IngressSpec{
					Rules: rules,
					DefaultBackend: &networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: "test-service",
							Port: networkingv1.ServiceBackendPort{
								Number: 7777,
							},
						},
					},
				},
			}
			Expect(k8sClient.Update(context.Background(), ingress)).Should(Succeed())

			By("Expecting ApiCheck to be updated")
			Eventually(func() bool {
				f := &checklyv1alpha1.ApiCheck{}
				err := k8sClient.Get(context.Background(), key, f)
				if err != nil {
					return false
				}

				if f.Spec.Endpoint != fmt.Sprintf("https://%s%s", updateHost, updatePath) {
					return false
				}

				if f.Spec.Success != "200" {
					return false
				}

				return true
			}, timeout, interval).Should(BeTrue())

			// Remove enabled label
			annotation["checkly.imgarena.com/enabled"] = "false"
			ingress = &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:        key.Name,
					Namespace:   key.Namespace,
					Annotations: annotation,
				},
				Spec: networkingv1.IngressSpec{
					Rules: rules,
					DefaultBackend: &networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: "test-service",
							Port: networkingv1.ServiceBackendPort{
								Number: 7777,
							},
						},
					},
				},
			}
			Expect(k8sClient.Update(context.Background(), ingress)).Should(Succeed())

			// Expect ApiCheck to be deleted
			By("Expecting APICheck to be deleted")
			Eventually(func() error {
				f := &checklyv1alpha1.ApiCheck{}
				return k8sClient.Get(context.Background(), key, f)
			}, timeout, interval).ShouldNot(Succeed())

			// Delete
			By("Expecting to delete successfully")
			Eventually(func() error {
				f := &networkingv1.Ingress{}
				k8sClient.Get(context.Background(), key, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())

			By("Expecting delete to finish")
			Eventually(func() error {
				f := &networkingv1.Ingress{}
				return k8sClient.Get(context.Background(), key, f)
			}, timeout, interval).ShouldNot(Succeed())

		})

		// Testing failures
		It("Some failures", func() {
			testHost := "foo.bar"
			testPath := "baz"
			// testGroup := "ingress-group"
			testSuccessCode := "200"

			key := types.NamespacedName{
				Name:      "test-fail-ingress",
				Namespace: "default",
			}

			annotation := make(map[string]string)
			annotation["checkly.imgarena.com/enabled"] = "false"
			annotation["checkly.imgarena.com/path"] = testPath
			annotation["checkly.imgarena.com/success"] = testSuccessCode

			rules := make([]networkingv1.IngressRule, 0)
			rules = append(rules, networkingv1.IngressRule{
				Host: testHost,
			})

			ingress := &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:        key.Name,
					Namespace:   key.Namespace,
					Annotations: annotation,
				},
				Spec: networkingv1.IngressSpec{
					Rules: rules,
					DefaultBackend: &networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: "test-service",
							Port: networkingv1.ServiceBackendPort{
								Number: 7777,
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(context.Background(), ingress)).Should(Succeed())

			// Test group annotation missing
			annotation["checkly.imgarena.com/enabled"] = "true"
			ingress.Annotations = annotation
			Expect(k8sClient.Update(context.Background(), ingress)).Should(Succeed())

			// Delete
			By("Expecting to delete successfully")
			Eventually(func() error {
				f := &networkingv1.Ingress{}
				k8sClient.Get(context.Background(), key, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())

			By("Expecting delete to finish")
			Eventually(func() error {
				f := &networkingv1.Ingress{}
				return k8sClient.Get(context.Background(), key, f)
			}, timeout, interval).ShouldNot(Succeed())

		})
	})

})
