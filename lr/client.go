package lr

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func dynamicClient(appid, customerid string, method, url string, payload io.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("X-CustomerId", customerid)
	req.Header.Add("X-AppId", appid)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+GetM2MToken())
	req.Header.Add("Access-Control-Allow-Origin", "*")

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

func client(method, url string, payload io.Reader, appId string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("X-CustomerId", "7c9254057e2044c5b3fadf8bf0b3dd31")
	req.Header.Add("X-AppId", func() string {
		if appId == "" {
			return "99207378"
		}
		return appId
	}())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+GetM2MToken())
	req.Header.Add("Access-Control-Allow-Origin", "*")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fmt.Println(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Get(url, appId string) ([]byte, error) {
	method := "GET"

	payload := strings.NewReader(``)

	return client(method, url, payload, appId)
}

func Post(url string, payload io.Reader) ([]byte, error) {

	method := "POST"

	return client(method, url, payload, "")
}

func DynamicPost(appid, customerid string, url string, payload io.Reader) ([]byte, error) {

	method := "POST"

	return dynamicClient(appid, customerid, method, url, payload)
}

func Put(url string, payload io.Reader) ([]byte, error) {

	method := "PUT"

	return client(method, url, payload, "")
}
