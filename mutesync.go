package mutesync

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strconv"
)

const port = 8249
const pathAuth = "authenticate"
const pathState = "state"

func InitMutesyncClient(host string) (mc MutesyncClient, err error) {
    mc.Host = host

    mc.Token, err = authenticate(host)
    if err != nil {
        return mc, err
    }

    return mc, nil
}

func getStatus(host, token string) (Status, error) {
    var sr statusResp

    body, err := doRequest(host, pathState, &token)
    if err != nil {
        return sr.Status, err
    }

    err = json.Unmarshal(body, &sr)
    if err != nil {
        return sr.Status, err
    }

    return sr.Status, nil
}

func isInMeeting(host, token string) (bool, error) {
    status, err := getStatus(host, token)
    if err != nil {
        return false, err
    }

    return status.InMeeting, nil
}

func isMuted(host, token string) (bool, error) {
    status, err := getStatus(host, token)
    if err != nil {
        return false, err
    }

    return status.Muted, nil
}

func authenticate(host string) (string, error) {
    var ar authResp

    body, err := doRequest(host, pathAuth, nil)
    if err != nil {
        return ar.Token, err
    }

    err = json.Unmarshal(body, &ar)
    if err != nil {
        return ar.Token, err
    }

    return ar.Token, nil
}

func doRequest(host, path string, token *string) (body []byte, err error){
    url := "http://" + host + ":" + strconv.Itoa(port) + "/" + path
    req, err := http.NewRequest(http.MethodGet, url, nil)
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

    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return body, err
    }

    switch resp.StatusCode {
    case http.StatusOK:
        break
    case http.StatusForbidden:
        return body, ErrAuthFailed{Reason: getAuthFailedReason(body), Path: path}
    default:
        return body, ErrUnexpectedResponse{}
    }

    return body, nil
}

func getAuthFailedReason(body []byte) (reason string) {
    var er errorResp

    err := json.Unmarshal(body, &er)
    if err != nil {
        return "unknown error"
    }

    return er.ErrorMsg
}
