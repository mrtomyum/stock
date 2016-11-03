package model

import sys "github.com/mrtomyum/sys/model"

type Unit struct {
	sys.Base
	Name   string `json:"name" db:"name"`
	NameEn string `json:"name_en" db:"name_en"`
}
