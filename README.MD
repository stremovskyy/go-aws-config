# GO AWS Config

**go-aws-config** is a very simple wrapper around the AWS SDK for Go. It provides a simple way to get configuration values from AWS App Config.

### Installation

```bash
go get github.com/awslabs/go-config
```

### Usage

#### With environment variables
```go
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
```

#### With credentials in code
```go
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
	awsKey := "AKIAIOSFODNN7EXAMPLE"
	awsSecret := "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	
    // Create new options for client
    opts := &go_aws_config.Options{
        Region:           "eu-central-1",
        ApplicationID:    "test",
        EnvironmentID:    "TestEnvironment",
        Profile:          "test-config",
        PollingInterval:  60,
        CredentialsInEnv: false,
        AccessKeyID:      awsKey,
        SecretAccessKey:  awsSecret,
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
```




### License

This library is licensed under the Apache 2.0 License.

### Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
Please make sure to update tests as appropriate.

### Authors

* **Karmadon** - *Initial work* - [Karmadon](https://github.com/karmadon)

### Acknowledgments

* [AWS SDK for Go](https://github.com/aws/aws-sdk-go/)
* [AWS AppConfig](https://aws.amazon.com/appconfig/)
* [AWS AppConfig Go SDK](https://docs.aws.amazon.com/sdk-for-go/api/service/appconfig/)

### TODO

* [ ] Add support for more configuration formats
* [ ] Add support for subscribing to configuration changes
