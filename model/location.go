package model

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type LocType int

const (
	ROOT LocType = 1 + iota
	STORE
	VAN
	MACHINE
	COLUMN
	VENDOR
	INSPECTION
	DAMAGE
)

func (lt LocType) MarshalJSON() ([]byte, error) {
	ltString, ok := map[LocType]string{
		ROOT:       "ROOT",
		STORE:      "STORE",
		VAN:        "VAN",
		MACHINE:    "MACHINE",
		COLUMN:     "COLUMN",
		VENDOR:     "VENDOR",
		INSPECTION: "INSPECTION",
		DAMAGE:     "DAMAGE",
	}[lt]
	if !ok {
		return nil, fmt.Errorf("invalid Location Type value %v", lt)
	}
	return json.Marshal(ltString)
}

type Location struct {
	Base
	Name     JsonNullString `json:"name"`
	Code     JsonNullString `json:"code"`
	Type     LocType        `json:"type"` // TODO: return LocType ENUM in string
	ParentID uint64         `json:"-" db:"parent_id"`
	Child    []*Location    `json:"nodes,omitempty"`
}

func (this *Location) Size() int {
	var size int = len(this.Child)
	for _, c := range this.Child {
		size += c.Size()
	}
	return size
}

func (this *Location) Add(nodes ...*Location) bool {
	var size = this.Size()
	for _, node := range nodes {
		if node.ParentID == this.ID {
			this.Child = append(this.Child, node)
		} else {
			for _, c := range this.Child {
				if c.Add(node) {
					break
				}
			}
		}
	}
	return this.Size() == size + len(nodes)
}

func (l *Location) All(db *sqlx.DB) ([]*Location, error) {
	locations := []*Location{}
	sql := `SELECT * FROM location`
	err := db.Select(&locations, sql)
	if err != nil {
		log.Fatal("Error in model.Select..", err)
		return nil, err
	}
	return locations, nil
}

func (l *Location) Show(db *sqlx.DB) ([]*Location, error) {
	// TODO: แก้ Select ให้สามารถ Recursive  Parent_id ได้
	sql := `
		SELECT * FROM location
		WHERE id = ?
		OR parent_id = ?
	`
	locations := []*Location{}
	err := db.Select(&locations, sql, l.ID, l.ID)
	if err != nil {
		log.Fatal("Error in model.Select..", err)
		return nil, err
	}
	return locations, nil
}
