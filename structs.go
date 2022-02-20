package mutesync

import "fmt"

type authResp struct {
    Token string `json:"token"`
}

type statusResp struct {
    Status Status `json:"data"`
}

type errorResp struct {
    ErrorMsg string `json:"error"`
}

type Status struct {
    InMeeting bool   `json:"in_meeting"`
    Hostname  string `json:"hostname"`
    UserID    string `json:"user-id"`
    Muted     bool   `json:"muted"`
}

type ErrAuthFailed struct {
    Reason string
    Path string
}

func (ef ErrAuthFailed) Error() string {
    return fmt.Sprintf("error at %s. %s", ef.Path, ef.Reason)
}

type ErrUnexpectedResponse struct {
}

func (ur ErrUnexpectedResponse) Error() string {
    return "unexpected response from mutesync"
}
