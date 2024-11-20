package external

import (
	"context"
	"time"

	"github.com/checkly/checkly-go-sdk"
	checklyv1alpha1 "github.com/checkly/checkly-operator/api/checkly/v1alpha1"
)

func checklyAlertChannel(alertChannel *checklyv1alpha1.AlertChannel, opsGenieConfig checkly.AlertChannelOpsgenie, webhookConfig checkly.AlertChannelWebhook) (ac checkly.AlertChannel, err error) {
	ac = checkly.AlertChannel{
		SendRecovery: &alertChannel.Spec.SendRecovery,
		SendFailure:  &alertChannel.Spec.SendFailure,
		SendDegraded: &alertChannel.Spec.SendDegraded,
		SSLExpiry:    &alertChannel.Spec.SSLExpiry,
	}

	if (alertChannel.Spec.SSLExpiryThreshold > 0) && (alertChannel.Spec.SSLExpiryThreshold < 30) {
		ac.SSLExpiryThreshold = &alertChannel.Spec.SSLExpiryThreshold
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

	if webhookConfig.Name != "" { // Struct has []KeyValue types which can't be compared
		ac.Type = "WEBHOOK" // Type has to be all caps, see https://developers.checklyhq.com/reference/postv1alertchannels
		ac.Webhook = &webhookConfig
	}
	return
}

func CreateAlertChannel(alertChannel *checklyv1alpha1.AlertChannel, opsGenieConfig checkly.AlertChannelOpsgenie, webhookConfig checkly.AlertChannelWebhook, client checkly.Client) (ID int64, err error) {

	ac, err := checklyAlertChannel(alertChannel, opsGenieConfig, webhookConfig)
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

func UpdateAlertChannel(alertChannel *checklyv1alpha1.AlertChannel, opsGenieConfig checkly.AlertChannelOpsgenie, webhookConfig checkly.AlertChannelWebhook, client checkly.Client) (err error) {
	ac, err := checklyAlertChannel(alertChannel, opsGenieConfig, webhookConfig)
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
