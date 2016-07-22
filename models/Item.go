package models

import (
	"database/sql"
	"log"
)

type Items []*Item

func (this *Item) Index(db *sql.DB) ([]*Item, error) {
	sql := `SELECT
		id,
		catetory_id,
		sku,
		name,
		std_price,
		str_cost,
		base_unit
		FROM item WHERE
		id=?`
	rows, err := db.Query(sql,this.ID)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var items Items
	for rows.Next() {
		i := new(Item)
		rows.Scan(
			&i.ID,
			&i.CategoryID,
			&i.SKU,
			&i.Name,
			&i.StdPrice,
			&i.StdCost,
			&i.BaseUnit,
		)
		items = append(items, i)
	}
	return items, err
}