package main

import (
	"github.com/savsgio/atreugo/v11"
)

func main() {

	server := atreugo.New(atreugo.Config{
		Addr:             serverURL,
		GracefulShutdown: true,
	})
	server.UseBefore()

	authCtx := server.NewGroupPath("/auth")

	authCtx.POST("/org", createOrg)
}

func createOrg(ctx *atreugo.RequestCtx) error {

	// Handle the request
	return ctx.JSONResponse("Org created successfully", 201)
}
