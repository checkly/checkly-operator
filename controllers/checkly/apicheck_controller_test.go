// Example code used for influence: https://github.com/Azure/azure-databricks-operator/blob/0f722a710fea06b86ecdccd9455336ca712bf775/controllers/dcluster_controller_test.go

package checkly

// import (
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )

// var _ = Describe("ApiCheck Controller", func() {

// 	// Define utility constants for object names and testing timeouts/durations and intervals.
// 	const (
// 		timeout  = time.Second * 10
// 		duration = time.Second * 10
// 		interval = time.Millisecond * 250
// 	)

// 	BeforeEach(func() {
// 		// Add any setup steps that needs to be executed before each test
// 	})

// 	AfterEach(func() {
// 		// Add any teardown steps that needs to be executed after each test
// 	})

// 	// Add Tests for OpenAPI validation (or additonal CRD features) specified in
// 	// your API definition.
// 	// Avoid adding tests for vanilla CRUD operations because they would
// 	// test Kubernetes API server, which isn't the goal here.
// 	Context("ApiCheck without group", func() {
// 		It("Should create successfully", func() {

// 			key := types.NamespacedName{
// 				Name:      "test-apicheck",
// 				Namespace: "default",
// 			}

// 			groupKey := types.NamespacedName{
// 				Name:      "test-group",
// 				Namespace: "default",
// 			}

// 			apiCheck := &checklyv1alpha1.ApiCheck{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      key.Name,
// 					Namespace: key.Namespace,
// 				},
// 				Spec: checklyv1alpha1.ApiCheckSpec{
// 					Team:     groupKey.Name,
// 					Endpoint: "http://bar.baz/quoz",
// 					Success:  "200",
// 				},
// 			}

// 			group := &checklyv1alpha1.Group{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      groupKey.Name,
// 					Namespace: groupKey.Namespace,
// 				},
// 				Spec: checklyv1alpha1.GroupSpec{
// 					Locations: []string{"eu-west-1"},
// 				},
// 			}

// 			// Create
// 			Expect(k8sClient.Create(context.Background(), group)).Should(Succeed())
// 			Expect(k8sClient.Create(context.Background(), apiCheck)).Should(Succeed())

// 			By("Expecting submitted")
// 			Eventually(func() bool {
// 				f := &checklyv1alpha1.ApiCheck{}
// 				err := k8sClient.Get(context.Background(), key, f)
// 				if err != nil {
// 					return false
// 				}
// 				return true
// 			}, timeout, interval).Should(BeTrue())

// 			// // Delete
// 			// By("Expecting to delete successfully")
// 			// Eventually(func() error {
// 			// 	f := &databricksv1.Dcluster{}
// 			// 	k8sClient.Get(context.Background(), key, f)
// 			// 	return k8sClient.Delete(context.Background(), f)
// 			// }, timeout, interval).Should(Succeed())

// 			// By("Expecting to delete finish")
// 			// Eventually(func() error {
// 			// 	f := &databricksv1.Dcluster{}
// 			// 	return k8sClient.Get(context.Background(), key, f)
// 			// }, timeout, interval).ShouldNot(Succeed())
// 		})
// 	})
// })
