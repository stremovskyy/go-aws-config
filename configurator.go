package go_aws_config

import (
	"context"
	"io"
)

type AWSConfigurator interface {
	Prepare(ctx context.Context) error
	LoadConfigBytes() ([]byte, error)
	LoadIntoYaml(dest interface{}) error
	LoadIntoJson(dest interface{}) error
	Reader() (io.Reader, error)
}
