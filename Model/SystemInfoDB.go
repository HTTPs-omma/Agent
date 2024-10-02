package Model

import (
	"errors"
	"fmt"
	"time"
)

type SystemInfoDB struct {
	dbName string
}

type DsystemInfoDB struct {
	ID            int       `json:"id"`
	Uuid          string    `json:"uuid"`
	HostName      string    `json:"host_name"`
	OsName        string    `json:"os_name"`
	OsVersion     string    `json:"os_version"`
	Family        string    `json:"family"`
	Architecture  string    `json:"architecture"`
	KernelVersion string    `json:"kernel_version"`
	BootTime      time.Time `json:"boot_time"`
	IP            string    `json:"ip"`  // IP 필드 추가
	MAC           string    `json:"mac"` // MAC 필드 추가
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewSystemInfoDB() *SystemInfoDB {
	sysDB := &SystemInfoDB{"SystemInfo"}
	sysDB.CreateTable()
	return sysDB
}

func (s *SystemInfoDB) CreateTable() error {
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
			IP string,                               -- IP 추가
			MAC string,                              -- MAC 추가
			createAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updateAt DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`
	sqlStmt = fmt.Sprintf(sqlStmt, s.dbName)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlModifyTrigger := fmt.Sprintf(`
		CREATE TRIGGER IF NOT EXISTS update_ModificationTime
		AFTER UPDATE ON %s
		FOR EACH ROW
		BEGIN	
			UPDATE %s SET
				updateAt = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
	`, s.dbName, s.dbName)

	_, err = db.Exec(sqlModifyTrigger)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) InsertRecord(data *DsystemInfoDB) error {
	// 데이터베이스에는 단 하나의 Row만을 보장해야 함
	isExist, err := s.existRecord()
	if err != nil {
		return err
	}
	if isExist == true {
		err = s.UpdateRecord(data)
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

	query := fmt.Sprintf(`INSERT INTO %s (uuid, HostName, OsName, OsVersion, Family, Architecture, KernelVersion, BootTime, IP, MAC) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, s.dbName)
	stmt, err := db.Prepare(query)
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(data.Uuid, data.HostName, data.OsName, data.OsVersion, data.Family, data.Architecture, data.KernelVersion, data.BootTime, data.IP, data.MAC)

	if err != nil {
		return err
	}

	return nil
}

/*
selectRecords()를 통해 반환된 DsystemInfoDB 객체의 값을 수정한 후,
수정된 객체를 updateRecord 함수의 매개변수로 전달하십시오
*/
func (s *SystemInfoDB) UpdateRecord(data *DsystemInfoDB) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := s.SelectRecords()
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return errors.New("SystemInfo 테이블에 저장된 데이터가 없습니다.")
	}
	row := rows[0]
	data.Uuid = row.Uuid

	query := fmt.Sprintf(`UPDATE %s SET HostName = ?, OsName = ?, OsVersion = ?, Family = ?, Architecture = ?, KernelVersion = ?, BootTime = ?, IP = ?, MAC = ?`, s.dbName)
	_, err = db.Exec(query, data.HostName, data.OsName, data.OsVersion, data.Family, data.Architecture, data.KernelVersion, data.BootTime, data.IP, data.MAC)
	if err != nil {
		return err
	}

	return nil
}

func (s *SystemInfoDB) SelectRecords() ([]DsystemInfoDB, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM %s`, s.dbName)
	row, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var rows []DsystemInfoDB

	for row.Next() {
		var data DsystemInfoDB
		err = row.Scan(&data.ID, &data.Uuid, &data.HostName, &data.OsName, &data.OsVersion, &data.Family, &data.Architecture, &data.KernelVersion, &data.BootTime, &data.IP, &data.MAC, &data.CreatedAt, &data.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rows = append(rows, data)
	}

	return rows, nil
}

/*
*
하나 이상의 row 행이 있는지 검사한다.
*/
func (s *SystemInfoDB) existRecord() (bool, error) {
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
