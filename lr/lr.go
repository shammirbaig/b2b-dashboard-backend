package lr

const (
	lrURL = "https://devapi.lrinternal.com"
)

func generateM2MToken() string {
	// Generate a new M2M token
	return "m2mToken"
}

func Test() {
	Get(lrURL + "/v2/manage/permissions")
}
