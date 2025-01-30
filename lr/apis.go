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

func createAppUrl() string {
	return lrURL + "/v2/manage/app"
}

func createOrgUrl() string {
	return lrURL + "/v2/manage/organizations"
}

func sendInvitationUrl(orgId string) string {
	url := lrURL + "/v2/manage/organizations/{orgId}/invitations"
	return strings.Replace(url, "{orgId}", orgId, 1)
}

func getUsersOfOrgUrl(orgID string) string {

	url := lrURL + "/v2/manage/organizations/{orgId}/orgcontext"

	return strings.Replace(url, "{orgId}", orgID, 1)
}

func getAllRolesOfAnOrg(orgId string) string {
	url := lrURL + "/v2/manage/organizations/{orgId}/roles"
	return strings.Replace(url, "{orgId}", orgId, 1)
}

func getOrgsOfTenantUrl() string {
	return lrURL + "/v2/manage/organizations"
}

func turnB2BApp(cid, appid string) string {
	url := "https://devpartner.lrinternal.com/v2/customer/{cid}/app/{appid}/feature?partnerKey=C6qM8fngbypF7z9BP5DzkzRfd4QShSQMP3JATjSnnAzn5zhNHXokerLcKYQbmo9pBibH36miKi&partnerSecret=AnyzXDkBt5sbJb6SDpjrYj4eLBAgmr49JETxePHhApmHnTdBYJoBrFtTbdpiPoKSsYhdReFNe9gTXFPfbJb69"
	u := strings.Replace(url, "{cid}", cid, 1)
	return strings.Replace(u, "{appid}", appid, 1)
}
