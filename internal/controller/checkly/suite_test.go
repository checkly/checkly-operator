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

package checkly

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/checkly/checkly-go-sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	checklyv1alpha1 "github.com/checkly/checkly-operator/apis/checkly/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	Expect(os.Setenv("USE_EXISTING_CLUSTER", "true")).To(Succeed())
	Expect(os.Setenv("TEST_ASSET_KUBECTL", "../testbin/bin/kubectl")).To(Succeed())
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = checklyv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	// Stub checkly client
	testClient := checkly.NewClient(
		"http://localhost:5555",
		"foobarbaz",
		nil,
		nil,
	)
	testClient.SetAccountId("1234567890")
	go func() {
		http.HandleFunc("/v1/checks", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]string)
			resp["id"] = "2"
			jsonResp, _ := json.Marshal(resp)
			w.Write(jsonResp)
			return
		})
		http.HandleFunc("/v1/checks/2", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			method := r.Method
			switch method {
			case "PUT":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				resp := make(map[string]string)
				resp["id"] = "2"
				jsonResp, _ := json.Marshal(resp)
				w.Write(jsonResp)
			case "DELETE":
				w.WriteHeader(http.StatusNoContent)
			}
			return
		})
		http.HandleFunc("/v1/check-groups", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]interface{})
			resp["id"] = 1
			jsonResp, _ := json.Marshal(resp)
			w.Write(jsonResp)
			return
		})
		http.HandleFunc("/v1/check-groups/1", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			method := r.Method
			switch method {
			case "PUT":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				resp := make(map[string]interface{})
				resp["id"] = 1
				jsonResp, _ := json.Marshal(resp)
				w.Write(jsonResp)
			case "DELETE":
				w.WriteHeader(http.StatusNoContent)
			}
			return
		})
		http.HandleFunc("/v1/alert-channels", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]interface{})
			resp["id"] = 3
			jsonResp, _ := json.Marshal(resp)
			w.Write(jsonResp)
			return
		})
		http.HandleFunc("/v1/alert-channels/3", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			method := r.Method
			switch method {
			case "PUT":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				resp := make(map[string]interface{})
				resp["id"] = 3
				jsonResp, _ := json.Marshal(resp)
				w.Write(jsonResp)
			case "DELETE":
				w.WriteHeader(http.StatusNoContent)
			}
			return
		})
		http.ListenAndServe(":5555", nil)
	}()

	err = (&ApiCheckReconciler{
		Client:    k8sManager.GetClient(),
		Scheme:    k8sManager.GetScheme(),
		ApiClient: testClient,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&GroupReconciler{
		Client:    k8sManager.GetClient(),
		Scheme:    k8sManager.GetScheme(),
		ApiClient: testClient,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&AlertChannelReconciler{
		Client:    k8sManager.GetClient(),
		Scheme:    k8sManager.GetScheme(),
		ApiClient: testClient,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
