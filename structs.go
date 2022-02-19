package mutesync

type authResponse struct {
	Token string `json:"token"`
}

type statusResponse struct {
	Status Status `json:"data"`
}

type Status struct {
	InMeeting bool   `json:"in_meeting"`
	Hostname  string `json:"hostname"`
	UserID    string `json:"user-id"`
	Muted     bool   `json:"muted"`
}
