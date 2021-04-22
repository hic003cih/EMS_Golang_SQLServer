package model

type SysPlant struct {
	PlantID     int    `json:"id"`
	PlantName   string `json:"plant_name"`
	PlantCode   string `json:"plant_code"`
	PlantDesc   string `json:"plant_desc"`
	PlantRemark string `json:"plant_remark"`
}
