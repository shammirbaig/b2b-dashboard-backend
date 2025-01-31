package main

import (
	"backend/lr"
	"encoding/json"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/savsgio/atreugo/v11"
)

func main() {

	server := atreugo.New(atreugo.Config{
		Addr:             serverURL,
		GracefulShutdown: true,
	})
	godotenv.Load(".env")

	server.UseBefore(func(ctx *atreugo.RequestCtx) error {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "*")
		ctx.Response.Header.Set("Content-Type", "application/json;charset=utf-8")

		return ctx.Next()
	})

	authCtx := server.NewGroupPath("/auth")

	lr.NewMongoClient()

	authCtx.GET("/test", func(ctx *atreugo.RequestCtx) error {
		return ctx.TextResponse("Hello, World!")
	})
	authCtx.POST("/login", login)
	authCtx.POST("/org/{id}/create", createOrg)
	authCtx.POST("/org/{id}/invitations", inviteUser)
	authCtx.GET("/org/{id}/invitations", getAllInvitationsOfOrganization)
	authCtx.GET("/org/{id}/users", getAllUsersOfAnOrganization)
	authCtx.GET("/org/{id}/roles", getAllRolesOfAnOrg)
	authCtx.GET("/orgs", getAllOrganizationsOfTenant)
	authCtx.GET("/org/{orgId}/user/{userId}/roles", getAllRolesOfUserInOrg)

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
		return ctx.ErrorResponse(err, http.StatusBadRequest)
	}

	// Handle the request
	uid, data, err := lr.Login(permission.Email, permission.Password)
	if err != nil {
		return ctx.JSONResponse(err.Error(), http.StatusForbidden)
	}

	resp := struct {
		UserId            string                    `json:"userId"`
		OrganizationsList []lr.OrganizationResponse `json:"organizationsList"`
		Token             string                    `json:"token"`
	}{
		UserId:            uid,
		OrganizationsList: data,
		Token:             "",
	}

	return ctx.JSONResponse(resp, http.StatusCreated)
}

func createOrg(ctx *atreugo.RequestCtx) error {

	// Get the request body
	var org = struct {
		Uid     string `json:"uid"`
		Name    string `json:"name"`
		OldName string `json:"oldname"`
	}{}
	if err := json.Unmarshal(ctx.PostBody(), &org); err != nil {
		return ctx.ErrorResponse(err, http.StatusBadRequest)
	}

	tenantOrgID := ctx.UserValue("id").(string)

	// Handle the request
	if err := lr.CreateOrg(org.Uid, tenantOrgID, org.Name, org.OldName); err != nil {
		return ctx.ErrorResponse(err, http.StatusForbidden)
	}

	return ctx.JSONResponse(nil, 201)
}

func inviteUser(ctx *atreugo.RequestCtx) error {

	// Get the request body
	var invite lr.SendInvitation
	if err := json.Unmarshal(ctx.PostBody(), &invite); err != nil {
		return ctx.ErrorResponse(err, http.StatusBadRequest)
	}

	orgId := ctx.UserValue("id").(string)

	// Handle the request
	if err := lr.InviteUser(orgId, invite); err != nil {
		return ctx.ErrorResponse(err, http.StatusForbidden)
	}

	return ctx.JSONResponse(nil, http.StatusCreated)
}

func getAllInvitationsOfOrganization(ctx *atreugo.RequestCtx) error {

	orgId := ctx.UserValue("id").(string)

	// Handle the request
	data, err := lr.GetAllInvitationsOfOrganization(orgId)
	if err != nil {
		return ctx.ErrorResponse(err, http.StatusForbidden)
	}

	return ctx.JSONResponse(*data, http.StatusOK)
}

func getAllUsersOfAnOrganization(ctx *atreugo.RequestCtx) error {

	orgId := ctx.UserValue("id").(string)

	// Handle the request
	data, err := lr.GetAllUsersOfAnOrganization(orgId)
	if err != nil {
		return ctx.ErrorResponse(err, http.StatusForbidden)
	}

	res := struct {
		Data []lr.UserRole `json:"Data"`
	}{
		Data: data,
	}

	return ctx.JSONResponse(res, http.StatusOK)
}

func getAllRolesOfAnOrg(ctx *atreugo.RequestCtx) error {

	orgId := ctx.UserValue("id").(string)

	// Handle the request
	data, err := lr.GetAllRolesOfAnOrg(orgId)
	if err != nil {
		return ctx.ErrorResponse(err, http.StatusForbidden)
	}

	res := struct {
		Data []lr.RoleResponse `json:"Data"`
	}{
		Data: data,
	}
	return ctx.JSONResponse(res, http.StatusOK)
}

func getAllOrganizationsOfTenant(ctx *atreugo.RequestCtx) error {

	orgId := string(ctx.QueryArgs().Peek("orgId"))
	// Handle the request
	data, err := lr.GetAllOrganizationsOfTenant(orgId)
	if err != nil {
		return ctx.ErrorResponse(err, http.StatusForbidden)
	}

	res := struct {
		Data []lr.AllOrganizationsResponse `json:"Data"`
	}{
		Data: data,
	}

	return ctx.JSONResponse(res, http.StatusOK)
}

func getAllRolesOfUserInOrg(ctx *atreugo.RequestCtx) error {

	orgId := ctx.UserValue("orgId").(string)
	userId := ctx.UserValue("userId").(string)

	// Handle the request
	data, err := lr.GetAllRolesOfUserInOrg(orgId, userId)
	if err != nil {
		return ctx.ErrorResponse(err, http.StatusForbidden)
	}

	res := struct {
		Data []lr.RoleResponse `json:"Data"`
	}{
		Data: data,
	}
	return ctx.JSONResponse(res, http.StatusOK)
}
