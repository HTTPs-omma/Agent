package Model

import (
	"fmt"
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

func NewSystemInfoDB() *SystemInfoDB {
	sysDB := &SystemInfoDB{"SystemInfo"}
	return sysDB
}

func (s *SystemInfoDB) createTable() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS %s (
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
	sqlStmt = fmt.Sprintf(sqlStmt, s.dbName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlCreateTrigger := `
		CREATE TRIGGER IF NOT EXISTS set_createTime_ModificationTime
		AFTER INSERT ON %s
		FOR EACH ROW
		BEGIN
			UPDATE tableName SET
				createTime = CURRENT_TIMESTAMP,
				ModificationTime = CURRENT_TIMESTAMP
				WHERE id = NEW.id;
		END;
	`
	sqlCreateTrigger = fmt.Sprintf(sqlCreateTrigger, s.dbName)

	sqlModifyTrigger := `
		CREATE TRIGGER IF NOT EXISTS update_ModificationTime
		AFTER UPDATE ON %s
		FOR EACH ROW
		BEGIN
			UPDATE SystemInfo SET
				ModificationTime = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
	`
	sqlModifyTrigger = fmt.Sprintf(sqlModifyTrigger, s.dbName)

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

func (s *SystemInfoDB) insertValue(data *DsystemInfoDB) error {
	err := s.createTable()
	if err != nil {
		return err
	}
	isExist, err := s.selectExists()
	if err != nil {
		return err
	}
	if isExist == true {
		err = s.updateValue(data)
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

	query := fmt.Sprintf(`INSERT INTO %s (uuid, HostName,
       OsName, OsVersion, Family, Architecture, KernelVersion,
       BootTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, s.dbName)

	_, errorss := db.Query(`NSERT INTO SystemInfo (uuid, HostName, 
       OsName, OsVersion, Family, Architecture, KernelVersion,
       BootTime) VALUES (?,""?, "", "", "", "", "", "")
	`)

	if errorss != nil {
		fmt.Errorf("에러!")
		return err
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(data.Uuid, data.HostName, data.OsName,
		data.OsVersion, data.Family, data.Architecture,
		data.KernelVersion, data.BootTime)
	//fmt.Println(rst.LastInsertId())
	fmt.Println("debug=-===============")

	if err != nil {
		return err
	}

	return nil
}

/*
selectValue()를 통해 반환된 DsystemInfoDB 객체의 값을 수정한 후,
수정된 객체를 updateValue 함수의 매개변수로 전달합시오
*/
func (s *SystemInfoDB) updateValue(data *DsystemInfoDB) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

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
	defer db.Close()

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
	defer db.Close()

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
	defer db.Close()

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s)`, s.dbName)
	var exists bool
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
