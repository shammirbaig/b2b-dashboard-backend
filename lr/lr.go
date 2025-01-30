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

func Test() {
	orgs, _ := GetUserOrgs("7c9254057e2044c5b3fadf8bf0b3dd31")

	fmt.Println("Response:", orgs)
}

func TestLogin() {
	orgs, _ := Login("nike@mailinator.com", "123456")
	fmt.Println("Response:", orgs)
}

func TestGetAllUsersOfAnOrganization() {
	users, _ := GetAllUsersOfAnOrganization("org_Z5pkOZ-0eGkkbhQ1")
	fmt.Println("Response:", users)
}
