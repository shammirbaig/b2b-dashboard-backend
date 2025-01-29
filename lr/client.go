package lr

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Get(url string) {
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("X-CustomerId", "7c9254057e2044c5b3fadf8bf0b3dd31")
	req.Header.Add("X-AppId", "99207378")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImQzMTNjZTcxZTY0MDRjZjE4ZGVhNTNlOWViNmI2MWMyIiwidHlwIjoiYXQrand0In0.eyJpc3MiOiJodHRwczovL2RlbW90ZXN0aW5nLmRldmh1Yi5scmludGVybmFsLmNvbS8iLCJzdWIiOiI2OTlmNDBhZS04MDQwLTQyZmUtYTFjZC1lMWViMTBjZjczZWZAY2xpZW50IiwiYXVkIjpbImh0dHBzOi8vZGV2YXBpLmxyaW50ZXJuYWwuY29tL2lkZW50aXR5L3YyL21hbmFnZSJdLCJleHAiOjE3MzgxNzI3OTEsIm5iZiI6MTczODE2NTU5MSwiaWF0IjoxNzM4MTY1NTkxLCJqdGkiOiJlMTQ1YWNmZi05N2MxLTQwZTMtOGUzOC0wZjgxYzZmYTUzZWIiLCJjaWQiOiJkZGZmOGE2My1jYmMzLTQ3MjMtODQxNS1iOTEwYzRkODc3MGQiLCJzaWQiOiI2NzlhNGQ1Ny03ODgwLTQzMWMtOTkxYi02YjEzMGI0ZDAyMGYiLCJzY3AiOlsiYWxsIiwiY3JlYXRlOnVzZXJzIiwicmVhZDp1c2VycyIsInVwZGF0ZTp1c2VycyIsImRlbGV0ZTp1c2VycyIsImNyZWF0ZTplbWFpbF92ZXJpZmljYXRpb25fdG9rZW4iLCJjcmVhdGU6Zm9yZ290X3Bhc3N3b3JkX3Rva2VuIiwidXBkYXRlOnVzZXJzX2FjY2Vzc190b2tlbiIsInJlYWQ6dXNlcnNfcHJpdmFjeV9wb2xpY3lfaGlzdG9yeSIsImNyZWF0ZTpyb2xlcyIsInJlYWQ6cm9sZXMiLCJ1cGRhdGU6cm9sZXMiLCJkZWxldGU6cm9sZXMiLCJjcmVhdGU6Y3VzdG9tX29iamVjdHMiLCJyZWFkOmN1c3RvbV9vYmplY3RzIiwidXBkYXRlOmN1c3RvbV9vYmplY3RzIiwiZGVsZXRlOmN1c3RvbV9vYmplY3RzIiwiY3JlYXRlOnNvdHQiLCJjcmVhdGU6dG9rZW5zIiwicmV2b2tlOnRva2VucyIsInJlYWQ6Y29uc2VudGxvZ3MiLCJ2ZXJpZnk6cmVhdXRoIiwicmVhZDp3ZWJob29rcyIsImNyZWF0ZTp3ZWJob29rcyIsInVwZGF0ZTp3ZWJob29rcyIsImRlbGV0ZTp3ZWJob29rcyIsInJlYWQ6d2ViaG9va19sb2dzIl0sImd0eSI6ImNsaWVudF9jcmVkZW50aWFscyJ9.T9djdHeXKbrDBV-sHgh3UvDKl--jgVFa-uYOEMQRGtpi7d_sNOerOahbF9Mn_sfPrP5dbd-piH-2HPhpmmOeMFQfwaAQOXkXw4_JfsTRuCIr4Jx3aYDdCd8HCX8JdbNRNpyrGJ54rZf8mG5uh0H_aA3bpULmDqDC6mqW6wL3EEc")
	req.Header.Add("Cookie", "__cf_bm=poEledIwwQoD3mx542YKDtJh8qCGJIXZwvkV.uampTc-1737548546-1.0.1.1-uXnnSebq3lSMAeDbqVjCAOoYrLE_E_XCE8CYKKXm6RXxiGZOsonnZGVZZf.ShWd53KhaJqq6BXEaNYTCh7rtqg")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}
