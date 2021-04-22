package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"my-app/model"
	"time"
)

type DB interface {
	GetPlants() ([]*model.SysPlant, error)
	InsertPlant(model.SysPlant) (int, error)
	GetPlant(int) (*model.SysPlant, error)
	UpdatePlant(int, model.SysPlant) error
	RemovePlant(int) error
}

// type MySQLDB struct {
// 	db *sql.DB
// }
type MSSQLDB struct {
	db *sql.DB
}

func NewDB(db *sql.DB) DB {
	return MSSQLDB{db: db}
}

func (d MSSQLDB) GetPlants() ([]*model.SysPlant, error) {
	ctx := context.Background()
	// Check if database is alive.
	err := d.db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}
	//rows, err := d.db.Query("select name, details from technologies")
	//取DB資料
	log.Println("取資料")
	// rows, err := d.db.Query("select id,plant_code,plant_name,plant_desc,plant_remark from EMS.dbo.sys_plant")

	rows, err := d.db.QueryContext(ctx, "select id,plant_code,plant_name,plant_desc,plant_remark from dbo.sys_plant")

	log.Println(rows.Scan())
	if err != nil {
		return nil, err
	}
	//最後關閉DB連結
	defer rows.Close()

	//把modelPlant存在memory中,做一個array,可存多個SysPlant的Model
	var plant []*model.SysPlant

	//對取出來的DB值做循環
	for rows.Next() {
		log.Println("test")
		//建一個SysPlant的model
		t := new(model.SysPlant)
		//把DB取出來的值用scan存到新建的model中
		err = rows.Scan(&t.PlantID, &t.PlantCode, &t.PlantName, &t.PlantDesc, &t.PlantRemark)
		log.Println(t.PlantID)
		if err != nil {
			return nil, err
		}
		//把存有DB值的model,append到plant的array
		plant = append(plant, t)
	}
	log.Println(plant)
	return plant, nil
} /*func (d MySQLDB) CreateTechnologies() ([]*model.SysPlant, error) {
	//rows, err := d.db.Query("select name, details from technologies")
	//取DB資料
	rows, err := d.db.Exec("select id,plant_code,plant_name,plant_desc,plant_remark from sys_plant")
	if err != nil {
		return nil, err
	}
	//最後關閉DB連結
	defer rows.Close()

	//把modelPlant存在memory中,做一個array,可存多個SysPlant的Model
	var plant []*model.SysPlant


	return plant, nil
}*/
func (d MSSQLDB) GetPlant(plantID int) (*model.SysPlant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// row := d.db.QueryRowContext(ctx, `
	// id,
	// plant_code,
	// plant_name,
	// plant_desc,
	// plant_remark
	// rom sys_plant
	// WHERE id = ?`, plantID)
	tsql := fmt.Sprintf("Select id,plant_code,plant_name,plant_desc,plant_remark from sys_plant WHERE id = (%d);", plantID)
	row := d.db.QueryRowContext(ctx, tsql)

	plant := &model.SysPlant{}
	err := row.Scan(
		&plant.PlantID,
		&plant.PlantCode,
		&plant.PlantName,
		&plant.PlantDesc,
		&plant.PlantRemark,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return plant, nil
}
func (d MSSQLDB) RemovePlant(plantID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// _, err := d.db.ExecContext(ctx, `DELETE FROM sys_plant where id = ?`, plantID)

	tsql := fmt.Sprintf("DELETE FROM sys_plant where id = %d;",
		plantID)
	_, err := d.db.ExecContext(ctx, tsql)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
func (d MSSQLDB) InsertPlant(model model.SysPlant) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	log.Println("insert")
	log.Println(model.PlantName)
	// result, err := d.db.ExecContext(ctx, `INSERT INTO  EMS.dbo.sys_plant
	// (
	// plant_name,
	// plant_code,
	// plant_desc,
	// plant_remark) VALUES ( ?, ?, ?, ?)`,
	// 	model.PlantName,
	// 	model.PlantCode,
	// 	model.PlantDesc,
	// 	model.PlantRemark)
	tsql := fmt.Sprintf("INSERT INTO  EMS.dbo.sys_plant (plant_name, plant_code,plant_desc,plant_remark) VALUES ('%s','%s','%s','%s');",
		model.PlantName, model.PlantCode, model.PlantDesc, model.PlantRemark)
	result, err := d.db.ExecContext(ctx, tsql)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}
func (d MSSQLDB) UpdatePlant(plantID int, model model.SysPlant) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// _, err := d.db.ExecContext(ctx, `DELETE FROM sys_plant where id = ?`, plantID)
	log.Println(model.PlantName)
	tsql := fmt.Sprintf("Update sys_plant Set plant_name='%s',plant_code='%s',plant_remark='%s',plant_desc='%s' where id = %d;", model.PlantName, model.PlantCode, model.PlantRemark, model.PlantDesc,
		plantID)
	log.Println("update test2")
	_, err := d.db.ExecContext(ctx, tsql)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
