package model

import "github.com/mrtomyum/nava-sys/model"

type Category struct {
	model.Base
	ParentID uint64
	TH       string
	EN       string
}
