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

	// lr.Test()

	authCtx.POST("/login", login)
	authCtx.POST("/org/{id}/create", createOrg)
	authCtx.POST("/org/{id}/invite", inviteUser)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func login(ctx *atreugo.RequestCtx) error {

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

func createOrg(ctx *atreugo.RequestCtx) error {

	// Get the request body
	var org = struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(ctx.PostBody(), &org); err != nil {
		return ctx.ErrorResponse(err, 400)
	}

	tenantOrgID := ctx.UserValue("id").(string)

	// Handle the request
	if err := lr.CreateOrg(tenantOrgID, org.Name); err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	return ctx.JSONResponse(nil, 201)
}

func inviteUser(ctx *atreugo.RequestCtx) error {

	// Get the request body
	var invite lr.SendInvitation
	if err := json.Unmarshal(ctx.PostBody(), &invite); err != nil {
		return ctx.ErrorResponse(err, 400)
	}

	orgId := ctx.UserValue("id").(string)

	// Handle the request
	if err := lr.InviteUser(orgId, invite); err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	return ctx.JSONResponse(nil, 201)
}
