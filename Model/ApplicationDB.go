package Model

import (
	"fmt"
	"github.com/yusufpapurcu/wmi"
	"log"
)

type ApplicationDB struct {
	dbName string
}

func NewApplicationDB() (metaTable *ApplicationDB) {
	//appPtr := getDBPtr()

	return
}

func isTablePresent() bool {
	return true
}

func (t *ApplicationDB) createTable() {

	//var dst []Win32_Process
	//sqlStmt := `
	//	CREATE TABLE IF NOT EXISTS FileMetadata (
	//		id INTEGER PRIMARY KEY AUTOINCREMENT,    -- 내부 ID, 자동 증가
	//		uuid TEXT NOT NULL unique,                      -- UUIDv4
	//		name string,
	//	);
	//`
	//_, err := t.db.Exec(sqlStmt)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

type Win32_Product struct {
	Name            string
	Version         string
	InstallDate     string
	InstallLocation string
}

func createQeruy() {
	var dst []Win32_Product
	query := "SELECT Name, Version, InstallDate, InstallLocation FROM Win32_Product"
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Fatalf("wmi query failed: %v", err)
	}
	fmt.Printf("%+v\n", dst)
}
