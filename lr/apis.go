package lr

import (
	"strings"
)

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

func getAllInvitationsOfOrganization(orgID string) string {
	url := lrURL + "/v2/manage/organizations/{orgId}/invitations"
	return strings.Replace(url, "{orgId}", orgID, 1)
}
