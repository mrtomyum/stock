package model

import (
	"log"
	"github.com/jmoiron/sqlx"
)

type Item struct {
	Base
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	StdPrice   int64  `json:"std_price"`
	StdCost    int64  `json:"std_cost"`
	BaseUnitID uint64 `json:"baseunit_id"`
	CategoryID uint64 `json:"category_id"`
}

type ItemBarcode struct {
	ItemID uint64
	UnitID uint64
	Code   string
	Price  int64
}

type Items []*Item

func (i *Item) All(db *sqlx.DB) ([]*Item, error) {
	sql := `SELECT
		id,
		sku,
		name,
		std_price,
		std_cost,
		base_unit_id,
		category_id
		FROM item
		`
	rows, err := db.Queryx(sql)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var items Items
	for rows.Next() {
		i := new(Item)
		rows.Scan(
			&i.ID,
			&i.SKU,
			&i.Name,
			&i.StdPrice,
			&i.StdCost,
			&i.BaseUnitID,
			&i.CategoryID,
		)
		items = append(items, i)
	}
	return items, err
}

func (i *Item) FindItemByID(db *sqlx.DB) ([]*Item, error) {
	sql := "SELECT * FROM item WHERE id = ?"
	rows, err := db.Queryx(sql)
	if err != nil {
		log.Println("FindItemByID/Query Error", err)
	}
	defer rows.Close()

	var items Items
	err = rows.StructScan(&items)
	if err != nil {
		log.Println("FindItemByID/rows.StructScan:", err)
		return nil, err
	}
	return items, nil
}