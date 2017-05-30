package model

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
)

type Location struct {
	Base
	Code     string      `json:"code"`
	Type     LocType     `json:"type"`
	ParentId uint64      `json:"-" db:"parent_id"`
	Movable  bool        `json:"movable"`
	Child    []*Location `json:"nodes,omitempty" db:"-"`
}

func (l *Location) Size() int {
	var size int = len(l.Child)
	for _, c := range l.Child {
		size += c.Size()
	}
	return size
}

func (l *Location) AddTree(nodes ...*Location) bool {
	var size = l.Size()
	for _, node := range nodes {
		if node.ParentId == l.Id {
			l.Child = append(l.Child, node)
		} else {
			for _, c := range l.Child {
				if c.AddTree(node) {
					break
				}
			}
		}
	}
	return l.Size() == size+len(nodes)
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

// Method Get จะดึง Location ของตัวมันเองและลูกๆ ที่มันเป็นแม่ *ยังไม่สามารถ Reqursive ลงไปหาหลานๆได้
func (l *Location) Get(db *sqlx.DB) ([]*Location, error) {
	// TODO: แก้ Select ให้สามารถ Recursive  Parent_id ได้
	sql := `
		SELECT * FROM location
		WHERE id = ?
		OR parent_id = ?
		`
	locations := []*Location{}
	err := db.Select(&locations, sql, l.Id, l.Id)
	if err != nil {
		log.Fatal("Error in model.Select..", err)
		return nil, err
	}
	return locations, nil
}

func (l *Location) Insert(db *sqlx.DB) (*Location, error) {
	sql := `
		INSERT INTO location (
			code,
			type,
			parent_id
		)
		VALUES (?, ?, ?, ?)
	`
	log.Println("Test Location receiver:", l.Code, l.Type, l.ParentId)
	res, err := db.Exec(sql,
		l.Code,
		l.Type,
		l.ParentId,
	)
	if err != nil {
		log.Println("Error db.Exec in model.Location.Show", err)
		return nil, err
	}
	id, _ := res.LastInsertId()
	newLoc := Location{}
	err = db.Get(&newLoc, "SELECT * FROM location WHERE id =?", id)
	if err != nil {
		log.Println("Error db.GET in model.Location.Show", err)
	}
	return &newLoc, nil
}

type LocType int

const (
	ROOT       LocType = iota
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