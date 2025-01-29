package lr

import (
	"io"
	"net/http"
	"strings"
)

func client(method, url string, payload io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("X-CustomerId", "7c9254057e2044c5b3fadf8bf0b3dd31")
	req.Header.Add("X-AppId", "99207378")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+generateM2MToken())

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

func Get(url string) ([]byte, error) {
	method := "GET"

	payload := strings.NewReader(``)

	return client(method, url, payload)
}

func Post(url string, payload io.Reader) ([]byte, error) {

	method := "POST"

	return client(method, url, payload)
}
