package go_aws_config

type AWSConfigurator interface {
	Prepare() error
	LoadConfigBytes() ([]byte, error)
	LoadIntoYaml(dest interface{}) error
	LoadIntoJson(dest interface{}) error
}
