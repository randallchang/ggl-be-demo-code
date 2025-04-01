package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/randallchang/ggl-be-demo-code/internal/api"
	"github.com/randallchang/ggl-be-demo-code/internal/service"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			gin.New,
			service.NewTaskService,
			fx.Annotate(
				service.NewTaskService,
				fx.As(new(service.Service)),
			),
			api.NewHandler,
		),
		fx.Invoke(registerHooks),
	)

	app.Run()
}

func registerHooks(lifecycle fx.Lifecycle, router *gin.Engine, handler *api.Handler) {
	api.SetupRoutes(router, handler)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := router.Run(":8080"); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			return nil
		},
	})
}
