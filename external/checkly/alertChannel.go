package external

import (
	"context"
	"time"

	"github.com/checkly/checkly-go-sdk"
	checklyv1alpha1 "github.com/checkly/checkly-operator/apis/checkly/v1alpha1"
)

func checklyAlertChannel(alertChannel *checklyv1alpha1.AlertChannel, opsGenieConfig checkly.AlertChannelOpsgenie) (ac checkly.AlertChannel, err error) {
	sslExpiry := false

	ac = checkly.AlertChannel{
		SendRecovery: &alertChannel.Spec.SendRecovery,
		SendFailure:  &alertChannel.Spec.SendFailure,
		SSLExpiry:    &sslExpiry,
	}

	if opsGenieConfig != (checkly.AlertChannelOpsgenie{}) {
		ac.Type = "OPSGENIE" // Type has to be all caps, see https://developers.checklyhq.com/reference/postv1alertchannels
		ac.Opsgenie = &opsGenieConfig
		return
	}

	if alertChannel.Spec.Email != (checkly.AlertChannelEmail{}) {
		ac.Type = "EMAIL" // Type has to be all caps, see https://developers.checklyhq.com/reference/postv1alertchannels
		ac.Email = &checkly.AlertChannelEmail{
			Address: alertChannel.Spec.Email.Address,
		}
		return
	}
	return
}

func CreateAlertChannel(alertChannel *checklyv1alpha1.AlertChannel, opsGenieConfig checkly.AlertChannelOpsgenie, client checkly.Client) (ID int64, err error) {

	ac, err := checklyAlertChannel(alertChannel, opsGenieConfig)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	gotAlertChannel, err := client.CreateAlertChannel(ctx, ac)
	if err != nil {
		return
	}

	ID = gotAlertChannel.ID

	return
}

func UpdateAlertChannel(alertChannel *checklyv1alpha1.AlertChannel, opsGenieConfig checkly.AlertChannelOpsgenie, client checkly.Client) (err error) {
	ac, err := checklyAlertChannel(alertChannel, opsGenieConfig)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = client.UpdateAlertChannel(ctx, alertChannel.Status.ID, ac)
	if err != nil {
		return
	}

	return
}

func DeleteAlertChannel(alertChannel *checklyv1alpha1.AlertChannel, client checkly.Client) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = client.DeleteAlertChannel(ctx, alertChannel.Status.ID)
	if err != nil {
		return
	}

	return
}
