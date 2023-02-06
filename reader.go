package go_aws_config

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfigdata"
	"gopkg.in/yaml.v2"
)

type client struct {
	options       *Options
	appConfigData *appconfigdata.AppConfigData
	configToken   *string
}

// NewClient creates a new AWSConfigurator
func NewClient(options *Options) AWSConfigurator {
	if options == nil {
		options = &Options{
			Region:           "eu-central-1",
			ApplicationID:    "test",
			EnvironmentID:    "TestEnvironment",
			Profile:          "test-config",
			PollingInterval:  60,
			CredentialsInEnv: true,
		}
	}

	if options.Region == "" {
		options.Region = "eu-central-1"
	}

	return &client{options: options}
}

// Prepare prepares the client to load the configuration
func (c *client) Prepare() error {
	awsCfg := &aws.Config{
		Region: aws.String("eu-central-1"),
	}

	if c.options.CredentialsInEnv {
		awsCfg.Credentials = credentials.NewEnvCredentials()
	} else if c.options.AccessKeyID != "" && c.options.SecretAccessKey != "" {
		awsCfg.Credentials = credentials.NewStaticCredentials(c.options.AccessKeyID, c.options.SecretAccessKey, c.options.Token)
	}

	ses, err := session.NewSession(awsCfg)
	if err != nil {
		return fmt.Errorf("failed to create AWS session, %v", err)
	}

	c.appConfigData = appconfigdata.New(ses)

	sessInp := appconfigdata.StartConfigurationSessionInput{
		ApplicationIdentifier:                aws.String(c.options.ApplicationID),
		ConfigurationProfileIdentifier:       aws.String(c.options.Profile),
		EnvironmentIdentifier:                aws.String(c.options.EnvironmentID),
		RequiredMinimumPollIntervalInSeconds: aws.Int64(c.options.PollingInterval),
	}

	configSession, err := c.appConfigData.StartConfigurationSession(&sessInp)
	if err != nil {
		return fmt.Errorf("failed to start AWS appconfig session, %v", err)
	}

	c.configToken = configSession.InitialConfigurationToken

	return nil
}

// LoadConfigBytes loads the configuration as bytes
func (c *client) LoadConfigBytes() ([]byte, error) {
	getLatestConfig := &appconfigdata.GetLatestConfigurationInput{
		ConfigurationToken: c.configToken,
	}

	latestConfiguration, err := c.appConfigData.GetLatestConfiguration(getLatestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest AWS appconfig configuration, %v", err)
	}

	return latestConfiguration.Configuration, nil
}

// LoadIntoYaml loads the configuration (in YAML format) into the given struct
func (c *client) LoadIntoYaml(dest interface{}) error {
	config, err := c.LoadConfigBytes()
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(config, dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal AWS appconfig configuration, %v", err)
	}

	return nil
}

// LoadIntoJson loads the configuration (in JSON format) into the given struct
func (c *client) LoadIntoJson(dest interface{}) error {
	config, err := c.LoadConfigBytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal AWS appconfig configuration, %v", err)
	}

	return nil
}
