package models

// Agent represents an agent responsible for enforcing blocking rules.
type Agent struct {
    ID       string `json:"id"`
    Hostname string `json:"hostname"`
    Status   string `json:"status"`
} 