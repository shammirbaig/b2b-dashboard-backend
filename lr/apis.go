package lr

import "strings"

const (
	lrURL  = "https://devapi.lrinternal.com"
	apikey = "apikey=6453562a-5411-466c-96f1-4a6c09817ce5"

	loginUrl = lrURL + "/identity/v2/auth/login?" + apikey
)

func getOrgsOfUserUrl(uid string) string {

	url := lrURL + "/identity/v2/manage/account/{uid}/orgcontext"

	return strings.Replace(url, "{uid}", uid, 1)
}
