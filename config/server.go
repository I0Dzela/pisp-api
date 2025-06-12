package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/urfave/cli/v2"
)

type ServerConfig struct {
	Port        uint64 `validate:"required,number,gte=0,lte=65535"`
	Scheme      string `validate:"required,oneof=http https"`
	Environment string `validate:"required,oneof=development production"`
	LogLevel    uint64 `validate:"required,min=0,max=6"`
	SwaggerHost string `validate:"required"`
}

var Server = new(ServerConfig)

func (c *ServerConfig) Setup() []cli.Flag {
	return []cli.Flag{
		&cli.Uint64Flag{
			Name:        "port",
			Aliases:     []string{"p"},
			Value:       1234,
			Usage:       "server port",
			Destination: &c.Port,
			EnvVars:     []string{"PORT"},
		},
		&cli.StringFlag{
			Name:        "scheme",
			Aliases:     []string{"s"},
			Value:       "https",
			Usage:       "server scheme",
			Destination: &c.Scheme,
			EnvVars:     []string{"SCHEME"},
		},
		&cli.StringFlag{
			Name:        "env",
			Value:       "development",
			Usage:       "environment",
			Destination: &c.Environment,
			EnvVars:     []string{"ENV"},
		},
		&cli.Uint64Flag{
			Name:        "logLevel",
			Value:       6,
			Usage:       "jwt secret",
			Destination: &c.LogLevel,
			EnvVars:     []string{"LOG_LEVEL"},
		},
		&cli.StringFlag{
			Name:        "swagger-host",
			Value:       "https://pisp.local:6030",
			Usage:       "swagger host",
			Destination: &c.SwaggerHost,
			EnvVars:     []string{"SWAGGER_HOST"},
		},
	}
}

func (c *ServerConfig) Validate() error {
	validate := validator.New()

	return validate.Struct(c)
}
