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

func Test() {
	orgs, _ := GetUserOrgs("7c9254057e2044c5b3fadf8bf0b3dd31")

	fmt.Println("Response:", orgs)
}
