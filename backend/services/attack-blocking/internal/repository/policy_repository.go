package repository

// PolicyRepository defines methods for managing policies in storage.
type PolicyRepository interface {
    CreatePolicy(policy interface{}) error
    GetPolicy(id string) (interface{}, error)
    ListPolicies() ([]interface{}, error)
    DeletePolicy(id string) error
} 