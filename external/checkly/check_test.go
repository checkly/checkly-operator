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

package external

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/checkly/checkly-go-sdk"
)

func TestChecklyCheck(t *testing.T) {

	data1 := Check{
		Name:            "foo",
		Namespace:       "bar",
		Frequency:       15,
		MaxResponseTime: 2000,
		Endpoint:        "https://foo.bar/baz",
		SuccessCode:     "403",
		Muted:           true,
	}

	testData, _ := checklyCheck(data1)

	if testData.Name != data1.Name {
		t.Errorf("Expected %s, got %s", data1.Name, testData.Name)
	}

	if testData.Frequency != data1.Frequency {
		t.Errorf("Expected %d, got %d", data1.Frequency, testData.Frequency)
	}

	if testData.MaxResponseTime != data1.MaxResponseTime {
		t.Errorf("Expected %d, got %d", data1.MaxResponseTime, testData.MaxResponseTime)
	}

	if testData.Muted != data1.Muted {
		t.Errorf("Expected %t, got %t", data1.Muted, testData.Muted)
	}

	if testData.ShouldFail != true {
		t.Errorf("Expected %t, got %t", true, testData.ShouldFail)
	}

	data2 := Check{
		Name:        "foo",
		Namespace:   "bar",
		Endpoint:    "https://foo.bar/baz",
		SuccessCode: "200",
	}

	testData, _ = checklyCheck(data2)

	if testData.Frequency != 5 {
		t.Errorf("Expected %d, got %d", 5, testData.Frequency)
	}

	if testData.MaxResponseTime != 15000 {
		t.Errorf("Expected %d, got %d", 15000, testData.MaxResponseTime)
	}

	if testData.ShouldFail != false {
		t.Errorf("Expected %t, got %t", false, testData.ShouldFail)
	}

	failData := Check{
		Name:        "fail",
		Namespace:   "bar",
		Endpoint:    "https://foo.bar/baz",
		SuccessCode: "foo",
	}

	_, err := checklyCheck(failData)
	if err == nil {
		t.Error("Expected error, got nil")
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

func TestShouldFail(t *testing.T) {
	testTrue := "401"
	testFalse := "200"
	testErr := "foo"

	testResponse, err := shouldFail(testTrue)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}
	if testResponse != true {
		t.Errorf("Expected true, got %t", testResponse)
	}

	testResponse, err = shouldFail(testFalse)
	if err != nil {
		t.Errorf("Expected no error, got %e", err)
	}
	if testResponse != false {
		t.Errorf("Expected false, got %t", testResponse)
	}

	_, err = shouldFail(testErr)
	if err == nil {
		t.Errorf("Expected error, got none")
	}

}
