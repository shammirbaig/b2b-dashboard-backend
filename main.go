package main

import (
	"backend/lr"
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/savsgio/atreugo/v11"
)

func main() {

	server := atreugo.New(atreugo.Config{
		Addr:             serverURL,
		GracefulShutdown: true,
	})
	server.UseBefore()

	authCtx := server.NewGroupPath("/auth")

	godotenv.Load(".env")

	lr.NewMongoClient()

	//lr.TestGetAllInvitationsOfOrganization()

	authCtx.POST("/login", login)
	authCtx.POST("/org/{id}/create", createOrg)
	authCtx.POST("/org/{id}/invite", inviteUser)
	authCtx.GET("/org/{id}/invitations", getAllInvitationsOfOrganization)
	authCtx.GET("/org/{id}/users", getAllUsersOfAnOrganization)
	authCtx.GET("/org/{id}/roles", getAllRolesOfAnOrg)
	authCtx.GET("/orgs", getAllOrganizationsOfTenant)

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
		Name    string `json:"name"`
		OldName string `json:"oldname"`
	}{}
	if err := json.Unmarshal(ctx.PostBody(), &org); err != nil {
		return ctx.ErrorResponse(err, 400)
	}

	tenantOrgID := ctx.UserValue("id").(string)

	// Handle the request
	if err := lr.CreateOrg(tenantOrgID, org.Name, org.OldName); err != nil {
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

func getAllInvitationsOfOrganization(ctx *atreugo.RequestCtx) error {

	orgId := ctx.UserValue("id").(string)

	// Handle the request
	data, err := lr.GetAllInvitationsOfOrganization(orgId)
	if err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	res := struct {
		Data []lr.InvitationResponse `json:"Data"`
	}{
		Data: data,
	}

	return ctx.JSONResponse(res, 200)
}

func getAllUsersOfAnOrganization(ctx *atreugo.RequestCtx) error {

	orgId := ctx.UserValue("id").(string)

	// Handle the request
	data, err := lr.GetAllUsersOfAnOrganization(orgId)
	if err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	res := struct {
		Data []lr.UserRole `json:"Data"`
	}{
		Data: data,
	}

	return ctx.JSONResponse(res, 200)
}

func getAllRolesOfAnOrg(ctx *atreugo.RequestCtx) error {

	orgId := ctx.UserValue("id").(string)

	// Handle the request
	data, err := lr.GetAllRolesOfAnOrg(orgId)
	if err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	res := struct {
		Data []lr.RoleResponse `json:"Data"`
	}{
		Data: data,
	}
	return ctx.JSONResponse(res, 200)
}

func getAllOrganizationsOfTenant(ctx *atreugo.RequestCtx) error {

	// Handle the request
	data, err := lr.GetAllOrganizationsOfTenant()
	if err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	res := struct {
		Data []lr.AllOrganizationsResponse `json:"Data"`
	}{
		Data: data,
	}

	return ctx.JSONResponse(res, 200)
}
