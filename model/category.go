package model

import "github.com/mrtomyum/sys/model"

type Category struct {
	model.Base
	ParentId uint64 `json:"parent_id" db:"parent_id"`
	Name     string `json:"name" db:"name"`
	NameEn   string `json:"name_en" db:"name_en"`
}
