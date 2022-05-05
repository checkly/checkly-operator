package external

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/checkly/checkly-go-sdk"
)

func TestChecklyCheck(t *testing.T) {

	data := Check{
		Name:            "foo",
		Namespace:       "bar",
		Frequency:       15,
		MaxResponseTime: 2000,
		Locations:       []string{"basement"},
		Endpoint:        "https://foo.bar/baz",
		SuccessCode:     "200",
	}

	testData := checklyCheck(data)

	if testData.Name != data.Name {
		t.Errorf("Expected %s, got %s", data.Name, testData.Name)
	}

	if testData.Frequency != data.Frequency {
		t.Errorf("Expected %d, got %d", data.Frequency, testData.Frequency)
	}

	if testData.MaxResponseTime != data.MaxResponseTime {
		t.Errorf("Expected %d, got %d", data.Frequency, testData.Frequency)
	}

	data = Check{
		Name:        "foo",
		Namespace:   "bar",
		Endpoint:    "https://foo.bar/baz",
		SuccessCode: "200",
	}

	testData = checklyCheck(data)

	if testData.Frequency != 5 {
		t.Errorf("Expected %d, got %d", data.Frequency, testData.Frequency)
	}

	if testData.MaxResponseTime != 15000 {
		t.Errorf("Expected %d, got %d", data.Frequency, testData.Frequency)
	}

	return
}

func TestChecklyCheckActions(t *testing.T) {

	expectedCheckID := "2"
	expectedGroupID := 1
	testData := Check{
		Name:            "foo",
		Namespace:       "bar",
		Frequency:       15,
		MaxResponseTime: 2000,
		Locations:       []string{"basement"},
		Endpoint:        "https://foo.bar/baz",
		SuccessCode:     "200",
		ID:              "",
	}

	// Test errors
	testClientFail := checkly.NewClient(
		"http://localhost:5556",
		"foobarbaz",
		nil,
		nil,
	)
	// Create
	_, err := Create(testData, testClientFail)
	if err == nil {
		t.Error("Expected error, got none")
	}

	// Update
	err = Update(testData, testClientFail)
	if err == nil {
		t.Error("Expected error, got none")
	}

	// Delete
	err = Delete(expectedCheckID, testClientFail)
	if err == nil {
		t.Error("Expected error, got none")
	}

	// Test happy path
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
			resp["id"] = expectedCheckID
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
				resp["id"] = expectedCheckID
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
			resp["id"] = expectedGroupID
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
				resp["id"] = expectedGroupID
				jsonResp, _ := json.Marshal(resp)
				w.Write(jsonResp)
			case "DELETE":
				w.WriteHeader(http.StatusNoContent)
			}
			return
		})
		http.ListenAndServe(":5555", nil)
	}()

	testID, err := Create(testData, testClient)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}

	if testID != expectedCheckID {
		t.Errorf("Expected %s, got %s", expectedCheckID, testID)
	}

	testData.ID = expectedCheckID

	err = Update(testData, testClient)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}

	err = Delete(expectedCheckID, testClient)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}

	return
}
