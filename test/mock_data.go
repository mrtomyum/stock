package test

import "github.com/mrtomyum/stock/model"

var MockBatchSales string = `
	[
		{
		  "recorded": "2016-08-22T00:00:10+07:00",
		  "machine_id": 1,
		  "column_no": 6,
		  "counter": 100
		},
		{
		  "recorded": "2016-08-22T00:00:11+07:00",
		  "machine_id": 1,
		  "column_no": 7,
		  "counter": 100
		},
		{
		  "recorded": "2016-08-22T00:00:12+07:00",
		  "machine_id": 1,
		  "column_no": 8,
		  "counter": 100
		}
	]
`

// Mock User table
var mockUsers = []*model.User{
	{Name: "tom", Title: model.ADMIN},
	{Name: "kwang", Title: model.ADMIN},
	{Name: "eak", Title: model.STOREMAN},
	{Name: "lek", Title: model.ROUTEMAN},
	{Name: "tam", Title: model.ROUTEMAN},
	{Name: "est", Title: model.ROUTEMAN},
}

