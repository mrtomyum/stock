package model

import "github.com/mrtomyum/sys/model"

type Category struct {
	model.Base
	ParentID uint64 `json:"parent_id" db:"parent_id"`
	NameTh   string `json:"name_th" db:"name_th"`
	NameEn   string `json:"name_en" db:"name_en"`
}
