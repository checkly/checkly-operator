package external

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/checkly/checkly-go-sdk"
)

func checklyClient() (client checkly.Client, ctx context.Context, cancel context.CancelFunc, err error) {
	baseUrl := "https://api.checklyhq.com"
	apiKey := os.Getenv("CHECKLY_API_KEY")
	if apiKey == "" {
		err = errors.New("Checkly.com API key environment variable is undefined")
		return
	}

	accountId := os.Getenv("CHECKLY_ACCOUNT_ID")
	if accountId == "" {
		err = errors.New("Checkly.com Account ID environment variable is undefined")
		return
	}

	client = checkly.NewClient(
		baseUrl,
		apiKey,
		nil, //custom http client, defaults to http.DefaultClient
		nil, //io.Writer to output debug messages
	)

	client.SetAccountId(accountId)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	return
}

func checkValueString(x string, y string) (value string) {
	if x == "" {
		value = y
	} else {
		value = x
	}
	return
}

func checkValueInt(x int, y int) (value int) {
	if x == 0 {
		value = y
	} else {
		value = x
	}
	return
}

func checkValueArray(x []string, y []string) (value []string) {
	if len(x) == 0 {
		value = y
	} else {
		value = x
	}
	return
}
