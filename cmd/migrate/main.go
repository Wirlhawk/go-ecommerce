package main

import (
	"log"
	"os"

	configs "github.com/Wirlhawk/go-ecommerce/config"
	"github.com/Wirlhawk/go-ecommerce/db"
	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db, err := db.NewMySQLStorage(mysqlCfg.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "mysql", driver)
	if err != nil {
		log.Fatal(err.Error())
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err.Error())
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err.Error())
		}
	}
}
