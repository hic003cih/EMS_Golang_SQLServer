package model

type SysPlant_mysql struct {
	Id          int    `json:"id"`
	Plantname   string `json:"plant_name"`
	Plantcode   string `json:"plant_code"`
	Plantdesc   string `json:"plant_desc"`
	Plantremark string `json:"plant_remark"`
}
