package models

// BlockingRule represents a rule for blocking malicious traffic.
type BlockingRule struct {
    ID          string `json:"id"`
    RuleType    string `json:"rule_type"`
    Description string `json:"description"`
    Enabled     bool   `json:"enabled"`
} 