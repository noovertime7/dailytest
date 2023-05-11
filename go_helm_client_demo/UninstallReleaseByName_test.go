package main

import (
	"bytes"
	"fmt"
	helmSdk "github.com/mittwald/go-helm-client"
	"io"
	"testing"
	"time"
)

func Test_UninstallReleaseByName_test(t *testing.T) {
	var outputBuffer bytes.Buffer
	opt := &helmSdk.Options{
		Namespace: "default", // Change this to the namespace you wish the client to operate in.
		Debug:     true,
		Linting:   true,
		Output:    &outputBuffer,
	}
	go func() {
		for {
			byteMsg, err := outputBuffer.ReadByte()
			if err != nil {
				if err == io.EOF {
					t.Log("eof")
					return
				}
				t.Errorf(err.Error())
				return
			}
			fmt.Println(string(byteMsg))
		}

	}()

	helmClient, err := helmSdk.New(opt)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = helmClient.UninstallReleaseByName("minio")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log("uninstall success")

	time.Sleep(10 * time.Minute)
}
