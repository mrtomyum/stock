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

type ItemView struct {
	Item
	CategoryEN string `json:"category_en" db:"category_en"`
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

func (i *Item) FindItemByID(db *sqlx.DB) (ItemView, error) {
	var iv ItemView
	sql := `SELECT
		item.sku,
		item.name,
		item.std_price,
		item.std_cost,
		item.baseunit_id,
		item.category_id,
		category.EN as category_en
		FROM item LEFT JOIN category
		ON item.category_id = category_id
		WHERE item.id = ?
	`
	err := db.QueryRowx(sql, i.ID).StructScan(&iv)
	//err := db.QueryRow(sql, i.ID).Scan(
	//	&i.SKU,
	//	&i.Name,
	//	&i.StdPrice,
	//	&i.StdCost,
	//	&i.BaseUnitID,
	//	&i.Category.EN,
	//)
	if err != nil {
		log.Println("Error: FindItemByID/Query Error", err)
		return iv, err
	}

	return iv, nil
}

func (i *Item) New(db *sqlx.DB) (Item, error) {
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
	var item Item
	if err != nil {
		log.Println("Error=>Item.New/db.Exec:> ", err)
		return item, err
	}
	lastID, _ := rs.LastInsertId()

	// Check result
	err = db.QueryRowx("SELECT * FROM item WHERE id =?", lastID).
		StructScan(&item)
	if err != nil {
		return item, err
	}
	log.Println("Success Insert New Item: ", i)
	return item, nil
}