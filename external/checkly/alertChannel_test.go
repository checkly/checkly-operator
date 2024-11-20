package external

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/checkly/checkly-go-sdk"
	checklyv1alpha1 "github.com/checkly/checkly-operator/api/checkly/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestChecklyAlertChannel(t *testing.T) {
	acName := "foo"
	acEmailAddress := "foo@bar.baz"

	dataEmpty := checklyv1alpha1.AlertChannel{
		ObjectMeta: metav1.ObjectMeta{
			Name: acName,
		},
		Spec: checklyv1alpha1.AlertChannelSpec{
			SendRecovery: false,
		},
	}

	opsGenieConfigEmpty := checkly.AlertChannelOpsgenie{}
	webhookConfigEmpty := checkly.AlertChannelWebhook{}

	returned, err := checklyAlertChannel(&dataEmpty, opsGenieConfigEmpty, webhookConfigEmpty)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}

	if returned.Opsgenie != nil {
		t.Errorf("Expected empty Opsgenie config, got %s", returned.Opsgenie)
	}

	dataEmail := dataEmpty
	dataEmail.Spec.Email = checkly.AlertChannelEmail{
		Address: acEmailAddress,
	}

	returned, err = checklyAlertChannel(&dataEmail, opsGenieConfigEmpty, webhookConfigEmpty)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}

	if returned.Email.Address != acEmailAddress {
		t.Errorf("Expected %s, got %s", acEmailAddress, returned.Email.Address)
	}

	dataOpsGenieFull := checkly.AlertChannelOpsgenie{
		APIKey:   "foo-bar",
		Region:   "US",
		Priority: "999",
		Name:     "baz",
	}

	returned, err = checklyAlertChannel(&dataEmpty, dataOpsGenieFull, webhookConfigEmpty)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}

	if returned.Opsgenie == nil {
		t.Error("Expected Opsgenie field to tbe populated, it's empty")
	}

	if returned.Opsgenie.Priority != "999" {
		t.Errorf("Expected %s, got %s", "999", returned.Opsgenie.Priority)
	}

	if returned.Opsgenie.Region != "US" {
		t.Errorf("Expected %s, got %s", "US", returned.Opsgenie.Region)
	}

	if returned.Email != nil {
		t.Errorf("Expected nil, got %s", returned.Email)
	}

	if returned.Webhook != nil { // Can't test against nil because []KeyValue pairs are present
		t.Errorf("Expected nil, got %v+", returned.Webhook)
	}

	dataWebhookFull := checkly.AlertChannelWebhook{
		Name:            "test",
		URL:             "http://foo.bar",
		WebhookType:     "GET",
		Method:          "POST",
		Template:        "",
		WebhookSecret:   "foobar",
		Headers:         []checkly.KeyValue{},
		QueryParameters: []checkly.KeyValue{},
	}

	returned, err = checklyAlertChannel(&dataEmpty, opsGenieConfigEmpty, dataWebhookFull)
	if err != nil {
		t.Errorf("Expected error, got %e", err)
	}

	if returned.Webhook == nil {
		t.Errorf("Expected Webhook field to be populated, it's empty")
	}

	if returned.Webhook.Method != "POST" {
		t.Errorf("Expected %s, got %s", "POST", returned.Webhook.Method)
	}

	if returned.Opsgenie != nil {
		t.Errorf("Expected nil, got %s", returned.Opsgenie)
	}

	if returned.Email != nil {
		t.Errorf("Expected nil, got %s", returned.Email)
	}

}

func TestAlertChannelActions(t *testing.T) {
	// Generate a different number each time
	rand.Seed(time.Now().UnixNano())
	expectedAlertChannelID := rand.Intn(100)

	acName := "foo"

	testData := &checklyv1alpha1.AlertChannel{
		ObjectMeta: metav1.ObjectMeta{
			Name: acName,
		},
		Spec: checklyv1alpha1.AlertChannelSpec{
			SendRecovery: false,
		},
		Status: checklyv1alpha1.AlertChannelStatus{
			ID: int64(expectedAlertChannelID),
		},
	}

	opsGenieConfigEmpty := checkly.AlertChannelOpsgenie{}
	webhookConfigEmpty := checkly.AlertChannelWebhook{}

	// Test errors
	testClient := checkly.NewClient(
		"http://localhost:5557",
		"foobarbaz",
		nil,
		nil,
	)
	testClient.SetAccountId("1234567890")

	// Create fail
	_, err := CreateAlertChannel(testData, opsGenieConfigEmpty, webhookConfigEmpty, testClient)
	if err == nil {
		t.Error("Expected error, got none")
	}

	// Update fail
	err = UpdateAlertChannel(testData, opsGenieConfigEmpty, webhookConfigEmpty, testClient)
	if err == nil {
		t.Error("Expected error, got none")
	}

	// Delete fail
	err = DeleteAlertChannel(testData, testClient)
	if err == nil {
		t.Error("Expected error, got none")
	}

	go func() {
		http.HandleFunc("/v1/alert-channels", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			resp := make(map[string]interface{})
			resp["id"] = expectedAlertChannelID
			jsonResp, _ := json.Marshal(resp)
			w.Write(jsonResp)
			return
		})
		http.HandleFunc(fmt.Sprintf("/v1/alert-channels/%d", expectedAlertChannelID), func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			method := r.Method
			switch method {
			case "PUT":
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				resp := make(map[string]interface{})
				resp["id"] = expectedAlertChannelID
				jsonResp, _ := json.Marshal(resp)
				w.Write(jsonResp)
			case "DELETE":
				w.WriteHeader(http.StatusNoContent)
			}
			return
		})
		http.ListenAndServe(":5557", nil)
	}()

	// Create success
	testID, err := CreateAlertChannel(testData, opsGenieConfigEmpty, webhookConfigEmpty, testClient)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}
	if testID != int64(expectedAlertChannelID) {
		t.Errorf("Expected %d, got %d", testID, int64(expectedAlertChannelID))
	}

	// Update success
	err = UpdateAlertChannel(testData, opsGenieConfigEmpty, webhookConfigEmpty, testClient)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}

	// Delete success
	err = DeleteAlertChannel(testData, testClient)
	if err != nil {
		t.Errorf("Expecte no error, got %e", err)
	}

}
