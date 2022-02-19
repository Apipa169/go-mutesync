package mutesync

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

const Port int = 8249

func GetStatus(host, token string) (status Status, err error) {
	var statusResponse statusResponse

	req, err := http.NewRequest(http.MethodGet, "http://" + host + ":" + strconv.Itoa(Port) + "/state", nil)
	if err != nil {
		return statusResponse.Status, err
	}

	req.Header.Set("Authorization", "Token " + token)
	req.Header.Set("x-mutesync-api-version", "1")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return statusResponse.Status, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return statusResponse.Status, err
	}

	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		return statusResponse.Status, err
	}

	return statusResponse.Status, nil
}

func Authenticate(host string) (token string, err error) {
	resp, err := http.Get("http://" + host + ":" + strconv.Itoa(Port) + "/authenticate")
	if err != nil {
		return token, err
	}

	if resp.StatusCode == http.StatusForbidden {
		return token, errors.New("make sure mutesync allows external apps")
	}

	if resp.StatusCode != http.StatusOK {
		return token, errors.New("unexpected response from mutesync")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}

	var authResponse authResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return token, err
	}

	token = authResponse.Token

	return token, nil
}
