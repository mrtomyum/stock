# Counter

## ขั้นตอนปฏิบัติงาน

เฟสแรกตอนนี้จะยังไม่แทรคสต๊อค

1. User ยิงรหัสตู้ UI จะ GET /v1/counter/last/machine_code/:code
1. Backend จะ response JSON machine_last_counter มาให้
1. User กรอกยอด Counter ทีละ Column
1. เมื่อครบทุกคอลัมน์ User Submit เพื่อ "สร้างใบนำส่งเงิน"
- ให้ UI ตรวจสอบการป้อนข้อมูลต่างๆ
- และส่งผลลัพธ์กลับด้วย POST /v1/counter/new/machine_code/:code
1. Server Backend Process ยอดนำส่งเงิน โดยเทียบผลต่าง Counter และคูณ ราคาขายล่าสุดทีละ column เพื่อ บันทึกลงตารางนำส่งเงินของ User นี้ เป็นการตั้งใบนำส่งเงิน แต่ยังไม่ Complete ต้องให้ Cashier ทำการบันทึกยอดนับเงินก่อน