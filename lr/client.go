package lr

import (
	"io"
	"net/http"
	"strings"
)

func Get(url string) ([]byte, error) {
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("X-CustomerId", "7c9254057e2044c5b3fadf8bf0b3dd31")
	req.Header.Add("X-AppId", "99207378")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+generateM2MToken())
	req.Header.Add("Cookie", "__cf_bm=poEledIwwQoD3mx542YKDtJh8qCGJIXZwvkV.uampTc-1737548546-1.0.1.1-uXnnSebq3lSMAeDbqVjCAOoYrLE_E_XCE8CYKKXm6RXxiGZOsonnZGVZZf.ShWd53KhaJqq6BXEaNYTCh7rtqg")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
