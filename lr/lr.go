package lr

import (
	"strings"
)

func generateM2MToken() string {
	// Generate a new M2M token
	return "m2mToken"
}

// Login returns the orgs that the user is a part of
func Login(email, password string) ([]string, error) {

	data, err := Post(loginUrl, strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	if err != nil {
		return nil, err
	}

	// get the UID
	uid := string(data)

	// using the above UID, get the orgs that the user is a part of
	GetUserRoles(uid)

	return nil, nil
}

func GetUserRoles(userID string) ([]string, error) {
	data, err := Get(getRolesOfUserInOrgUrl(userID))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func CreateAnOrg(orgName string) (string, error) {
	data, err := Post(lrURL+"/v2/manage/organizations", strings.NewReader(`{"name":"`+orgName+`"}`))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func Test() {
	Get(lrURL + "/v2/manage/permissions")
}
