package Model

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type SystemInfoDB struct {
	dbName string
}

type DsystemInfoDB struct {
	Uuid          string
	HostName      string
	OsName        string
	OsVersion     string
	Family        string
	Architecture  string
	KernelVersion string
	BootTime      time.Time
	createAt      time.Time
	updateAt      time.Time
}

func (s *SystemInfoDB) createTable() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()
	s.dbName = "SystemInfo"

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS FileMetadata (
			id INTEGER PRIMARY KEY AUTOINCREMENT,    -- 내부 ID, 자동 증가
			uuid TEXT NOT NULL unique,               -- UUIDv4
			HostName string,
			OsName string,
			OsVersion string,
			Family string,
			Architecture string,
			KernelVersion string,
			BootTime DATETIME,
			createAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updateAt DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	sqlCreateTrigger := `
		CREATE TRIGGER set_createTime_ModificationTime
		AFTER INSERT ON SystemInfo
		FOR EACH ROW
		BEGIN
			UPDATE tableName SET
				createTime = CURRENT_TIMESTAMP,
				ModificationTime = CURRENT_TIMESTAMP
				WHERE id = NEW.id;
		END;
	`

	sqlModifyTrigger := `
		CREATE TRIGGER update_ModificationTime
		AFTER UPDATE ON SystemInfo
		FOR EACH ROW
		BEGIN
			UPDATE SystemInfo SET
				ModificationTime = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(sqlCreateTrigger)
	if err != nil {
		return err
	}

	_, err = db.Exec(sqlModifyTrigger)
	if err != nil {
		return err
	}

	return nil
}

/*
반드시 하나를 유지 하자
*/
func (s *SystemInfoDB) insertValue(data DsystemInfoDB) error {
	isExist, err := s.selectExists()
	if err != nil {
		return err
	}
	if isExist == true {
		err := s.updateValue(data)
		if err != nil {
			return err
		}
		return nil
	}

	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO %s (Uuid, HostName,
       OsName, OsVersion, Family, Architecture, KernelVersion,
       BootTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, s.dbName)

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) updateValue(data DsystemInfoDB) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}

	field := reflect.ValueOf(data)
	t := field.Type()
	for i := 0; i < field.NumField(); i++ {

		switch t {
		case reflect.TypeOf(""):
			if strings.Compare() {

			}
		case reflect.TypeOf(time.Time{}):
			fmt.Println("Field is a time.Time : ", field.Interface().(time.Time))
		default:
			fmt.Println("Unknown type : ", t)
		}

	}

	query := fmt.Sprintf(`UPDATE %s SET uuid = ?, HostName = ?, OsName = ?, OsVersion = ?, Family = ?, Architecture = ?, KernelVersion = ?, BootTime = ?`, s.dbName)
	_, err = db.Exec(query, data.Uuid, data.HostName, data.OsName, data.OsVersion, data.Family, data.Architecture, data.KernelVersion, data.BootTime)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) deleteValue(uuid string) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE Uuid = ?`, s.dbName)
	_, err = db.Exec(query, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) selectValue() (*DsystemInfoDB, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT * FROM %s `, s.dbName)
	row := db.QueryRow(query)

	var data DsystemInfoDB
	err = row.Scan(&data.Uuid, &data.HostName, &data.OsName,
		&data.OsVersion, &data.Family, &data.Architecture, &data.KernelVersion,
		&data.BootTime, &data.createAt, &data.updateAt)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (s *SystemInfoDB) selectExists() (bool, error) {
	db, err := getDBPtr()
	if err != nil {
		return false, err
	}

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE Uuid = ?)`, s.dbName)
	var exists bool
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
