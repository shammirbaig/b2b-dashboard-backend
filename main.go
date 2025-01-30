package main

import (
	"backend/lr"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/savsgio/atreugo/v11"
)

func main() {

	server := atreugo.New(atreugo.Config{
		Addr:             serverURL,
		GracefulShutdown: true,
	})

	server.UseBefore(func(ctx *atreugo.RequestCtx) error {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "*")
		ctx.Response.Header.Set("Content-Type", "application/json;charset=utf-8")

		if string(ctx.Path()) != "/auth/login" {
			tokenString := string(ctx.Request.Header.Peek("Authorization"))
			if tokenString == "" {
				return ctx.ErrorResponse(fmt.Errorf("missing token"), 401)
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})

			if err != nil || !token.Valid {
				return ctx.ErrorResponse(fmt.Errorf("invalid token"), 401)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return ctx.ErrorResponse(fmt.Errorf("invalid claims"), 401)
			}

			ctx.SetUserValue("uid", claims["uid"])

		}

		return ctx.Next()
	})

	authCtx := server.NewGroupPath("/auth")

	godotenv.Load(".env")

	lr.NewMongoClient()

	authCtx.GET("/test", func(ctx *atreugo.RequestCtx) error {
		return ctx.TextResponse("Hello, World!")
	})
	authCtx.POST("/login", login)

	// Authenticated routes
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
		return ctx.ErrorResponse(err, 400)
	}

	// Handle the request
	res, data, err := lr.Login(ctx, permission.Email, permission.Password)
	if err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	resp := struct {
		UserId            string                    `json:"userId"`
		OrganizationsList []lr.OrganizationResponse `json:"organizationsList"`
		Token             string                    `json:"token"`
	}{
		UserId:            *res.Profile.IdentityResponse.Identity.Uid,
		OrganizationsList: data,
		Token:             createSessionToken(*res.Profile.IdentityResponse.Identity.Uid),
	}

	return ctx.JSONResponse(resp, 201)
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
	if err := lr.CreateOrg(ctx, tenantOrgID, org.Name, org.OldName); err != nil {
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
	if err := lr.InviteUser(ctx, orgId, invite); err != nil {
		return ctx.ErrorResponse(err, 500)
	}

	return ctx.JSONResponse(nil, 201)
}

func getAllInvitationsOfOrganization(ctx *atreugo.RequestCtx) error {

	orgId := ctx.UserValue("id").(string)

	// Handle the request
	data, err := lr.GetAllInvitationsOfOrganization(ctx, orgId)
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
	data, err := lr.GetAllUsersOfAnOrganization(ctx, orgId)
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
	data, err := lr.GetAllRolesOfAnOrg(ctx, orgId)
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

	orgId := string(ctx.QueryArgs().Peek("orgId"))
	// Handle the request
	data, err := lr.GetAllOrganizationsOfTenant(ctx, orgId)
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

func getAllRolesOfUserInOrg(ctx *atreugo.RequestCtx) error {

	orgId := ctx.UserValue("orgId").(string)
	userId := ctx.UserValue("userId").(string)

	// Handle the request
	data, err := lr.GetAllRolesOfUserInOrg(orgId, userId)
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

func createSessionToken(userId string) string {
	// Define the token claims
	claims := jwt.MapClaims{
		"uid": userId,
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Token expiration time
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		log.Fatalf("Error signing token: %v", err)
	}

	return tokenString
}
