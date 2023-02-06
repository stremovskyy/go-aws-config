package main

import (
	"fmt"
	"github.com/karmadon/go-aws-config"
)

type Configuration struct {
	// DB configuration
	DBHost string `yaml:"db_host" json:"db_host"`
	DBPort string `yaml:"db_port" json:"db_port"`
}

func main() {
	// Create new options for client
	opts := &go_aws_config.Options{
		Region:           "eu-central-1",
		ApplicationID:    "test",
		EnvironmentID:    "TestEnvironment",
		Profile:          "test-config",
		PollingInterval:  60,
		CredentialsInEnv: true,
	}

	// Create new client
	configurator := go_aws_config.NewClient(opts)

	// Prepare client
	err := configurator.Prepare()
	if err != nil {
		fmt.Printf("failed to prepare client, %v", err)
		return
	}

	// create a struct to load configuration into
	dbConfiguration := &Configuration{}

	// load configuration into struct
	err = configurator.LoadIntoYaml(dbConfiguration)
	if err != nil {
		fmt.Printf("failed to load configuration into struct, %v", err)
		return
	}

	// PROFIT!
	fmt.Printf("db host: %s", dbConfiguration.DBHost)
}
