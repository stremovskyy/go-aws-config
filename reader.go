package go_aws_config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfigdata"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
)

const (
	defaultRegion           = "eu-central-1"
	defaultPollingInterval  = 60
	defaultCredentialsInEnv = true
)

type client struct {
	options       *Options
	appConfigData *appconfigdata.AppConfigData
	configToken   *string
}

// NewClient creates a new AWSConfigurator
func NewClient(options *Options) (AWSConfigurator, error) {
	if options == nil {
		options = &Options{
			CredentialsInEnv: defaultCredentialsInEnv,
		}
	}

	// Set default values for missing options
	if options.Region == "" {
		options.Region = defaultRegion
	}
	if options.PollingInterval == 0 {
		options.PollingInterval = defaultPollingInterval
	}

	return &client{options: options}, nil
}

// Prepare prepares the client to load the configuration from AWS AppConfig
func (c *client) Prepare(ctx context.Context) error {
	awsCfg := &aws.Config{
		Region: aws.String(c.options.Region),
	}

	if c.options.CredentialsInEnv {
		awsCfg.Credentials = credentials.NewEnvCredentials()
	} else if c.options.AccessKeyID != "" && c.options.SecretAccessKey != "" {
		awsCfg.Credentials = credentials.NewStaticCredentials(c.options.AccessKeyID, c.options.SecretAccessKey, c.options.Token)
	}

	sess, err := session.NewSession(awsCfg)
	if err != nil {
		return fmt.Errorf("failed to create AWS session, %v", err)
	}

	c.appConfigData = appconfigdata.New(sess)

	sessInp := appconfigdata.StartConfigurationSessionInput{
		ApplicationIdentifier:                aws.String(c.options.ApplicationID),
		ConfigurationProfileIdentifier:       aws.String(c.options.Profile),
		EnvironmentIdentifier:                aws.String(c.options.EnvironmentID),
		RequiredMinimumPollIntervalInSeconds: aws.Int64(c.options.PollingInterval),
	}

	// Start the configuration session
	configSession, err := c.appConfigData.StartConfigurationSessionWithContext(ctx, &sessInp)
	if err != nil {
		return fmt.Errorf("failed to start AWS AppConfig session, %v", err)
	}

	// Set the configuration token on the client
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

func (c *client) Reader() (io.Reader, error) {
	config, err := c.LoadConfigBytes()
	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(config)), nil
}
