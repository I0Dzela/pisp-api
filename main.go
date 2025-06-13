package main

import (
	"github.com/I0Dzela/pisp-api/cmd"
	"github.com/I0Dzela/pisp-api/logger"
	cl "github.com/I0Dzela/pisp-common/logger"
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=common.config.yaml openapi/common.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=facekit.config.yaml openapi/facekit.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=file.config.yaml openapi/file.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=relation.config.yaml openapi/relation.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=user.config.yaml openapi/user.yaml

func main() {
	logX := logger.NewLogger(logrus.WithField(cl.EventField("main")))

	app := &cli.App{
		Commands: []*cli.Command{cmd.Server},
	}

	if err := app.Run(os.Args); errors.Wrap(err, 0) != nil {
		logX.Fatalf(errors.Wrap(err, 0), "run args: %v", os.Args)
	}
}
