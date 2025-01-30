package lr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Login returns the orgs that the user is a part of
func Login(ctx context.Context, email, password string) (*AuthResponse, []OrganizationResponse, error) {

	data, err := Post(ctx, loginUrl, strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	var resp AuthResponse
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&resp); err != nil {
		return nil, nil, err
	}

	orgs, err := GetUserOrgs(ctx, *resp.Profile.Uid)
	if err != nil {
		return nil, nil, err
	}

	fmt.Println(orgs)
	return &resp, orgs, nil
}

func GetUserOrgs(ctx context.Context, uid string) ([]OrganizationResponse, error) {
	data, err := Get(ctx, getOrgsOfUserUrl(uid), "")
	if err != nil {
		return nil, err
	}

	var orgsResp struct {
		Data []OrganizationResponse `json:"Data"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&orgsResp); err != nil {
		return nil, err
	}

	return orgsResp.Data, nil
}

func GetAllUsersOfAnOrganization(ctx context.Context, orgID string) ([]UserRole, error) {
	data, err := Get(ctx, getUsersOfOrgUrl(orgID), "")
	if err != nil {
		return nil, err
	}

	var usersResp struct {
		Data []UserRole `json:"Data"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&usersResp); err != nil {
		return nil, err
	}

	return usersResp.Data, nil
}

func CreateOrg(ctx context.Context, tenantOrgID string, newOrgName string, oldOrgName string) error {

	// create an app : 7c9254057e2044c5b3fadf8bf0b3dd31

	// AppName should be current org name and mark its name as appname
	payload := `
		"AppName": "{oldOrgName}", 
		"Domain": "login2website.com; localhost:3000",
		"CallbackUrl": "login2website.com; localhost",
		"DevDomain": "login2website.com"
	}`

	payload = strings.Replace(payload, "{oldOrgName}", oldOrgName, 1)

	res, err := Post(ctx, createAppUrl(), strings.NewReader(payload))
	if err != nil {
		return err
	}

	var appData AppResponse
	if err := json.NewDecoder(bytes.NewReader(res)).Decode(&appData); err != nil {
		return err
	}

	feature := strings.NewReader(`
	{
		"Data": [
			{
				"feature": "enable_b2b_identity",
				"status": true
			}
		]
	}`)

	// turn on b2b feature
	_, err = Put(ctx, turnB2BApp(appData.OwnerId, strconv.Itoa(appData.AppId)), feature)
	if err != nil {
		return err
	}

	// create a relation between orgid and appid
	if err := CreateAppidToOrgidMapping(mongoClient, appData.AppId, tenantOrgID); err != nil {
		return err
	}

	// create an org for the above passed data
	orgPayload := `{
		"Name": "{orgName}",
		"Domains": [
			{
				"DomainName": "loginradius.com",
				"IsVerified": true
			}
		],
		"IsAuthRestrictedToDomain": false
	}`

	orgPayload = strings.Replace(orgPayload, "{orgName}", newOrgName, 1)

	// this creates in nike-com db but how do we differentiate the hirearcy
	orgRes, err := DynamicPost(ctx, strconv.Itoa(appData.AppId), appData.OwnerId, createOrgUrl(), strings.NewReader(orgPayload))
	if err != nil {
		return err
	}

	var orgDataResp Organizations
	if err := json.NewDecoder(bytes.NewReader(orgRes)).Decode(&orgDataResp); err != nil {
		return err
	}

	return nil
}

func InviteUser(ctx context.Context, orgId string, invite SendInvitation) error {

	payload, err := json.Marshal(invite)
	if err != nil {
		return err
	}

	_, err = Post(ctx, sendInvitationUrl(orgId), bytes.NewReader(payload))
	if err != nil {
		return err
	}

	return nil
}

func GetAllRolesOfAnOrg(ctx context.Context, orgId string) ([]RoleResponse, error) {
	data, err := Get(ctx, getAllRolesOfAnOrg(orgId), "")
	if err != nil {
		return nil, err
	}

	var rolesResp struct {
		Data []RoleResponse `json:"Data"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&rolesResp); err != nil {
		return nil, err
	}

	return rolesResp.Data, nil
}

func GetAllOrganizationsOfTenant(ctx context.Context, orgId string) ([]AllOrganizationsResponse, error) {

	//get AppId from the orgId-appId relation
	if orgId != "" {
		appId, err := GetAppIdFromOrgIdMapping(mongoClient, orgId)
		if err != nil {
			return nil, err
		}

		data, err := Get(ctx, getOrgsOfTenantUrl(), strconv.Itoa(appId))
		if err != nil {
			return nil, err
		}

		var orgsResp struct {
			Data []AllOrganizationsResponse `json:"Data"`
		}

		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&orgsResp); err != nil {
			return nil, err
		}

		return orgsResp.Data, nil

	}

	data, err := Get(ctx, getOrgsOfTenantUrl(), "")
	if err != nil {
		return nil, err
	}

	var orgsResp struct {
		Data []AllOrganizationsResponse `json:"Data"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&orgsResp); err != nil {
		return nil, err
	}

	return orgsResp.Data, nil
}

func GetAllInvitationsOfOrganization(ctx context.Context, orgID string) ([]InvitationResponse, error) {
	data, err := Get(ctx, getAllInvitationsOfOrganization(orgID), "")
	if err != nil {
		return nil, err
	}

	var invitationsResp struct {
		Data       []InvitationResponse `json:"Data"`
		TotalCount int                  `json:"TotalCount"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&invitationsResp); err != nil {
		return nil, err
	}

	return invitationsResp.Data, nil
}

func GetAllRolesOfUserInOrg(orgID, uid string) ([]RoleResponse, error) {

	data, err := Get(context.Background(), getAllRolesOfUserInOrg(orgID, uid), "")
	if err != nil {
		return nil, err
	}

	var rolesResp struct {
		Data []RoleResponse `json:"Data"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&rolesResp); err != nil {
		return nil, err
	}

	return rolesResp.Data, nil
}

func Test() {
	if err := CreateOrg(context.Background(), "qwer", "niketest123", ""); err != nil {
		fmt.Println("Error:", err)
	}
}

func TestLogin() {
	// orgs, _ := Login("nike@mailinator.com", "123456")
	// fmt.Println("Response:", orgs)
}

func TestGetAllUsersOfAnOrganization() {
	users, _ := GetAllUsersOfAnOrganization(context.Background(), "org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", users)
}

func TestGetAllRolesOfAnOrg() {
	roles, _ := GetAllRolesOfAnOrg(context.Background(), "org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", roles)
}

func TestGetAllOrganizationsOfTenant() {
	orgs, _ := GetAllOrganizationsOfTenant(context.Background(), "org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", orgs)
}

func TestGetAllInvitationsOfOrganization() {
	invitations, _ := GetAllInvitationsOfOrganization(context.Background(), "org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", invitations)
}

func TestGetAllRolesOfUserInOrg() {
	roles, _ := GetAllRolesOfUserInOrg("org_Z5pkOZ-0eGkkbhQ1", "auth_0eGkkbhQ1")
	fmt.Println("Response:", roles)
}
