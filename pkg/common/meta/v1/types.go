package v1

// ListMeta describes metadata that synthetic resources must have, including lists and
type ListMeta struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"pageSize,omitempty"`
}
