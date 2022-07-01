package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfigdata"
	"github.com/aws/jsii-runtime-go"
)

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

	// Option if not defining a struct in advance for unmarshall

	// var featureFlagResults map[string]interface{}
	// json.Unmarshal(result.Configuration, &featureFlagResults)
	// CecilsChoice := featureFlagResults["allegiance"].(map[string]interface{})
	// fmt.Println(CecilsChoice["choice"])
	htmlDarkKnightOutput := `			<section>
	<img src="assets/img/dark_knight_cecil.png" alt="image of Cecil as a Dark Knight">
</section>
<section>
	<p>Cecil chose the path of the Dark Knight. The new and unfamiliar path of the alternative was too
	 discomforting for him to stray from the path set out for him by his adopted family. Cecil's mother
	 was disappointed, but she ultimately understood.</p>
</section>`
	htmlPaladinOutput := `			<section>
	<img src="assets/img/paladin_cecil.png" alt="image of Cecil as a Paladin">
</section>
<section>
	<p>Cecil heeded the words of his mother and took his rightful place as the prophesied Paladin.
	 Contrary to his fears, selecting the light has grown his power exponentially. He then realized his former general
	 had him take the dark sword to hold him back rather than to empower him.</p>
</section>`
	fmt.Println(htmlDarkKnightOutput)
	fmt.Println(htmlPaladinOutput)

	var featureFlagResults featureflag
	json.Unmarshal(result.Configuration, &featureFlagResults)
	CecilsChoice := featureFlagResults.Allegiance.Choice
	fmt.Println(CecilsChoice)
}
