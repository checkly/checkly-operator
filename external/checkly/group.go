package external

import (
	"github.com/checkly/checkly-go-sdk"
)

type Group struct {
	Name          string
	Namespace     string
	ID            int64
	Locations     []string
	Activated     bool
	AlertChannels []string
}

func checklyGroup(group Group) (check checkly.Group) {

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
		Name:                group.Name,
		Activated:           true,
		Muted:               true, // muted for development
		DoubleCheck:         false,
		LocalSetupScript:    "",
		LocalTearDownScript: "",
		Concurrency:         2,
		Locations:           checkValueArray(group.Locations, []string{"eu-west-1"}),
		Tags: []string{
			group.Namespace,
			"checkly-operator",
		},
		AlertSettings:             alertSettings,
		UseGlobalAlertSettings:    false,
		AlertChannelSubscriptions: []checkly.AlertChannelSubscription{},
	}

	return
}

func GroupCreate(group Group) (ID int64, err error) {
	groupSetup := checklyGroup(group)

	client, ctx, cancel, _ := checklyClient()
	defer cancel()

	gotGroup, err := client.CreateGroup(ctx, groupSetup)
	if err != nil {
		return
	}

	ID = gotGroup.ID

	return
}

func GroupUpdate(group Group) (err error) {

	groupSetup := checklyGroup(group)

	client, ctx, cancel, _ := checklyClient()
	defer cancel()

	_, err = client.UpdateGroup(ctx, group.ID, groupSetup)
	if err != nil {
		return
	}

	return
}

func GroupDelete(ID int64) (err error) {
	client, ctx, cancel, _ := checklyClient()
	defer cancel()

	err = client.DeleteGroup(ctx, ID)

	return
}
