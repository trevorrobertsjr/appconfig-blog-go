package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfigdata"
	"github.com/aws/jsii-runtime-go"
)

const (
	ApplicationIdentifier          string = "appconfig-blog-go"
	ConfigurationProfileIdentifier string = "WhichSide"
	EnvironmentIdentifier          string = "dev"
)

func main() {
	// Specifying the region when creating the AWS Session
	mySession, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		fmt.Println(err)
	}

	// Create a AppConfigData client from the AWS Session.
	svc := appconfigdata.New(mySession)

	// Retrieve an AppConfig token to then make the request
	// for the latest version of the feature flag
	token, err := svc.StartConfigurationSession(&appconfigdata.StartConfigurationSessionInput{
		// The ApplicationIdentifier for the application that depends on the feature flag.
		ApplicationIdentifier: jsii.String(ApplicationIdentifier),

		// The configuration profile ID or the configuration profile name.
		ConfigurationProfileIdentifier: jsii.String(ConfigurationProfileIdentifier),

		// The AppConfig environment ID (ex: dev, beta, prod, etc.).
		EnvironmentIdentifier: jsii.String(EnvironmentIdentifier),
	})
	if err != nil {
		fmt.Println(err)
	}
	result, err := svc.GetLatestConfiguration(&appconfigdata.GetLatestConfigurationInput{
		ConfigurationToken: jsii.String(*token.InitialConfigurationToken),
	})
	if err != nil {
		fmt.Println(err)
	}
	// resultStruct := result.Configuration
	var results map[string]interface{}
	json.Unmarshal(result.Configuration, &results)
	birds := results["allegiance"].(map[string]interface{})
	// json.Unmarshal(result.Configuration, &resultStruct)
	fmt.Println(birds["choice"])
}
