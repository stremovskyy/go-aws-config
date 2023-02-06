package go_aws_config

import "io"

type AWSConfigurator interface {
	Prepare() error
	LoadConfigBytes() ([]byte, error)
	LoadIntoYaml(dest interface{}) error
	LoadIntoJson(dest interface{}) error
	Reader() (io.Reader, error)
}
