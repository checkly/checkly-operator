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
	"time"

	"github.com/checkly/checkly-go-sdk"
)

type Group struct {
	Name          string
	ID            int64
	Locations     []string
	Activated     bool
	AlertChannels []checkly.AlertChannelSubscription
	Labels        map[string]string
}

func checklyGroup(group Group) (check checkly.Group) {

	tags := getTags(group.Labels)
	tags = append(tags, "checkly-operator")

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

	check = checkly.Group{
		Name:                      group.Name,
		Activated:                 true,
		Muted:                     false, // muted for development
		DoubleCheck:               false,
		LocalSetupScript:          "",
		LocalTearDownScript:       "",
		Concurrency:               2,
		Locations:                 checkValueArray(group.Locations, []string{"eu-west-1"}),
		Tags:                      tags,
		AlertSettings:             alertSettings,
		UseGlobalAlertSettings:    false,
		AlertChannelSubscriptions: group.AlertChannels,
	}

	return
}

func GroupCreate(group Group, client checkly.Client) (ID int64, err error) {
	groupSetup := checklyGroup(group)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	gotGroup, err := client.CreateGroup(ctx, groupSetup)
	if err != nil {
		return
	}

	ID = gotGroup.ID

	return
}

func GroupUpdate(group Group, client checkly.Client) (err error) {

	groupSetup := checklyGroup(group)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.UpdateGroup(ctx, group.ID, groupSetup)
	if err != nil {
		return
	}

	return
}

func GroupDelete(ID int64, client checkly.Client) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = client.DeleteGroup(ctx, ID)

	return
}
