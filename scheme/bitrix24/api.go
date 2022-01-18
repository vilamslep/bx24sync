package bitrix24

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BitrixRestResponseFind struct {
	Result []BitrxiRestResult `json:"result"`
	Total  int                `json:"total"`
}

type BitrixRestResponseUpdate struct {
	Result bool `json:"result"`
}

type BitrixRestResponseAdd struct {
	Result int `json:"result"`
}

type BitrxiRestResult struct {
	ID string `json:"ID"`
}

func checkResponseFind(res *http.Response) (response BitrixRestResponseFind, err error) {

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		if err != nil {
			return response, err
		}
		return response, fmt.Errorf(fmt.Sprintf("%s: %s", "bad request", string(body)))
	} else {
		return fillResponseFind(body)
	}
}

func checkResponseUpdate(res *http.Response) (response BitrixRestResponseUpdate, err error) {

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		if err != nil {
			return response, err
		}
		return response, fmt.Errorf(fmt.Sprintf("%s: %s", "bad request", string(body)))
	} else {
		return fillResponseUpdate(body)
	}
}

func checkResponseAdd(res *http.Response) (response BitrixRestResponseAdd, err error) {

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		if err != nil {
			return response, err
		}
		return response, fmt.Errorf(fmt.Sprintf("%s: %s", "bad request", string(body)))
	} else {
		return fillResponseAdd(body)
	}
}

func fillResponseFind(bodyRaw []byte) (response BitrixRestResponseFind, err error) {
	err = json.Unmarshal(bodyRaw, &response)
	return response, err
}

func fillResponseUpdate(bodyRaw []byte) (response BitrixRestResponseUpdate, err error) {
	err = json.Unmarshal(bodyRaw, &response)
	return response, err
}

func fillResponseAdd(bodyRaw []byte) (response BitrixRestResponseAdd, err error) {
	err = json.Unmarshal(bodyRaw, &response)
	return response, err
}

func ExecReq(method string, url string, rd io.Reader) (*http.Response, error) {

	if req, err := http.NewRequest(method, url, rd); err == nil {
		req.Header.Set("Content-Type", "application/json")
		client := http.Client{Timeout: time.Second * 300}
		return client.Do(req)
	} else {
		return nil, err
	}
}