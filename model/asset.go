package model

import (
	sys "github.com/mrtomyum/sys/model"
)

// Asset ทะเบียนสินทรัพย์ จะต้องมีการกำหนดผู้รับผิดชอบแต่ละชิ้น
// โดยให้ตารางอื่นอ้าง AssetId กลับมา
type Asset struct {
	sys.Base
	Name    string
	Code    string
	UserId  uint64
	BuyDate Date
}

type Vehicle struct {
	sys.Base
	Name      string // V1, V2,...
	NamePlate string // ทะเบียนรถ
	Brand     string // ยี่ห้อ
	AssetId   uint64
	Asset     Asset
}

// ข้อมูลฝ่ายบริการเติมสินค้า ลงทะเบียนเบิกกุญแจรถตอนเช้า จะมีผลกับการขาย VanSale  ไม่ต้องระบุรหัสรถ และผู้ขับ RouteMan อีก
type RouteMan struct {
	sys.Base
	UserId    uint64
	Driver    sys.User
	VehicleID uint64
}
