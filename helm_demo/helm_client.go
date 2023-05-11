package main

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"log"
	"os"
)

func main() {
	chartPath := "./charts.tgz"
	namespace := "default"
	releaseName := "minio"

	settings := cli.New()

	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace,
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	// define values
	vals := map[string]interface{}{
		"namespace":    namespace,
		"replicaCount": "1",
		"image": map[string]interface{}{
			"name":       "image",
			"pullPolicy": "IfNotPresent",
		},
		"nameOverride":     releaseName,
		"fullnameOverride": releaseName,
		"service": map[string]interface{}{
			"type": "ClusterIP",
			"port": "8080",
		},
		"hostAliases": map[string]interface{}{
			"hostAliasesSet": false,
			"hostAliasesMap": "",
		},
		"imagePull": map[string]interface{}{
			"secretsSet": true,
			"secrets":    "secret-harbor",
		},
		"envConfig": map[string]interface{}{
			"configSet": false,
			"config":    "",
		},
		"resources": map[string]interface{}{
			"resourcesSet": false,
			"requests": map[string]interface{}{
				"memory": "",
				"cpu":    "",
			},
			"limits": map[string]interface{}{
				"memory": "",
				"cpu":    "",
			},
			"affinity": map[string]interface{}{
				"affinitySet": false,
				"key":         "",
				"operator":    "In",
				"values":      "",
			},
			"tolerations": map[string]interface{}{
				"tolerationsSet": false,
				"key":            "__TOLERATIONS_KEY__",
				"operator":       "Equal",
				"value":          "__TOLERATIONS_VALUE__",
				"effect":         "NoSchedule",
			},
			"livenessProbe": map[string]interface{}{
				"livenessGetSet":     false,
				"livenessGet":        "__LIVENESS_GET__",
				"livenessCommandSet": false,
				"livenessCommand":    "__LIVENESS_COMMAND__",
			},
			"command": map[string]interface{}{
				"commandSet":   false,
				"startCommand": "__COMMAND__",
			},
			"AnnoTations": map[string]interface{}{
				"AnnoTationSet": false,
			},
		},
	}

	// load chart from the path
	chart, err := loader.Load(chartPath)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = namespace
	client.ReleaseName = releaseName
	// client.DryRun = true - very handy!

	// install the chart here
	rel, err := client.Run(chart, vals)
	if err != nil {
		panic(err)
	}

	log.Printf("Installed Chart from path: %s in namespace: %s\n", rel.Name, rel.Namespace)
	// this will confirm the values set during installation
	log.Println(rel.Config)
}
