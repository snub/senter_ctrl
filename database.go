package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = nil

func connectDatabase(host string, port string, username string, password string, database string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db = d
    return nil
}

func saveSensorController(mac string, timestamp int) sql.Result {
	// Open doesn't open a connection. Validate DSN data:
	err := db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	result, err := db.Exec("INSERT INTO sensor_controller (mac_address, last_startup) VALUES (?, FROM_UNIXTIME(?))", mac, timestamp)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return result
}

func existsSensorController(mac string) bool {
	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}

	var sensorControllerId int64
	err = db.QueryRow("SELECT id FROM sensor_controller WHERE mac_address = ?", mac).Scan(&sensorControllerId)
	switch {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		panic(err.Error())
	default:
		return true
	}
}
