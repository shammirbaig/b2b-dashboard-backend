package lr

import (
	"fmt"
	"strings"
)

func generateM2MToken() string {
	// Generate a new M2M token
	return "m2mToken"
}

func Login(email, password string) (string, error) {

	data, err := Post(loginUrl, strings.NewReader(`{"email":"`+email+`","password":"`+password+`"}`))
	if err != nil {
		return "", err
	}

	fmt.Println(string(data))

	return string(data), nil
}

func Test() {
	Get(lrURL + "/v2/manage/permissions")
}
