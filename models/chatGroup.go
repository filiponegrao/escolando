package models

type ChatGroup struct {
	Register   Register          `json:"register"`
	Targets    []RegisterContact `json:"targets"`
	TargetName string            `json:"targetName"`
}
