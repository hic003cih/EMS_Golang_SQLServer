package main

import (
	"database/sql"
	"fmt"

	//_ "github.com/go-sql-driver/mysql"
	"log"
	"my-app/db"
	"my-app/web"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	//d, err := sql.Open("mysql", dataSource())
	log.Println("Begin DBConect")
	d, err := sql.Open("sqlserver", dataSource())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Connected!\n")
	defer d.Close()
	// CORS is enabled only in prod profile
	cors := os.Getenv("profile") == "prod"
	//先執行NewApp,把所有Handler註冊
	app := web.NewApp(db.NewDB(d), cors)
	//執行ListenAndServe
	err = app.Serve()

	log.Println("DBConect Sucess")
	log.Println("Error", err)
}

func dataSource() string {
	//host := "localhost"
	//pass := "password"

	server := "dev543"
	port := 1433
	user := "sa"
	password := "liteon1234"
	database := "EMS"

	// if os.Getenv("profile") == "prod" {
	// 	host = "db"
	// 	pass = os.Getenv("db_pass")
	// }

	// return "root:" + pass + "@tcp(" + host + ":3306)/EMS"
	//return "root:" + pass + "@tcp(" + host + ":3306)/emsgolang"

	// Create connection string
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
}
