package model

import (
	"log"
	"github.com/jmoiron/sqlx"
)

type Item struct {
	Base
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	StdPrice   int64  `json:"std_price" db:"std_price"`
	StdCost    int64  `json:"std_cost" db:"std_cost"`
	BaseUnitID uint64 `json:"baseunit_id" db:"baseunit_id"`
	CategoryID uint64 `json:"category_id" db:"category_id"`
}

type ItemBarcode struct {
	ItemID uint64
	UnitID uint64
	Code   string
	Price  int64
}

type Items []*Item

func (i *Item) All(db *sqlx.DB) (Items, error) {
	sql := `SELECT * FROM item`
	rows, err := db.Queryx(sql)
	if err != nil {
		log.Println("Error: db.Queryx in Item.All(): ", err)
	}
	defer rows.Close()
	var items Items
	for rows.Next() {
		i := new(Item)
		err = rows.StructScan(&i)
		if err != nil {
			log.Println("Error: rows.StructScan in Item.All(): ", err)
			return nil, err
		}
		items = append(items, i)
		log.Println("Read item:", i)
	}
	return items, nil
}

func (i *Item) FindItemByID(db *sqlx.DB) (Item, error) {
	var item Item
	sql := "SELECT * FROM item WHERE id = ?"
	err := db.QueryRowx(sql, i.ID).StructScan(&item)
	if err != nil {
		log.Println("Error: FindItemByID/Query Error", err)
		return item, err
	}
	return item, nil
}

func (i *Item) New(db *sqlx.DB) error {
	sql := `
		INSERT INTO item (
			sku,
			name,
			std_price,
			std_cost,
			baseunit_id,
			category_id
		) VALUES(
			?,?,?,?,?,?
		)
		`
	rs, err := db.Exec(sql,
		i.SKU,
		i.Name,
		i.StdPrice,
		i.StdCost,
		i.BaseUnitID,
		i.CategoryID,
	)
	if err != nil {
		log.Println("Error=>Item.New/db.Exec:> ", err)
		return err
	}
	lastID, _ := rs.LastInsertId()

	// Check result
	err = db.QueryRowx("SELECT * FROM item WHERE id =?", lastID).
		StructScan(&i)
	if err != nil {
		return err
	}
	log.Println("Success Insert New Item: ", i)
	return nil
}