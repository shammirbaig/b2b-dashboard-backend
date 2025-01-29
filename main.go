package main

import (
	"backend/lr"

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

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func createOrg(ctx *atreugo.RequestCtx) error {

	// Handle the request
	lr.Test()

	return ctx.JSONResponse("Org created successfully", 201)
}
