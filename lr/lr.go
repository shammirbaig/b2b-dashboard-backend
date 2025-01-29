package lr

const (
	lrURL = "https://devapi.lrinternal.com"
)

func generateM2MToken() {
	// Generate a new M2M token
}

func Test() {
	Get(lrURL + "/v2/manage/permissions")
}
