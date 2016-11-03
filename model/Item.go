package model

import (
	sys "github.com/mrtomyum/sys/model"
	"github.com/shopspring/decimal"
	"log"
)

type Item struct {
	sys.Base
	SKU        string          `json:"sku"`
	Name       string          `json:"name" db:"name"`
	NameEn     string          `json:"name_en" db:"name_en"`
	StdPrice   decimal.Decimal `json:"std_price" db:"std_price"`
	StdCost    decimal.Decimal `json:"std_cost" db:"std_cost"`
	BaseUnitId uint64          `json:"base_unit_id" db:"base_unit_id"`
	CategoryId uint64          `json:"category_id" db:"category_id"`
	BrandID    uint64          `json:"brand_id" db:"brand_id"`
}

type ItemView struct {
	//sys.Base
	Item
	BaseUnit   string `json:"base_unit" db:"base_unit"`
	BaseUnitEn string `json:"base_unit_en" db:"base_unit_en"`
	Category   string `json:"category" db:"category"`
	CategoryEn string `json:"category_en" db:"category_en"`
}

type ItemBar struct {
	sys.Base
	ItemId uint64 `json:"item_id" db:"item_id"`      // สินค้า
	UnitId uint64 `json:"unit_id" db:"unit_id"`      // หน่วยนับ
	Code   string `json:"code" db:"code"`            // รหัสบาร์โค้ด
	Ratio  int    `json:"ratio" db:"ratio"`          // เป็นตัวคูณ ratio in times of BaseUnit.
	Price  decimal.Decimal `json:"price" db:"price"` // ราคาขายเฉพาะของบาร์โค้ดนั้น(ถ้ามี)
}

type Brand struct {
	sys.Base
	Name   string
	NameEn string
}

func (i *Item) GetAll() ([]*Item, error) {
	sql := `
	SELECT
		id,
		created,
		updated,
		deleted,
		sku,
		name_th,
		name_en,
		std_price,
		std_cost,
		base_unit_id,
		category_id,
		brand_id
	FROM item`
	rows, err := DB.Queryx(sql)
	if err != nil {
		log.Println("Error: db.Queryx in Item.All(): ", err)
		return nil, err
	}
	defer rows.Close()
	var items []*Item
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

func (i *Item) GetItemView() (*ItemView, error) {
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
	err := DB.Get(&iv, sql, i.Id)
	if err != nil {
		log.Println("Error: model.Item.GetItemView/Query Error", err)
		return &iv, err
	}
	return &iv, nil
}

func (i *Item) Insert() (Item, error) {
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
	rs, err := DB.Exec(sql,
		i.SKU,
		i.Name,
		i.NameEn,
		i.StdPrice,
		i.StdCost,
		i.BaseUnitId,
		i.CategoryId,
	)
	var item Item
	if err != nil {
		log.Println("Error=>Item.New/db.Exec:> ", err)
		return item, err
	}
	lastID, _ := rs.LastInsertId()

	// Check result
	err = DB.QueryRowx("SELECT * FROM item WHERE id =?", lastID).
		StructScan(&item)
	if err != nil {
		return item, err
	}
	log.Println("Success Insert New Item: ", i)
	return item, nil
}

func (i *Item) Update() (*Item, error) {
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
	_, err := DB.Exec(sql,
		i.SKU,
		i.Name,
		i.NameEn,
		i.StdPrice,
		i.StdCost,
		i.BaseUnitId,
		i.CategoryId,
		i.Id,
	)
	if err != nil {
		log.Println("Error after db.Exec()")
		return nil, err
	}
	// Get updated record back from DB to confirm
	sql = `SELECT * FROM item WHERE id = ?`
	var updatedItem Item
	err = DB.Get(&updatedItem, sql, i.Id)
	if err != nil {
		log.Println("Error after db.Get()")
		return nil, err
	}
	return &updatedItem, nil
}
