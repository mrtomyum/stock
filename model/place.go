package model

type PlaceType int

const (
	NO_PLACE PlaceType = iota
	DORM
	FACTORY
	OFFICE
	SCHOOL
	PATHWAY
	MARKET_PLACE
	DEPARTMENT_STORE
	MODERN_TRADE
)

type Place struct {
	Base
	ClientID uint64    `json:"client_id"`
	NameTh   string    `json:"name_th" db:"name_th"`
	NameEn   string    `json:"name_en" db:"name_en"`
	Address  string    `json:"address"`
	Type     PlaceType `json:"type"`
	Lat      float64   `json:"lat"`
	Lng      float64   `json:"lng"`
}

func (p *Place) GetByMachineCode() string {
	return p.NameTh
}