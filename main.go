package main

import (
	"backend/lr"
	"encoding/json"

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
	lr.Test()
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func createOrg(ctx *atreugo.RequestCtx) error {

	// Get the request body
	var permission = struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := json.Unmarshal(ctx.PostBody(), &permission); err != nil {
		return ctx.ErrorResponse(err, 400)
	}

	// Handle the request
	data, err := lr.Login(permission.Email, permission.Password)
	if err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	return ctx.JSONResponse(data, 201)
}
