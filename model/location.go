package model

import (
	//"encoding/json"
	//"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"fmt"
	"encoding/json"
	"github.com/mrtomyum/nava-sys/model"
)

type LocType int

const (
	ROOT LocType = iota
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
	model.Base
	Name     string `json:"name"` //Todo: Has problem with custom type JsonNullString can't receive value from json.NewDecoder()
	Code     string `json:"code"`
	Type     LocType        `json:"type"`
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
	log.Println("call method Location.All()")
	locations := []*Location{}
	sql := `SELECT * FROM location`
	err := db.Select(&locations, sql)
	if err != nil {
		log.Fatal("Error in model.Select..", err)
		return nil, err
	}
	log.Println(locations)
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

func (l *Location) New(db *sqlx.DB) (*Location, error) {
	sql := `
		INSERT INTO location (
			name,
			code,
			type,
			parent_id
		)
		VALUES (?, ?, ?, ?)
	`
	log.Println("Test Location receiver:", l.Name, l.Code, l.Type, l.ParentID)
	rsp, err := db.Exec(sql,
		l.Name,
		l.Code,
		l.Type,
		l.ParentID,
	)
	if err != nil {
		log.Println("Error db.Exec in model.Location.Show", err)
		return nil, err
	}
	id, _ := rsp.LastInsertId()
	newLoc := Location{}
	err = db.Get(&newLoc, "SELECT * FROM location WHERE id =?", id)
	if err != nil {
		log.Println("Error db.GET in model.Location.Show", err)
	}
	return &newLoc, nil
}