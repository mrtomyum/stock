package model

type Location struct {
	ID       uint64  `json:"id"`
	ParentID uint64  `json:"parent_id"`
	Name     string  `json:"name"`
	Type     LocType `json:"type"`
}
