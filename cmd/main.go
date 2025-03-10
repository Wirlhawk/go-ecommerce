package main

import (
	"database/sql"
	"log"

	"github.com/Wirlhawk/go-ecommerce/cmd/api"
	"github.com/Wirlhawk/go-ecommerce/config"
	"github.com/Wirlhawk/go-ecommerce/db"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)

	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB Sucessfully Connected")
}
