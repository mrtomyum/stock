package model

import (
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/sys/model"
	"github.com/shopspring/decimal"
	"log"
)

type Item struct {
	sys.Base
	SKU        string          `json:"sku"`
	NameTh     string          `json:"name_th" db:"name_th"`
	NameEn     string          `json:"name_en" db:"name_en"`
	StdPrice   decimal.Decimal `json:"std_price" db:"std_price"`
	StdCost    decimal.Decimal `json:"std_cost" db:"std_cost"`
	BaseUnitID uint64          `json:"base_unit_id" db:"base_unit_id"`
	CategoryID uint64          `json:"category_id" db:"category_id"`
	BrandID    uint64          `json:"brand_id" db:"brand_id"`
}

type ItemView struct {
	//sys.Base
	Item
	BaseUnitTH string `json:"base_unit_th" db:"base_unit_th"`
	BaseUnitEN string `json:"base_unit_en" db:"base_unit_en"`
	CategoryTH string `json:"category_th" db:"category_th"`
	CategoryEN string `json:"category_en" db:"category_en"`
}

type ItemBarcode struct {
	sys.Base
	ItemID uint64
	UnitID uint64
	Code   string
	Price  decimal.Decimal
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

func (i *Item) GetItemView(db *sqlx.DB) (*ItemView, error) {
	log.Println("call model.Item.GetItemView()")
	var iv ItemView
	sql := `
	SELECT
		item.sku,
		item.name_th,
		item.name_en,
		item.std_price,
		item.std_cost,
		item.base_unit_id,
		item.category_id,
		unit.name_th as base_unit_th,
		unit.name_en as base_unit_en,
		category.name_th as category_th,
		category.name_en as category_en
	FROM item
	LEFT JOIN unit ON item.base_unit_id = unit.id
	LEFT JOIN category ON item.category_id = category.id
	WHERE item.id = ?
	`
	//err := db.QueryRowx(sql, i.ID).StructScan(&iv)
	err := db.Get(&iv, sql, i.ID)
	if err != nil {
		log.Println("Error: model.Item.GetItemView/Query Error", err)
		return &iv, err
	}
	return &iv, nil
}

func (i *Item) Insert(db *sqlx.DB) (Item, error) {
	sql := `
		INSERT INTO item (
			sku,
			name_th,
			name_en,
			std_price,
			std_cost,
			base_unit_id,
			category_id
		) VALUES(
			?,?,?,?,?,?
		)
		`
	rs, err := db.Exec(sql,
		i.SKU,
		i.NameTh,
		i.NameEn,
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

func (i *Item) Update(db *sqlx.DB) (*Item, error) {
	sql := `
		UPDATE item
		SET
			sku = ?,
			name_th = ?,
			name_en = ?,
			std_price = ?,
			std_cost = ?,
			base_unit_id = ?,
			category_id = ?
		WHERE id = ?
	`
	_, err := db.Exec(sql,
		i.SKU,
		i.NameTh,
		i.NameEn,
		i.StdPrice,
		i.StdCost,
		i.BaseUnitID,
		i.CategoryID,
		i.ID,
	)
	if err != nil {
		log.Println("Error after db.Exec()")
		return nil, err
	}
	// Get updated record back from DB to confirm
	sql = `SELECT * FROM item WHERE id = ?`
	var updatedItem Item
	err = db.Get(&updatedItem, sql, i.ID)
	if err != nil {
		log.Println("Error after db.Get()")
		return nil, err
	}
	return &updatedItem, nil
}
