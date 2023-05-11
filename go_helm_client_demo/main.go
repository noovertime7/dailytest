package main

import (
	"context"
	helmSdk "github.com/mittwald/go-helm-client"
	"helm.sh/helm/v3/pkg/repo"
	"log"
)

func main() {
	opt := &helmSdk.Options{
		Namespace: "default", // Change this to the namespace you wish the client to operate in.
		Debug:     true,
		Linting:   true,
	}

	helmClient, err := helmSdk.New(opt)
	if err != nil {
		panic(err)
	}

	//release, err := helmClient.GetRelease("minio")
	//if err != nil {
	//	log.Fatalf("GetRelease error %v", err)
	//}
	//fmt.Println(release)

	err = helmClient.AddOrUpdateChartRepo(repo.Entry{
		Name:               "chartmuseum",
		URL:                "http://localhost:8080",
		PassCredentialsAll: true,
	})
	if err != nil {
		log.Fatalf("AddOrUpdateChartRepo error %v", err)
	}
	log.Println("AddOrUpdateChartRepo success")

	err = helmClient.UpdateChartRepos()
	if err != nil {
		log.Fatalf("UpdateChartRepos error %v", err)
	}

	chartSpec := helmSdk.ChartSpec{
		ReleaseName: "minio",
		ChartName:   "chartmuseum/minio",
		Namespace:   "default",
		UpgradeCRDs: true,
		Wait:        true,
	}
	realse, err := helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec, nil)
	if err != nil {
		log.Fatalf("InstallChart error %v", err)
	}
	log.Println(realse)

}
