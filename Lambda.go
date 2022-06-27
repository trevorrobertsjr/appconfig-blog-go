package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfigdata"
	"github.com/aws/jsii-runtime-go"
)

// func handler(request events.S3Entity) (string, error) {
// 	log.Println("Hello World")

// 	// Read feature flags and modify what is returned accordingly
// 	resp, err := http.Get("http://localhost:2772/applications/appconfig-blog-go/environments/dev/configurations/WhichSide")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//We Read the response body on the line below.
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//Convert the body to type string
// 	sb := string(body)
// 	// log.Printf(sb)
// 	// Read DynamoDB
// 	// Return html
// 	return fmt.Sprintf(sb), nil
// }
// feature flag format:
// {"allegiance":{"choice":"darkknight","enabled":true}}
// Struct for unmarshall simplifies outputting field data as string in the Lambda
// may simplify to the map version of the implementation later
type featureflagdata struct {
	Choice  string
	Enabled bool
}
type featureflag struct {
	Allegiance featureflagdata
}

const (
	ApplicationIdentifier          string = "appconfig-blog-go"
	ConfigurationProfileIdentifier string = "WhichSide"
	EnvironmentIdentifier          string = "dev"
)

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (string, error) {
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

	// Option if not defining a struct in advance for unmarshall

	// var featureFlagResults map[string]interface{}
	// json.Unmarshal(result.Configuration, &featureFlagResults)
	// CecilsChoice := featureFlagResults["allegiance"].(map[string]interface{})
	// fmt.Println(CecilsChoice["choice"])

	var featureFlagResults featureflag
	json.Unmarshal(result.Configuration, &featureFlagResults)
	CecilsChoice := featureFlagResults.Allegiance.Choice
	fmt.Println(CecilsChoice)
	log.Println(CecilsChoice)
	return fmt.Sprintf(CecilsChoice), nil
}

func main() {
	lambda.Start(handler)
}
