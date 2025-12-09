package repository

// BlockingRepository defines methods for managing blocking rules in storage.
type BlockingRepository interface {
    CreateBlockingRule(rule interface{}) error
    GetBlockingRule(id string) (interface{}, error)
    ListBlockingRules() ([]interface{}, error)
    DeleteBlockingRule(id string) error
} 