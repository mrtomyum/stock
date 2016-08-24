package model

import (
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/nava-sys/model"
	"golang.org/x/text/currency"
	"log"
	"time"
)

type BatchCounter struct {
	sys.Base
	Recorded  *time.Time      `json:"recorded"`
	MachineID uint64          `json:"machine_id" db:"machine_id"`
	ColumnNo  int             `json:"column_no" db:"column_no"`
	Counter   int             `json:"counter"`
	//SalePrice currency.Amount `json:"-" db:"sale_price"` // SalePrice search data from Last update Price of this Machine.Column
}

func (s *BatchCounter) All(db *sqlx.DB) ([]*BatchCounter, error) {
	log.Println("call model.BatchSale.All()")
	sales := []*BatchCounter{}
	sql := `SELECT * FROM batch_sale`
	err := db.Select(&sales, sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(sales)
	return sales, nil
}

func NewBatchCounter(db *sqlx.DB, sales []*BatchCounter) ([]*BatchCounter, error) {
	// Call from controller.PostMachineBatchSale()
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	sql := `INSERT INTO batch_sale (
		recorded,
		machine_id,
		column_no,
		counter
		) VALUES(?,?,?,?)
	`
	var ids []uint64
	for _, c := range sales {
		res, err := tx.Exec(sql,
			c.Recorded,
			c.MachineID,
			c.ColumnNo,
			c.Counter,
		)
		if err != nil {
			log.Println("error in tx.Exec(), res =", res, "Error: ", err)
			errRollback := tx.Rollback()
			if errRollback != nil {
				log.Println("errRollback", errRollback)
				return nil, errRollback
			}
			log.Println("tx.Rollback()", err)
			return nil, err
		}
		id, err := res.LastInsertId()
		log.Println("id = ", id, "err = ", err)
		ids = append(ids, uint64(id))
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	////Read from written DB.
	log.Println(ids)
	readSales, err := SelectBatchSale(db, ids)
	if err != nil {
		log.Println("Error in SelectBatchSale() = ", err)
		return nil, err
	}
	return readSales, nil
}

func SelectBatchSale(db *sqlx.DB, ids []uint64) ([]*BatchCounter, error) {
	sql := `SELECT * FROM batch_sale WHERE id = ?`
	sales := []*BatchCounter{}
	for _, id := range ids {
		var s BatchCounter
		log.Print("id: ", id)
		err := db.Get(&s, sql, id)
		if err != nil {
			return nil, err
		}
		sales = append(sales, &s)
		log.Print(sales)
	}
	return sales, nil
}

// บันทึกราคาจากหน้าตู้ ด้วยมือถือ
type BatchPrice struct {
	sys.Base
	Recorded  *time.Time      `json:"recorded"`
	MachineID uint64          `json:"machine_id"`
	ColumnNo  int             `json:"column_no"`
	Price     currency.Amount `json:"price"`
}

func (s *BatchPrice) All(db *sqlx.DB) ([]*BatchPrice, error) {
	log.Println("call model.BatchPrice.All()")
	prices := []*BatchPrice{}
	sql := `SELECT * FROM batch_price`
	err := db.Select(&prices, sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(prices)
	return prices, nil
}

// RouteMan สามารถขายสินค้านอกตู้ได้ในหลายๆกรณี เช่นยังเติมของไม่เสร็จ และเก็บเงินสดจากการขายนำส่งต่างหากได้
// พฤติกรรมงานจะเหมือนกับการขาย POS แต่สรุปยอดบิลวันละ 1 ใบ หลายๆรายการรวมกัน ตัดสต๊อคจากท้ายรถ ไม่ใช่จากตู้
// ลองออกแบบ Type นี้ไว้รับข้อมูลดิบก่อนหาข้อมูลสินค้า ตรวจความถูกต้องก่อนใช้ VanSale เก็บ เขียนลง DB
type VanSaleRawData struct {
	Recorded   *time.Time      `json:"recorded"`
	Barcode    string          `json:"barcode"`
	Qty        int             `json:"qty"`
	PriceUnit  currency.Amount `json:"price_unit"`
	PriceTotal currency.Amount `json:"price_total"`
}

// ใช้ type นี้  map DB
type VanSale struct {
	sys.Base
	Recorded   *time.Time      `json:"recorded"`
	Barcode    string          `json:"barcode"`
	ItemID     uint64          `json:"item_id"`
	Qty        int             `json:"qty"`
	UnitPrice  currency.Amount `json:"unit_price"`
	TotalPrice currency.Amount `json:"total_price"`
}

// Design this struct for data from VMC telemetry system.
//type SaleStatus int
//const (
//	COMPLETED SaleStatus = iota
//	INCOMPLETED
//)
//type RealTimeSale struct {
//	sys.Base
//	Recorded  *time.Time      `json:"recorded"`
//	MachineID uint64          `json:"machine_id"`
//	ColumnNo  int             `json:"column_no"`
//	ItemID    uint64          `json:"item_id"`
//	Price     currency.Amount `json:"price"`
//	//Status SaleStatus `json:"status"`
//}
