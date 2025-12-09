package models

// Policy represents a security policy for attack blocking.
type Policy struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Active      bool   `json:"active"`
} 