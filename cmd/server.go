package cmd

import (
	"fmt"
	cl "github.com/I0Dzela/pisp-common/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urfave/cli/v2"

	"github.com/I0Dzela/pisp-api/config"
	"github.com/I0Dzela/pisp-api/logger"
)

const service = "api"

var Server = &cli.Command{
	Name:        "server",
	Description: fmt.Sprintf("start pisp %s server", service),
	Flags:       config.Server.Setup(),
	Action: func(_ *cli.Context) error {
		logX := logger.NewLogger(nil)

		if err := config.Server.Validate(); errors.Wrap(err, 0) != nil {
			logX.Fatal(errors.Wrap(err, 0))
		}

		logEnv(logX)

		gin.SetMode(gin.ReleaseMode)
		//gr := gin.Default()
		gr := gin.New()
		gr.UseH2C = true
		//gr.RemoveExtraSlash = true
		gr.Use(gin.Recovery())

		gr.Use(cl.DefaultLogger(service))

		gr.StaticFile("/openapi/common.yaml", "/app/openapi/common.yaml")

		facekitSwaggerUrl := Init(gr, "facekit")
		logX.Infof("Specs served at: %s", facekitSwaggerUrl)

		fileSwaggerUrl := Init(gr, "file")
		logX.Infof("Specs served at: %s", fileSwaggerUrl)

		relationSwaggerUrl := Init(gr, "relation")
		logX.Infof("Specs served at: %s", relationSwaggerUrl)

		userSwaggerUrl := Init(gr, "user")
		logX.Infof("Specs served at: %s", userSwaggerUrl)

		addr := fmt.Sprintf(":%d", config.Server.Port)
		logX.Infof("rest server listening at %s", addr)
		if config.Server.Scheme != "https" {
			if err := gr.Run(addr); errors.Wrap(err, 0) != nil {
				logX.Fatal(errors.Wrap(err, 0))
			}
		} else {
			if err := gr.RunTLS(addr, "pisp.local.crt", "pisp.local.key"); errors.Wrap(err, 0) != nil {
				logX.Fatal(errors.Wrap(err, 0))
			}
		}

		return nil
	},
}

func logEnv(logX cl.LoggerX) {
	logX.Infof("environment: %s", config.Server.Environment)
	logX.Infof("log level: %d", config.Server.LogLevel)
}

func Init(gr *gin.Engine, serverName string) string {
	swaggerUrl := ginSwagger.URL(fmt.Sprintf("%s/openapi/%s.yaml", config.Server.SwaggerHost, serverName))
	swaggerInstanceName := ginSwagger.InstanceName(serverName)

	gr.GET(fmt.Sprintf("/%s/swagger/*any", serverName), ginSwagger.WrapHandler(swaggerFiles.NewHandler(), swaggerUrl, swaggerInstanceName))
	gr.StaticFile(fmt.Sprintf("openapi/%s.yaml", serverName), fmt.Sprintf("/app/openapi/%s.yaml", serverName))

	return fmt.Sprintf("%s/%s/swagger/index.html", config.Server.SwaggerHost, serverName)
}
