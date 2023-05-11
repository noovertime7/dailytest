package main

import (
	helmSdk "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/action"
	"testing"
)

func Test_ListDeployedReleases(t *testing.T) {
	opt := &helmSdk.Options{
		Namespace: "default", // Change this to the namespace you wish the client to operate in.
		Debug:     true,
		Linting:   true,
	}
	helmClient, err := helmSdk.New(opt)
	if err != nil {
		t.Fatalf(err.Error())
	}

	releases, err := helmClient.ListReleasesByStateMask(action.ListFailed)
	if err != nil {
		t.Fatalf(err.Error())
	}
	for _, release := range releases {
		t.Log(release)
	}
}
