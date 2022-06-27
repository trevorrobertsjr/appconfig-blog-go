package main

import (
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
		envres, err = appconfig.NewEnvironment(ctx, "exampleEnvironment", &appconfig.EnvironmentArgs{
			Description:   pulumi.String("Example AppConfig Environment"),
			ApplicationId: appres.ID(),
			Monitors: appconfig.EnvironmentMonitorArray{
				&appconfig.EnvironmentMonitorArgs{
					AlarmArn:     pulumi.Any(aws_cloudwatch_metric_alarm.Example.Arn),
					AlarmRoleArn: pulumi.Any(aws_iam_role.Example.Arn),
				},
			},
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Environment"),
			},
		})
		if err != nil {
			return err
		}
		deployres, err := appconfig.NewDeployment(ctx, "example", &appconfig.DeploymentArgs{
			ApplicationId:          appres.ID(),
			ConfigurationProfileId: pulumi.Any(aws_appconfig_configuration_profile.Example.Configuration_profile_id),
			ConfigurationVersion:   pulumi.Any(aws_appconfig_hosted_configuration_version.Example.Version_number),
			DeploymentStrategyId:   depstratres.ID(),
			Description:            pulumi.String("My example deployment"),
			EnvironmentId:          pulumi.Any(aws_appconfig_environment.Example.Environment_id),
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Deployment"),
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
