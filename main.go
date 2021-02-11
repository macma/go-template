package main

import (
	"database/sql"
	"eass/bpe-scheduler/config"
	"eass/bpe-scheduler/db"
	"fmt"
)

func main() {
	var cfgEnv config.EnvConfig
	config.InitViperConfig(&cfgEnv, "dev.env")
	config.Config = &cfgEnv

	var cfgDbconn *sql.DB
	cfgDbconn = db.InitDB()
	config.DbConnection = cfgDbconn
	fmt.Print("aaa")
}
