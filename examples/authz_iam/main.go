// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "github.com/hacku7/component-base-sdk/pkg/meta/v1"
	"github.com/hacku7/component-base-sdk/pkg/util/homedir"
	"github.com/hacku7/hbird-sdk/hbird/service/iam"
	"github.com/hacku7/hbird-sdk/tools/clientcmd"
	"github.com/ory/ladon"
	"path/filepath"
)

func main() {
	var iamconfig *string
	if home := homedir.HomeDir(); home != "" {
		iamconfig = flag.String(
			"iamconfig",
			filepath.Join(home, ".iam", "config"),
			"(optional) absolute path to the iamconfig file",
		)
	} else {
		iamconfig = flag.String("iamconfig", "", "absolute path to the iamconfig file")
	}
	flag.Parse()

	// use the current context in iamconfig
	config, err := clientcmd.BuildConfigFromFlags("", *iamconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the iamclient
	iamclient, err := iam.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	authzClient := iamclient.AuthzV1().Authz()

	request := &ladon.Request{
		Resource: "resources:articles:ladon-introduction",
		Action:   "delete",
		Subject:  "users:peter",
		Context: ladon.Context{
			"remoteIP": "192.168.0.5",
		},
	}

	// Authorize the request
	fmt.Println("Authorize request...")
	ret, err := authzClient.Authorize(context.TODO(), request, metav1.AuthorizeOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Authorize response: %s.\n", ret.ToString())
}
