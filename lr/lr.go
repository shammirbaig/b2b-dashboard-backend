package lr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Login returns the orgs that the user is a part of
func Login(email, password string) ([]Organizations, error) {

	data, err := Post(loginUrl, strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	if err != nil {
		return nil, err
	}

	var tokenResp struct {
		Uid string `json:"Uid"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&tokenResp); err != nil {
		return nil, err
	}

	orgs, err := GetUserOrgs(tokenResp.Uid)
	if err != nil {
		return nil, err
	}

	fmt.Println(orgs)
	return orgs, nil
}

func GetUserOrgs(uid string) ([]Organizations, error) {
	data, err := Get(getOrgsOfUserUrl(uid))
	if err != nil {
		return nil, err
	}

	var orgsResp struct {
		Data struct {
			Orgs []Organizations `json:"Orgs"`
		} `json:"data"`
	}

	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&orgsResp); err != nil {
		return nil, err
	}

	return orgsResp.Data.Orgs, nil
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

func Test() {
	if err := CreateOrg("qwer", "niketest123"); err != nil {
		fmt.Println("Error:", err)
	}

	// fmt.Println("Response:", orgs)
}
