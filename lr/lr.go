package lr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Login returns the orgs that the user is a part of
func Login(email, password string) (string, []OrganizationResponse, error) {

	data, err := Post(loginUrl, strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	var tokenResp struct {
		Profile struct {
			Uid string `json:"Uid"`
		} `json:"Profile"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&tokenResp); err != nil {
		return "", nil, err
	}

	orgs, err := GetUserOrgs(tokenResp.Profile.Uid)
	if err != nil {
		return "", nil, err
	}

	fmt.Println(orgs)
	return tokenResp.Profile.Uid, orgs, nil
}

func GetUserOrgs(uid string) ([]OrganizationResponse, error) {
	data, err := Get(getOrgsOfUserUrl(uid), "")
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

func GetAllUsersOfAnOrganization(orgID string) ([]UserRole, error) {
	data, err := Get(getUsersOfOrgUrl(orgID), "")
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

func CreateOrg(tenantOrgID string, newOrgName string, oldOrgName string) error {

	// create an app : 7c9254057e2044c5b3fadf8bf0b3dd31

	// AppName should be current org name and mark its name as appname
	payload := `
		"AppName": "{oldOrgName}", 
		"Domain": "login2website.com; localhost:3000",
		"CallbackUrl": "login2website.com; localhost",
		"DevDomain": "login2website.com"
	}`

	payload = strings.Replace(payload, "{oldOrgName}", oldOrgName, 1)

	res, err := Post(createAppUrl(), strings.NewReader(payload))
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
	_, err = Put(turnB2BApp(appData.OwnerId, strconv.Itoa(appData.AppId)), feature)
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
	orgRes, err := DynamicPost(strconv.Itoa(appData.AppId), appData.OwnerId, createOrgUrl(), strings.NewReader(orgPayload))
	if err != nil {
		return err
	}

	var orgDataResp Organizations
	if err := json.NewDecoder(bytes.NewReader(orgRes)).Decode(&orgDataResp); err != nil {
		return err
	}

	return nil
}

func InviteUser(orgId string, invite SendInvitation) error {

	payload, err := json.Marshal(invite)
	if err != nil {
		return err
	}

	_, err = Post(sendInvitationUrl(orgId), bytes.NewReader(payload))
	if err != nil {
		return err
	}

	return nil
}

func GetAllRolesOfAnOrg(orgId string) ([]RoleResponse, error) {
	data, err := Get(getAllRolesOfAnOrg(orgId), "")
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

func GetAllOrganizationsOfTenant(orgId string) ([]AllOrganizationsResponse, error) {

	//get AppId from the orgId-appId relation
	appId, err := GetAppIdFromOrgIdMapping(mongoClient, orgId)
	if err != nil {
		return nil, err
	}

	data, err := Get(getOrgsOfTenantUrl(), strconv.Itoa(appId))
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

func GetAllInvitationsOfOrganization(orgID string) ([]InvitationResponse, error) {
	data, err := Get(getAllInvitationsOfOrganization(orgID), "")
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
	data, err := Get(getAllRolesOfUserInOrg(orgID, uid), "")
	if err != nil {
		return nil, err
	}

	var userRolesResp struct {
		Data []UserRole `json:"Data"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&userRolesResp); err != nil {
		return nil, err
	}

	var rolesData []RoleResponse
	for _, userRole := range userRolesResp.Data {
		roles, err := Get(getRoleById(userRole.RoleId), "")
		if err != nil {
			return nil, err
		}

		var roleResp RoleResponse
		if err := json.NewDecoder(bytes.NewReader(roles)).Decode(&roleResp); err != nil {
			return nil, err
		}

		rolesData = append(rolesData, roleResp)
	}

	return rolesData, nil
}

func Test() {
	if err := CreateOrg("qwer", "niketest123", ""); err != nil {
		fmt.Println("Error:", err)
	}
}

func TestLogin() {
	// orgs, _ := Login("nike@mailinator.com", "123456")
	// fmt.Println("Response:", orgs)
}

func TestGetAllUsersOfAnOrganization() {
	users, _ := GetAllUsersOfAnOrganization("org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", users)
}

func TestGetAllRolesOfAnOrg() {
	roles, _ := GetAllRolesOfAnOrg("org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", roles)
}

func TestGetAllOrganizationsOfTenant() {
	orgs, _ := GetAllOrganizationsOfTenant("org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", orgs)
}

func TestGetAllInvitationsOfOrganization() {
	invitations, _ := GetAllInvitationsOfOrganization("org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", invitations)
}

func TestGetAllRolesOfUserInOrg() {
	roles, _ := GetAllRolesOfUserInOrg("org_Z5pkOZ-0eGkkbhQ1", "auth_0eGkkbhQ1")
	fmt.Println("Response:", roles)
}
