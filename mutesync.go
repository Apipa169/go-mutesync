package mutesync

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strconv"
)

const port int = 8249
const pathAuth = "/authenticate"
const pathState = "/state"

func GetStatus(host, token string) (Status, error) {
    var statusResponse statusResponse

    body, err := doRequest(host, pathState, &token)
    if err != nil {
        return statusResponse.Status, err
    }

    err = json.Unmarshal(body, &statusResponse)
    if err != nil {
        return statusResponse.Status, err
    }

    return statusResponse.Status, nil
}

func Authenticate(host string) (string, error) {
    var authResponse authResponse

    body, err := doRequest(host, pathAuth, nil)
    if err != nil {
        return authResponse.Token, err
    }

    err = json.Unmarshal(body, &authResponse)
    if err != nil {
        return authResponse.Token, err
    }

    return authResponse.Token, nil
}

func doRequest(host, path string, token *string) (body []byte, err error){
    req, err := http.NewRequest(http.MethodGet, "http://" + host + ":" + strconv.Itoa(port) + path, nil)
    if err != nil {
        return body, err
    }

    req.Header.Set("x-mutesync-api-version", "1")
    if token != nil {
        req.Header.Set("Authorization", "Token " + *token)
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return body, err
    }

    switch resp.StatusCode {
    case http.StatusForbidden:
        return body, ErrAuthFailed{Reason: getAuthFailedReason(path), Path: path}
    case http.StatusOK:
        break
    default:
        return body, ErrUnexpectedResponse{}
    }

    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return body, err
    }

    return body, nil
}

func getAuthFailedReason(path string) (reason string) {
    switch path {
    case pathAuth:
        return "make sure mutesync allows external apps"
    case pathState:
        return "invalid token"
    }

    return "forbidden"
}
