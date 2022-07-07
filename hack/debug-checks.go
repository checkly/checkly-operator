package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/checkly/checkly-go-sdk"
	ctrl "sigs.k8s.io/controller-runtime"
)

// The script returns the data for a specific checkly check from the checklyhq API,
// it's meant to be used as a debug tool for issues with checks created.

func main() {

	setupLog := ctrl.Log.WithName("setup")

	var checklyID string

	flag.StringVar(&checklyID, "c", "", "Specify the checkly check ID")
	flag.Parse()

	if checklyID == "" {
		setupLog.Error(errors.New("ChecklyID is empty"), "exiting due to missing information")
		os.Exit(1)
	}

	baseUrl := "https://api.checklyhq.com"
	apiKey := os.Getenv("CHECKLY_API_KEY")
	if apiKey == "" {
		setupLog.Error(errors.New("checklyhq.com API key environment variable is undefined"), "checklyhq.com credentials missing")
		os.Exit(1)
	}

	accountId := os.Getenv("CHECKLY_ACCOUNT_ID")
	if accountId == "" {
		setupLog.Error(errors.New("checklyhq.com Account ID environment variable is undefined"), "checklyhq.com credentials missing")
		os.Exit(1)
	}

	client := checkly.NewClient(
		baseUrl,
		apiKey,
		nil, //custom http client, defaults to http.DefaultClient
		nil, //io.Writer to output debug messages
	)

	client.SetAccountId(accountId)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	returnedCheck, err := client.Get(ctx, checklyID)
	if err != nil {
		setupLog.Error(err, "failed to get check")
		os.Exit(1)
	}

	fmt.Printf("%+v", returnedCheck)

}
