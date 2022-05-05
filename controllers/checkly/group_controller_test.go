package checkly

import (
	"context"
	"time"

	checklyv1alpha1 "github.com/imgarena/checkly-operator/apis/checkly/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	// Add Tests for OpenAPI validation (or additonal CRD features) specified in
	// your API definition.
	// Avoid adding tests for vanilla CRUD operations because they would
	// test Kubernetes API server, which isn't the goal here.
	Context("Group", func() {
		It("Full reconciliation", func() {

			updatedLocations := []string{"eu-west-2", "eu-west-1"}
			groupKey := types.NamespacedName{
				Name:      "test-group",
				Namespace: "default",
			}

			group := &checklyv1alpha1.Group{
				ObjectMeta: metav1.ObjectMeta{
					Name:      groupKey.Name,
					Namespace: groupKey.Namespace,
				},
				Spec: checklyv1alpha1.GroupSpec{
					Locations: []string{"eu-west-1"},
				},
			}

			// Create
			Expect(k8sClient.Create(context.Background(), group)).Should(Succeed())

			By("Expecting submitted")
			Eventually(func() bool {
				f := &checklyv1alpha1.Group{}
				err := k8sClient.Get(context.Background(), groupKey, f)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			// Status.ID should be present
			By("Expecting group ID")
			Eventually(func() bool {
				f := &checklyv1alpha1.Group{}
				err := k8sClient.Get(context.Background(), groupKey, f)
				if f.Status.ID == 1 && err == nil {
					return true
				} else {
					return false
				}
			}, timeout, interval).Should(BeTrue())

			// Finalizer should be present
			By("Expecting finalizer")
			Eventually(func() bool {
				f := &checklyv1alpha1.Group{}
				err := k8sClient.Get(context.Background(), groupKey, f)
				if len(f.Finalizers) == 1 && err == nil {
					return true
				} else {
					return false
				}
			}, timeout, interval).Should(BeTrue())

			// Update
			updated := &checklyv1alpha1.Group{}
			Expect(k8sClient.Get(context.Background(), groupKey, updated)).Should(Succeed())

			updated.Spec.Locations = updatedLocations
			Expect(k8sClient.Update(context.Background(), updated)).Should(Succeed())

			By("Expecting update")
			Eventually(func() bool {
				f := &checklyv1alpha1.Group{}
				err := k8sClient.Get(context.Background(), groupKey, f)
				if len(f.Spec.Locations) == 2 && err == nil {
					return true
				} else {
					return false
				}
			}, timeout, interval).Should(BeTrue())

			// Delete
			By("Expecting to delete successfully")
			Eventually(func() error {
				f := &checklyv1alpha1.Group{}
				k8sClient.Get(context.Background(), groupKey, f)
				return k8sClient.Delete(context.Background(), f)
			}, timeout, interval).Should(Succeed())

			By("Expecting delete to finish")
			Eventually(func() error {
				f := &checklyv1alpha1.Group{}
				return k8sClient.Get(context.Background(), groupKey, f)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})
})
