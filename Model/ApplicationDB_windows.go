package Model

import (
	"fmt"
	"github.com/yusufpapurcu/wmi"
	"log"
	"time"
)

type ApplicationDB struct {
	dbName string
}

func NewApplicationDB() (metaTable *ApplicationDB) {
	appDB := &ApplicationDB{"Application"}
	return appDB
}

type DapplicationDB struct {
	ID              int       // 내부 ID, 자동 증가
	Name            string    // 제품 이름
	Version         string    // 제품 버전
	Language        string    // 제품의 언어
	Vendor          string    // 제품 공급자
	InstallDate2    string    // 설치 날짜
	InstallLocation string    // 패키지 설치 위치
	InstallSource   string    // 설치 소스 위치
	PackageName     string    // 원래 패키지 이름
	PackageCode     string    // 패키지 식별자
	RegCompany      string    // 제품을 사용하는 것으로 등록된 회사 이름
	RegOwner        string    // 제품을 사용하는 것으로 등록된 사용자 이름
	URLInfoAbout    string    // 제품에 대한 정보가 제공되는 URL
	Description     string    // 제품 설명
	CreateAt        time.Time // 레코드 생성 시간
	UpdateAt        time.Time // 레코드 업데이트 시간
}

func (a *ApplicationDB) createTable() error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,       -- 내부 ID, 자동 증가
			Name VARCHAR(255),                          -- 제품 이름
			Version VARCHAR(50),                        -- 제품 버전
			Language VARCHAR(10),                       -- 제품의 언어
			Vendor VARCHAR(255),                        -- 제품 공급자
			InstallDate2 VARCHAR(20),                   -- 설치 날짜
			InstallLocation TEXT,                       -- 패키지 설치 위치
			InstallSource TEXT,                         -- 설치 소스 위치
			PackageName VARCHAR(255),                   -- 원래 패키지 이름
			PackageCode VARCHAR(255) UNIQUE NOT NULL  	-- 패키지 식별자 UUID
			RegCompany VARCHAR(255),                    -- 제품을 사용하는 것으로 등록된 회사 이름
			RegOwner VARCHAR(255),                      -- 제품을 사용하는 것으로 등록된 사용자 이름
			URLInfoAbout TEXT,                          -- 제품에 대한 정보가 제공되는 URL
			Description TEXT,                           -- 제품 설명
		    isDeleted bool, 							-- apllication 제거 여부를 파악함
			createAt DATETIME DEFAULT CURRENT_TIMESTAMP, -- 레코드 생성 시간
			updateAt DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 레코드 업데이트 시간
		    deletedAt DATETIME							-- 제거된 시간
		);
	`
	sqlStmt = fmt.Sprintf(sqlStmt, a.dbName)

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
	`, a.dbName, a.dbName)

	_, err = db.Exec(sqlModifyTrigger)
	if err != nil {
		return err
	}

	return nil
}

/*
refer : https://learn.microsoft.com/en-us/previous-versions/windows/desktop/legacy/aa394378(v=vs.85)
class Win32_Product : CIM_Product

	{
	  uint16   AssignmentType;
	  string   Caption;
	  string   Description;
	  string   IdentifyingNumber;
	  string   InstallDate;
	  datetime InstallDate2;
	  string   InstallLocation;
	  sint16   InstallState;
	  string   HelpLink;
	  string   HelpTelephone;
	  string   InstallSource;
	  string   Language;
	  string   LocalPackage;
	  string   Name;
	  string   PackageCache;
	  string   PackageCode;
	  string   PackageName;
	  string   ProductID;
	  string   RegOwner;
	  string   RegCompany;
	  string   SKUNumber;
	  string   Transforms;
	  string   URLInfoAbout;
	  string   URLUpdateInfo;
	  string   Vendor;
	  uint32   WordCount;
	  string   Version;
	};
*/
type Win32_Product struct {
	Name            string // 제품 이름
	Version         string // 제품 버전
	Language        string // 제품의 언어
	Vendor          string // 제품 공급자
	InstallDate2    string // 설치 날짜
	InstallLocation string // 패키지 설치 위치
	InstallSource   string // 설치 소스 위치
	PackageName     string // 원래 패키지 이름
	PackageCode     string // 패키지 식별자
	RegCompany      string // 제품을 사용하는 것으로 등록된 회사 이름
	RegOwner        string // 제품을 사용하는 것으로 등록된 사용자 이름
	URLInfoAbout    string // 제품에 대한 정보가 제공되는 URL
	Description     string // 제품 설명
}

func createQeruy() []Win32_Product {
	var dst []Win32_Product
	query := "SELECT Name, Version, Language, Vendor, InstallDate2, InstallLocation, InstallSource, PackageName, PackageCode, RegCompany, RegOwner, URLInfoAbout, Description FROM Win32_Product"
	err := wmi.Query(query, &dst)
	if err != nil {
		log.Fatalf("wmi query failed: %v", err)
	}
	return dst
}

func (a *ApplicationDB) insertRecord(data *DapplicationDB) error {
	// ProductID 가 있는지 확인 후 중복되는 것이 없으면 insert 하기

	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO %s ( Name, Version, Language, Vendor, 
        InstallDate2, InstallLocation, InstallSource, PackageName, PackageCode, RegCompany, 
        RegOwner, URLInfoAbout, Description ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, a.dbName)

	stmt, err := db.Prepare(query)
	fmt.Println(query)
	defer stmt.Close()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(&data)

	if err != nil {
		return err
	}

	return nil
}

/*
selectRecords()를 통해 반환된 DsystemInfoDB 객체의 값을 수정한 후,
수정된 객체를 updateRecord 함수의 매개변수로 전달합시오
*/
func (a *ApplicationDB) updateRecord(data *DapplicationDB) error {
	db, err := getDBPtr()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := a.selectRecords()
	if err != nil {
		return err
	}
	row := rows[0]
	data. = row.Uuid

	query := fmt.Sprintf(`UPDATE %s SET HostName = ?, OsName = ?, OsVersion = ?, Family = ?, Architecture = ?, KernelVersion = ?, BootTime = ?`, s.dbName)
	_, err = db.Exec(query, data.HostName, data.OsName, data.OsVersion, data.Family, data.Architecture, data.KernelVersion, data.BootTime)
	if err != nil {
		return err
	}

	return nil
}


func (s *ApplicationDB) selectByPackageCode(packageCode string) (*DapplicationDB ,error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM %s WHERE PackageCode = %s `, s.dbName, packageCode)
	row, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var rows []DapplicationDB

	for row.Next() {
		var data DapplicationDB

		err = row.Scan(&data)
		if err != nil {
			return nil, err
		}
		rows = append(rows, data)
	}

	return rows, nil
}

func (s *ApplicationDB) deleteRecord(uuid string) error {
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

func (s *ApplicationDB) selectRecords() ([]DsystemInfoDB, error) {
	db, err := getDBPtr()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM %s `, s.dbName)
	row, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var rows []DsystemInfoDB

	for row.Next() {
		var data DsystemInfoDB

		err = row.Scan(&data.id, &data.Uuid, &data.HostName, &data.OsName,
			&data.OsVersion, &data.Family, &data.Architecture, &data.KernelVersion,
			&data.BootTime, &data.createAt, &data.updateAt)
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
func (s *ApplicationDB) existRecord() (bool, error) {
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
