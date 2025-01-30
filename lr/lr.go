package lr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
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

	for i, org := range orgsResp.Data {
		name, _ := GetAnOrganizationDetailsName(org.OrgId)
		orgsResp.Data[i].Name = name
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

	for i, user := range usersResp.Data {
		name, email, userName, err := GetProfileDetail(user.Uid)
		if err != nil {
			return nil, err
		}

		usersResp.Data[i].Name = name
		usersResp.Data[i].Email = email
		usersResp.Data[i].Username = userName
	}

	return usersResp.Data, nil
}

func CreateOrg(tenantOrgID string, newOrgName string, oldOrgName string) error {

	// create an app : 7c9254057e2044c5b3fadf8bf0b3dd31

	// AppName should be current org name and mark its name as appname
	var payload = struct {
		AppName     string `json:"AppName"`
		Domain      string `json:"Domain"`
		CallbackUrl string `json:"CallbackUrl"`
		DevDomain   string `json:"DevDomain"`
	}{
		AppName:     oldOrgName,
		Domain:      "loginradius.com",
		CallbackUrl: "https://loginradius.com",
		DevDomain:   "loginradius.com",
	}

	v, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	payloadReader := bytes.NewReader(v)

	res, err := Post(createAppUrl(), payloadReader)
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
	var orgPayload = struct {
		Name    string `json:"Name"`
		Domains []struct {
			DomainName string `json:"DomainName"`
			IsVerified bool   `json:"IsVerified"`
		} `json:"Domains"`
		IsAuthRestrictedToDomain bool `json:"IsAuthRestrictedToDomain"`
	}{
		Name: newOrgName,
		Domains: []struct {
			DomainName string "json:\"DomainName\""
			IsVerified bool   "json:\"IsVerified\""
		}{
			{
				DomainName: "loginradius.com",
				IsVerified: true,
			},
		},
		IsAuthRestrictedToDomain: false,
	}

	// orgPayload := `{
	// 	"Name": "{orgName}",
	// 	"Domains": [
	// 		{
	// 			"DomainName": "loginradius.com",
	// 			"IsVerified": true
	// 		}
	// 	],
	// 	"IsAuthRestrictedToDomain": false
	// }`

	v2, err := json.Marshal(orgPayload)
	if err != nil {
		return err
	}

	orgPayloadReader := bytes.NewReader(v2)

	// this creates in nike-com db but how do we differentiate the hirearcy
	_, err = DynamicPost(strconv.Itoa(appData.AppId), appData.OwnerId, createOrgUrl(), orgPayloadReader)
	if err != nil {
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

	tenantData, err := Get(getAllTenantRoles(), "")
	if err != nil {
		return nil, err
	}

	var tenantRolesResp struct {
		Data []RoleResponse `json:"Data"`
	}

	if err := json.NewDecoder(bytes.NewReader(tenantData)).Decode(&tenantRolesResp); err != nil {
		return nil, err
	}

	rolesResp.Data = append(rolesResp.Data, tenantRolesResp.Data...)

	return rolesResp.Data, nil
}

func GetAllOrganizationsOfTenant(orgId string) ([]AllOrganizationsResponse, error) {

	//get AppId from the orgId-appId relation
	appId, err := GetAppIdFromOrgIdMapping(mongoClient, orgId)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return make([]AllOrganizationsResponse, 0), nil
		}

		return nil, err
	}

	if appId == 0 {
		return make([]AllOrganizationsResponse, 0), nil
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

func GetAllInvitationsOfOrganization(orgID string) (*InvitationResponse, error) {
	data, err := Get(getAllInvitationsOfOrganization(orgID), "")
	if err != nil {
		return nil, err
	}

	var invitationsResp InvitationResponse
	if err := json.Unmarshal(data, &invitationsResp); err != nil {
		return nil, err
	}

	return &invitationsResp, nil
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

func GetAnOrganizationDetailsName(orgID string) (string, error) {
	data, err := Get(getOrganizationDetails(orgID), "")
	if err != nil {
		return "", err
	}

	var orgsResp OrganizationResponse

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&orgsResp); err != nil {
		return "", err
	}

	return orgsResp.Name, nil
}

func GetProfileDetail(uid string) (string, string, string, error) {
	data, err := Get(getProfileDetail(uid), "")
	if err != nil {
		return "", "", "", err
	}

	var profileResp struct {
		Uid       string       `json:"Uid"`
		FirstName string       `json:"FirstName"`
		LastName  string       `json:"LastName"`
		Email     []EmailValue `json:"Email"`
		UserName  string       `json:"UserName"`
	}

	if err := json.Unmarshal(data, &profileResp); err != nil {
		return "", "", "", err
	}

	return profileResp.FirstName + " " + profileResp.LastName, profileResp.Email[0].Value, profileResp.UserName, nil
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
