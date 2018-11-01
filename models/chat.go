package models

type Chat struct {
	Register  Register           `json:"register"`
	Responses []RegisterResponse `json:"responses"`
}
