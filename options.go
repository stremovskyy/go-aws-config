package go_aws_config

type Options struct {
	// AWS region
	Region string
	// AWS AppConfig application ID
	ApplicationID string
	// AWS AppConfig environment ID
	EnvironmentID string
	// AWS AppConfig configuration profile
	Profile string
	// AWS AppConfig poling interval in seconds
	PollingInterval int64
	// AWS AppConfig credentials in environment variables
	CredentialsInEnv bool
	// AWS Access Key ID
	AccessKeyID string
	// AWS Secret Access Key
	SecretAccessKey string
	// AWS Token
	Token string
}
