package models

type LocType int

const (
	STORE LocType = 1 + iota
	VAN
	MACHINE
	INSPECTION
	VENDOR
	DAMAGE
)

type Location struct {
	ID       uint64  `json:"id"`
	ParentID uint64  `json:"parent_id"`
	Name     string  `json:"name"`
	Type     LocType `json:"type"`
}

type Unit struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type ItemUnit struct {
	ID     uint64 `json:"id"`
	ItemID uint64 `json:"Item_id"`
	UnitID uint64 `json:"unit_id"`
	Ratio  int    `json:"ratio"`
	IsSale bool   `json:"is_sale"`
	IsBuy  bool   `json:"is_buy"`
}

type Item struct {
	ID         uint64 `json:"id"`
	CategoryID uint64 `json:"category_id"`
	SKU        string `json:"sku"`
	Name       string `json:"name"`
	StdPrice   int64  `json:"std_price"`
	StdCost    int64  `json:"std_cost"`
	BaseUnit   Unit   `json:"base_unit"`
}

type ItemBarcode struct {
	ItemID uint64
	UnitID uint64
	Code   string
	Price  int64
}

type Stock struct {
	LocationID uint64 `json:"location_id"`
	ItemID     uint64 `json:"item_id"`
	Quantity   int64  `json:"quantity"`
}

type machineType int

const (
	CAN machineType = 1 + iota
	CUP
	CUP_FRESH_COFFEE
	CUP_NOODLE
	SEE_THROUGH
)

type brand int

const (
	NATIONAL brand = 1 + iota
	SANDEN
	FUJI_ELECTRIC
	CIRBOX
)

type ClientType int

const (
	FACTORY ClientType = 1 + iota
	EDUCATION
	OFFICE
)

type Client struct {
	ID   uint64
	Name string
	Type ClientType
}

type Place struct {
	ID       uint64
	ClientID uint64
	Name     string
	Lat      float64
	Long     float64
}

type Machine struct {
	ID        uint64
	PlaceID   uint64
	Name      string
	Type      machineType
	Brand     brand
	Model     string
	Selection int //จำนวน Column หรือช่องเก็บสินค้า
	LocNumber []int
	//LocRow int  	//จำนวนแถว และคอลัมน์ไว้ทำ Schematic Profile  หน้าตู้
	//LocCol int
}

type carBrand int

const (
	SUZUKI carBrand = 1 + iota
	TATA
)

type Vehicle struct {
	ID        uint64
	Name      string   // V1, V2,...
	NamePlate string   // ทะเบียนรถ
	Brand     carBrand // ยี่ห้อ
}

type RouteMan struct {
	ID        uint64
	Name      string
	VehicleID uint64
}
