package lr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Login returns the orgs that the user is a part of
func Login(email, password string) ([]OrganizationResponse, error) {

	data, err := Post(loginUrl, strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	if err != nil {
		return nil, err
	}

	var tokenResp struct {
		Profile struct {
			Uid string `json:"Uid"`
		} `json:"Profile"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&tokenResp); err != nil {
		return nil, err
	}

	orgs, err := GetUserOrgs(tokenResp.Profile.Uid)
	if err != nil {
		return nil, err
	}

	fmt.Println(orgs)
	return orgs, nil
}

func GetUserOrgs(uid string) ([]OrganizationResponse, error) {
	data, err := Get(getOrgsOfUserUrl(uid))
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
	data, err := Get(getUsersOfOrgUrl(orgID))
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

func CreateOrg(tenantOrgID string, newOrgName string) error {

	// create an app : 7c9254057e2044c5b3fadf8bf0b3dd31

	// AppName should be current org name and mark its name as appname
	payload := strings.NewReader(`{
		"AppName": "niketest123", 
		"Domain": "login2website.com; localhost:3000",
		"CallbackUrl": "login2website.com; localhost",
		"DevDomain": "login2website.com"
	}`)

	res, err := Post(createAppUrl(), payload)
	if err != nil {
		return err
	}

	var appData AppResponse
	if err := json.NewDecoder(bytes.NewReader(res)).Decode(&appData); err != nil {
		return err
	}

	// fmt.Printf("AppData: %+v %+v\n", appData, appData.AppId)

	// create a relation between orgid and appid
	if err := CreateAppidToOrgidMapping(mongoClient, appData.AppId, tenantOrgID); err != nil {
		return err
	}

	// create an org for the above passed data
	orgPayload := strings.NewReader(`{
		"Name": "test-nike-hyd",
		"Domains": [
			{
				"DomainName": "loginradius.com",
				"IsVerified": true
			}
		],
		"IsAuthRestrictedToDomain": false
	}`)

	// this creates in nike-com db but how do we differentiate the hirearcy
	orgRes, err := Post(createOrgUrl(), orgPayload)
	if err != nil {
		return err
	}

	fmt.Println("OrgRes:", string(orgRes))
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

func Test() {
	if err := CreateOrg("qwer", "niketest123"); err != nil {
		fmt.Println("Error:", err)
	}

	// fmt.Println("Response:", orgs)
}

func TestLogin() {
	orgs, _ := Login("nike@mailinator.com", "123456")
	fmt.Println("Response:", orgs)
}

func TestGetAllUsersOfAnOrganization() {
	users, _ := GetAllUsersOfAnOrganization("org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", users)
}
