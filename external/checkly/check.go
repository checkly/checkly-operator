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
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/checkly/checkly-go-sdk"
)

// Check is a struct for the internal packages to help put together the checkly check
type Check struct {
	Name            string
	Namespace       string
	Frequency       int
	MaxResponseTime int
	Endpoint        string
	SuccessCode     string
	GroupID         int64
	ID              string
	Muted           bool
	Labels          map[string]string
	Assertions      []checkly.Assertion // Align with checkly's Assertion struct
	Method          string              // HTTP method to use for the check
}

func checklyCheck(apiCheck Check) (check checkly.Check, err error) {
	// Ensure `shouldFail` logic is handled
	shouldFail, err := shouldFail(apiCheck.SuccessCode)
	if err != nil {
		return
	}

	// Map tags from labels and namespace
	tags := getTags(apiCheck.Labels)
	tags = append(tags, "checkly-operator", apiCheck.Namespace)

	// Define alert settings
	alertSettings := checkly.AlertSettings{
		EscalationType: checkly.RunBased,
		RunBasedEscalation: checkly.RunBasedEscalation{
			FailedRunThreshold: 5,
		},
		TimeBasedEscalation: checkly.TimeBasedEscalation{
			MinutesFailingThreshold: 5,
		},
		Reminders: checkly.Reminders{
			Interval: 5,
		},
		SSLCertificates: checkly.SSLCertificates{
			Enabled:        false,
			AlertThreshold: 3,
		},
	}

	// Default assertion logic: If no assertions are provided, add a default
	assertions := apiCheck.Assertions
	if len(assertions) == 0 {
		assertions = []checkly.Assertion{
			{
				Source:     checkly.StatusCode,
				Comparison: checkly.Equals,
				Target:     apiCheck.SuccessCode,
			},
		}
	}

	// Determine the HTTP method; default to GET if not specified
	method := http.MethodGet
	if apiCheck.Method != "" {
		method = apiCheck.Method
	}

	// Create the Checkly API check structure
	check = checkly.Check{
		Name:                 apiCheck.Name,
		Type:                 checkly.TypeAPI,
		Frequency:            checkValueInt(apiCheck.Frequency, 5),
		DegradedResponseTime: 5000,
		MaxResponseTime:      checkValueInt(apiCheck.MaxResponseTime, 15000),
		Activated:            true,
		Muted:                apiCheck.Muted,
		ShouldFail:           shouldFail,
		DoubleCheck:          false,
		SSLCheck:             false,
		AlertSettings:        alertSettings,
		Locations:            []string{},
		Tags:                 tags,
		Request: checkly.Request{
			Method:          method,
			URL:             apiCheck.Endpoint,
			Assertions:      assertions,
			Headers:         []checkly.KeyValue{},
			QueryParameters: []checkly.KeyValue{},
			Body:            "",
			BodyType:        "NONE",
		},
		UseGlobalAlertSettings: false,
		GroupID:                apiCheck.GroupID,
	}

	return
}

// Create creates a new checklyhq.com check
func Create(apiCheck Check, client checkly.Client) (ID string, err error) {

	check, err := checklyCheck(apiCheck)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	gotCheck, err := client.Create(ctx, check)
	if err != nil {
		return
	}

	ID = gotCheck.ID

	return
}

// Update updates an existing checklyhq.com check
func Update(apiCheck Check, client checkly.Client) (err error) {

	check, err := checklyCheck(apiCheck)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.Update(ctx, apiCheck.ID, check)

	return
}

// Delete deletes an existing checklyhq.com check
func Delete(ID string, client checkly.Client) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = client.Delete(ctx, ID)

	return
}

func shouldFail(successCode string) (bool, error) {
	code, err := strconv.Atoi(successCode)
	if err != nil {
		return false, err
	}
	if code < 400 {
		return false, nil
	} else {
		return true, nil
	}
}
