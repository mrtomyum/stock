package model

type Category struct {
	Base
	ParentID uint64
	TH       string
	EN       string
}
