package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/appconfig"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create an AWS AppConfig resource for our application's feature flags
		applicationResource, err := appconfig.NewApplication(ctx, "blogApplication", &appconfig.ApplicationArgs{
			Description: pulumi.String("Blog AppConfig Application"),
			Name:        pulumi.String("blogAppConfigGo"),
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Application"),
			},
		})
		if err != nil {
			return err
		}

		// Specify the update frequency for changes to the application's feature flags
		deploymentStrategyResource, err := appconfig.NewDeploymentStrategy(ctx, "blogDeploymentStrategy", &appconfig.DeploymentStrategyArgs{
			DeploymentDurationInMinutes: pulumi.Int(1),
			Description:                 pulumi.String("Blog Deployment Strategy"),
			FinalBakeTimeInMinutes:      pulumi.Int(1),
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

		// Define which application environment (ex: dev, prod, alpha, beta)
		environmentResource, err := appconfig.NewEnvironment(ctx, "blogEnvironment", &appconfig.EnvironmentArgs{
			Description:   pulumi.String("Blog AppConfig Environment"),
			ApplicationId: applicationResource.ID(),
			Name:          pulumi.String("prod"),
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Environment"),
			},
		})
		if err != nil {
			return err
		}

		// Specify the type of resource our application will consume for the service
		// We are using feature flags in this application, but free form JSON input
		// is also supported.
		configurationProfileResource, err := appconfig.NewConfigurationProfile(ctx, "blogConfigurationProfile", &appconfig.ConfigurationProfileArgs{
			ApplicationId: applicationResource.ID(),
			Name:          pulumi.String("whichSide"),
			Description:   pulumi.String("Blog Configuration Profile"),
			LocationUri:   pulumi.String("hosted"),
			Type:          pulumi.String("AWS.AppConfig.FeatureFlags"),
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Configuration Profile"),
			},
		})
		if err != nil {
			return err
		}
		acceptableValues := make([]string, 2)
		acceptableValues[0] = "darkknight"
		acceptableValues[1] = "paladin"

		// Define the structure of the feature flag(s) our application will use
		rawFeatureFlagJSONInput, err := json.Marshal(map[string]interface{}{
			"flags": map[string]interface{}{
				"allegiance": map[string]interface{}{
					"name": "allegiance",
					"attributes": map[string]interface{}{
						"choice": map[string]interface{}{
							"constraints": map[string]interface{}{
								"type":     "string",
								"enum":     acceptableValues,
								"required": true,
							},
						},
					},
				},
			},
			"values": map[string]interface{}{
				"allegiance": map[string]interface{}{
					"enabled": "true",
					"choice":  "paladin",
				},
			},
			"version": "1",
		})
		if err != nil {
			return err
		}
		// Convert the feature flag input to a string for input to the
		// Hosted Configuration Version Resource
		featureFlagJSONInput := string(rawFeatureFlagJSONInput)

		// Create a configuration version for each change to the feature flag(s)
		hostedConfigurationVersionResource, err := appconfig.NewHostedConfigurationVersion(ctx, "blogHostedConfigurationVersion", &appconfig.HostedConfigurationVersionArgs{
			ApplicationId:          applicationResource.ID(),
			ConfigurationProfileId: configurationProfileResource.ConfigurationProfileId,
			Description:            pulumi.String("Blog Feature Flag Hosted Configuration Version"),
			ContentType:            pulumi.String("application/json"),
			Content:                pulumi.String(featureFlagJSONInput),
		})
		if err != nil {
			return err
		}

		// Deploy the feature flag.
		_, errDeploy := appconfig.NewDeployment(ctx, "blogDeployment", &appconfig.DeploymentArgs{
			ApplicationId:          applicationResource.ID(),
			ConfigurationProfileId: configurationProfileResource.ConfigurationProfileId,
			ConfigurationVersion:   pulumi.Sprintf("%v", hostedConfigurationVersionResource.VersionNumber),
			DeploymentStrategyId:   deploymentStrategyResource.ID(),
			Description:            pulumi.String("My Blog Deployment"),
			EnvironmentId:          environmentResource.EnvironmentId,
			Tags: pulumi.StringMap{
				"Type": pulumi.String("AppConfig Deployment"),
			},
		})

		if errDeploy != nil {
			return errDeploy
		}

		// Assume Role Policy for ecs tasks

		tmpJSON0, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				map[string]interface{}{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Sid":    "",
					"Principal": map[string]interface{}{
						"Service": "ecs-tasks.amazonaws.com",
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Deploy an IAM role with appconfig permisison
		json0 := string(tmpJSON0)
		// ecsTaskRole, err := iam.NewRole(ctx, "appConfigRole", &iam.RoleArgs{
		_, err = iam.NewRole(ctx, "appConfigRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(json0),
			ManagedPolicyArns: pulumi.StringArray{pulumi.String("arn:aws:iam::318168271290:policy/AppConfigAllAccess")},
			Tags: pulumi.StringMap{
				"tag-key": pulumi.String("tag-value"),
			},
		})
		if err != nil {
			return err
		}

		// // Deploy an ECS Cluster
		// ecsCluster, err := ecs.NewCluster(ctx, "appConfigAgentFargate", &ecs.ClusterArgs{

		// })
		// if err != nil {
		// 	return err
		// }
		// // Deploy an ECS Task Definition

		// // Task Definition Properties
		// ecsTaskDefTmpJSON0, err := json.Marshal([]interface{}{
		// 	map[string]interface{}{
		// 		"name":      "appConfigAgent",
		// 		"image":     "public.ecr.aws/aws-appconfig/aws-appconfig-agent:2.x",
		// 		"essential": true,
		// 		"portMappings": []map[string]interface{}{
		// 			map[string]interface{}{
		// 				"containerPort": 2772,
		// 				"hostPort":      2772,
		// 			},
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }
		// // Create Task Definition from properties described above
		// ecsTaskDefJSON0 := string(ecsTaskDefTmpJSON0)
		// ecsTaskDefinition, err := ecs.NewTaskDefinition(ctx, "AppConfigAgent", &ecs.TaskDefinitionArgs{
		// 	Family:               pulumi.String("AppConfigAgent"),
		// 	ContainerDefinitions: pulumi.String(ecsTaskDefJSON0),
		// 	TaskRoleArn: ecsTaskRole.Arn,
		// 	Cpu:                  pulumi.String("256"),
		// 	Memory:               pulumi.String("512"),
		// 	NetworkMode:          pulumi.String("awsvpc"),
		// 	RequiresCompatibilities: pulumi.StringArray{
		// 		pulumi.String("FARGATE"),
		// 	},
		// })
		// if err != nil {
		// 	return err
		// }

		// // Deploy a Fargate service to ECS
		// _, err = ecs.NewService(ctx, "AppConfigListener", &ecs.ServiceArgs{
		// 	Cluster:        ecsCluster.Arn,
		// 	TaskDefinition: ecsTaskDefinition.Arn,
		// 	DesiredCount:   pulumi.Int(1),
		// 	LaunchType: pulumi.String("FARGATE"),
		// 	NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
		// 			Subnets: pulumi.StringArray{pulumi.String("subnet-09dd613a058ca077c"),},
		// 			AssignPublicIp: pulumi.Bool(true),
		// 			SecurityGroups: pulumi.StringArray{pulumi.String("sg-0198dc928c67a40b3"),},
		// 		},

		// })

		// if err != nil {
		// 	return err
		// }

		// // Sample curl
		// // $ curl "http://44.212.32.110:2772/applications/blogAppConfigGo/environments/prod/configurations/whichSide?flag=allegiance"
		return nil
	})
}
