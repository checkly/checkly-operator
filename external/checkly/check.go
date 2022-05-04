package external

import (
	"context"
	"net/http"
	"time"

	"github.com/checkly/checkly-go-sdk"
)

// Check is a struct for the internal packages to help put together the checkly check
type Check struct {
	Name            string
	Namespace       string
	Frequency       int
	MaxResponseTime int
	Locations       []string
	Endpoint        string
	SuccessCode     string
	GroupID         int64
	ID              string
}

func checklyCheck(apiCheck Check) (check checkly.Check) {

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

	check = checkly.Check{
		Name:                 apiCheck.Name,
		Type:                 checkly.TypeAPI,
		Frequency:            checkValueInt(apiCheck.Frequency, 5),
		DegradedResponseTime: 5000,
		MaxResponseTime:      checkValueInt(apiCheck.MaxResponseTime, 15000),
		Activated:            true,
		Muted:                true, // muted for development
		ShouldFail:           false,
		DoubleCheck:          false,
		SSLCheck:             false,
		LocalSetupScript:     "",
		LocalTearDownScript:  "",
		Locations:            checkValueArray(apiCheck.Locations, []string{"eu-west-1"}),
		Tags: []string{
			apiCheck.Namespace,
			"checkly-operator",
		},
		AlertSettings:          alertSettings,
		UseGlobalAlertSettings: false,
		GroupID:                apiCheck.GroupID,
		Request: checkly.Request{
			Method:  http.MethodGet,
			URL:     apiCheck.Endpoint,
			Headers: []checkly.KeyValue{
				// {
				// 	Key:   "X-Test",
				// 	Value: "foo",
				// },
			},
			QueryParameters: []checkly.KeyValue{
				// {
				// 	Key:   "query",
				// 	Value: "foo",
				// },
			},
			Assertions: []checkly.Assertion{
				{
					Source:     checkly.StatusCode,
					Comparison: checkly.Equals,
					Target:     apiCheck.SuccessCode,
				},
			},
			Body:     "",
			BodyType: "NONE",
		},
	}

	return
}

// Create creates a new checkly.com check
func Create(apiCheck Check, client checkly.Client) (ID string, err error) {

	check := checklyCheck(apiCheck)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	gotCheck, err := client.Create(ctx, check)
	if err != nil {
		return
	}

	ID = gotCheck.ID

	return
}

// Update updates an existing checkly.com check
func Update(apiCheck Check, client checkly.Client) (err error) {

	check := checklyCheck(apiCheck)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.Update(ctx, apiCheck.ID, check)
	if err != nil {
		return
	}

	return
}

// Delete deletes an existing checkly.com check
func Delete(ID string, client checkly.Client) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = client.Delete(ctx, ID)

	return
}
