package model

import (
	"github.com/jmoiron/sqlx"
	"log"
	sys "github.com/mrtomyum/nava-sys/model"
)

type Item struct {
	sys.Base
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	StdPrice   int64  `json:"stdPrice" db:"std_price"`
	StdCost    int64  `json:"stdCost" db:"std_cost"`
	BaseUnitID uint64 `json:"baseUnitID" db:"base_unit_id"`
	CategoryID uint64 `json:"categoryID" db:"category_id"`
}

type ItemView struct {
	sys.Base
	Item
	BaseUnitTH string `json:"baseUnitTH" db:"base_unit_th"`
	BaseUnitEN string `json:"baseUnitEN" db:"base_unit_en"`
	CategoryTH string `json:"categoryTH" db:"category_th"`
	CategoryEN string `json:"categoryEN" db:"category_en"`
}

type ItemBarcode struct {
	sys.Base
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
		return nil, err
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
	sql := `
	SELECT
		item.sku,
		item.name,
		item.std_price,
		item.std_cost,
		item.baseunit_id,
		item.category_id,
		unit.th as baseunit_th,
		unit.en as baseunit_en,
		category.th as category_th,
		category.en as category_en
	FROM item
	LEFT JOIN unit ON item.baseunit_id = unit.id
	LEFT JOIN category ON item.category_id = category.id
	WHERE item.id = ?
	`
	err := db.QueryRowx(sql, i.ID).StructScan(&iv)

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
