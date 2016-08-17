package model

import (
	"github.com/jmoiron/sqlx"
	sys "github.com/mrtomyum/nava-sys/model"
	"golang.org/x/text/currency"
	"log"
	"time"
)

type BatchSale struct {
	sys.Base
	Recorded  *time.Time      `json:"recorded"`
	MachineID uint64          `json:"machine_id" db:"machine_id"`
	ColumnNo  int             `json:"column_no" db:"column_no"`
	Counter   int             `json:"counter"`
	SalePrice currency.Amount `json:"-" db:"sale_price"` // SalePrice search data from Last update Price of this Machine.Column
}

func (s *BatchSale) All(db *sqlx.DB) ([]*BatchSale, error) {
	log.Println("call model.BatchSale.All()")
	sales := []*BatchSale{}
	sql := `SELECT * FROM batch_sale`
	err := db.Select(&sales, sql)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(sales)
	return sales, nil
}

func NewBatchSale(db *sqlx.DB, columns []*BatchSale) error {
	// Call from controller.PostMachineBatchSale()
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	sql := `INSERT INTO batch_sale (
		recorded,
		machine_id,
		column_no,
		counter
		) VALUES(?,?,?,?)
	`
	for _, c := range columns {
		res, err := tx.Exec(sql,
			c.Recorded,
			c.MachineID,
			c.ColumnNo,
			c.Counter,
		)
		if err != nil {
			errTx := tx.Rollback()
			if errTx != nil {
				return errTx
			}
			return err
		}
		id, _ := res.LastInsertId()
		log.Println("Last Insert Id: ", id)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
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