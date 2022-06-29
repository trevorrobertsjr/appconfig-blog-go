package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/appconfig"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		appres, err := appconfig.NewApplication(ctx, "blogapplication", &appconfig.ApplicationArgs{
			Description: pulumi.String("Example AppConfig Application"),
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Application"),
			},
		})
		if err != nil {
			return err
		}
		depstratres, err := appconfig.NewDeploymentStrategy(ctx, "blogdeploymentstrategy", &appconfig.DeploymentStrategyArgs{
			DeploymentDurationInMinutes: pulumi.Int(3),
			Description:                 pulumi.String("Example Deployment Strategy"),
			FinalBakeTimeInMinutes:      pulumi.Int(4),
			GrowthFactor:                pulumi.Float64(10),
			GrowthType:                  pulumi.String("LINEAR"),
			ReplicateTo:                 pulumi.String("NONE"),
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Deployment Strategy"),
			},
		})
		if err != nil {
			return err
		}
		envres, err := appconfig.NewEnvironment(ctx, "blogEnvironment", &appconfig.EnvironmentArgs{
			Description:   pulumi.String("Example AppConfig Environment"),
			ApplicationId: appres.ID(),
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Environment"),
			},
		})
		if err != nil {
			return err
		}
		cfgprofileres, err := appconfig.NewConfigurationProfile(ctx, "blogConfigurationProfile", &appconfig.ConfigurationProfileArgs{
			ApplicationId: appres.ID(),
			Description:   pulumi.String("Example Configuration Profile"),
			LocationUri:   pulumi.String("hosted"),
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Configuration Profile"),
			},
		})
		if err != nil {
			return err
		}
		tmpJSON0, err := json.Marshal(map[string]interface{}{
			"flags": map[string]interface{}{
				"foo": map[string]interface{}{
					"name": "foo",
					"_deprecation": map[string]interface{}{
						"status": "planned",
					},
				},
				"bar": map[string]interface{}{
					"name": "bar",
					"attributes": map[string]interface{}{
						"someAttribute": map[string]interface{}{
							"constraints": map[string]interface{}{
								"type":     "string",
								"required": true,
							},
						},
						"someOtherAttribute": map[string]interface{}{
							"constraints": map[string]interface{}{
								"type":     "number",
								"required": true,
							},
						},
					},
				},
			},
			"values": map[string]interface{}{
				"foo": map[string]interface{}{
					"enabled": "true",
				},
				"bar": map[string]interface{}{
					"enabled":            "true",
					"someAttribute":      "Hello World",
					"someOtherAttribute": 123,
				},
			},
			"version": "1",
		})
		if err != nil {
			return err
		}
		json0 := string(tmpJSON0)
		hostedCfgVersion, err := appconfig.NewHostedConfigurationVersion(ctx, "blogConfigurationVersion", &appconfig.HostedConfigurationVersionArgs{
			ApplicationId:          appres.ID(),
			ConfigurationProfileId: cfgprofileres.ConfigurationProfileId,
			Description:            pulumi.String("Example Freeform Hosted Configuration Version"),
			ContentType:            pulumi.String("application/json"),
			Content:                pulumi.String(json0),
		})
		if err != nil {
			return err
		}
		_, errDeploy := appconfig.NewDeployment(ctx, "blogDeployment", &appconfig.DeploymentArgs{
			ApplicationId:          appres.ID(),
			ConfigurationProfileId: cfgprofileres.ConfigurationProfileId,
			ConfigurationVersion:   pulumi.Sprintf("%v", hostedCfgVersion.VersionNumber),
			DeploymentStrategyId:   depstratres.ID(),
			Description:            pulumi.String("My example deployment"),
			EnvironmentId:          envres.EnvironmentId,
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Deployment"),
			},
		})

		if errDeploy != nil {
			return errDeploy
		}
		return nil
	})
}
