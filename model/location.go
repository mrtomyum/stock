package model

import (
	"encoding/json"
	"fmt"
	sys "github.com/mrtomyum/sys/model"
	"log"
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
	typeStr, ok := map[LocType]string{
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
	return json.Marshal(typeStr)
}

type Location struct {
	sys.Base
	Code     string      `json:"code"`
	Type     LocType     `json:"type"`
	ParentID uint64      `json:"-" db:"parent_id"`
	Child    []*Location `json:"nodes,omitempty" db:"-"`
}

func (this *Location) Size() int {
	var size int = len(this.Child)
	for _, c := range this.Child {
		size += c.Size()
	}
	return size
}

func (this *Location) AddTree(nodes ...*Location) bool {
	var size = this.Size()
	for _, node := range nodes {
		if node.ParentID == this.ID {
			this.Child = append(this.Child, node)
		} else {
			for _, c := range this.Child {
				if c.AddTree(node) {
					break
				}
			}
		}
	}
	return this.Size() == size + len(nodes)
}

func (l *Location) All() ([]*Location, error) {
	log.Println("call method Location.All()")
	locations := []*Location{}
	sql := `SELECT * FROM location`
	err := DB.Select(&locations, sql)
	if err != nil {
		log.Fatal("Error in model.Select..", err)
		return nil, err
	}
	log.Println(locations)
	return locations, nil
}

func (l *Location) Get() ([]*Location, error) {
	// TODO: แก้ Select ให้สามารถ Recursive  Parent_id ได้
	sql := `
		SELECT * FROM location
		WHERE id = ?
		OR parent_id = ?
		`
	locations := []*Location{}
	err := DB.Select(&locations, sql, l.ID, l.ID)
	if err != nil {
		log.Fatal("Error in model.Select..", err)
		return nil, err
	}
	return locations, nil
}

func (l *Location) Insert() (*Location, error) {
	sql := `
		INSERT INTO location (
			code,
			type,
			parent_id
		)
		VALUES (?, ?, ?, ?)
	`
	log.Println("Test Location receiver:", l.Code, l.Type, l.ParentID)
	res, err := DB.Exec(sql,
		l.Code,
		l.Type,
		l.ParentID,
	)
	if err != nil {
		log.Println("Error db.Exec in model.Location.Show", err)
		return nil, err
	}
	id, _ := res.LastInsertId()
	newLoc := Location{}
	err = DB.Get(&newLoc, "SELECT * FROM location WHERE id =?", id)
	if err != nil {
		log.Println("Error db.GET in model.Location.Show", err)
	}
	return &newLoc, nil
}
