package Model

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type sqlite_master struct {
	Type     string
	name     string
	tbl_name string
	rootpage string
	sql      string
}

func TestSystemInfoDB_createTable(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "create DB test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			s := NewSystemInfoDB()
			s.createTable()
			if err != nil {
				t.Fatalf(s.dbName + " : DB를 생성할 수 없습니다.")
			}

			dbPtr, err := getDBPtr()
			if err != nil {
				t.Fatalf(s.dbName + " : DB 포인터를 가져올 수 없습니다. getDBPtr() 함수 오류")
			}

			//query := fmt.Sprintf("select * from sqlite_master")
			query := fmt.Sprintf("select * from sqlite_master where name = '%s'", s.dbName)

			dsys := &sqlite_master{}

			rst := dbPtr.QueryRow(query).Scan(&dsys.Type, &dsys.name, &dsys.tbl_name, &dsys.rootpage, &dsys.sql)
			if rst != nil {
				t.Fatalf(s.dbName + " : 생성된 테이블이 존재하지 않습니다.")
			}
			assert.Equal(t, dsys.name, s.dbName)

			//for rows.Next() {
			//	rows.Scan(&dsys.Type, &dsys.name, &dsys.tbl_name, &dsys.rootpage, &dsys.sql)
			//	fmt.Println("rst : ", dsys.name)
			//}

		})
	}
}

//func TestNewSystemInfoDB_insertRecord(t *testing.T) {
//	type fields struct {
//		dbName string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		{name: "create DB test"},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var err error
//			s := NewSystemInfoDB()
//
//			data := &DsystemInfoDB{
//				Uuid:          "123e4567-e89b-12d3-a456-426614174000",
//				HostName:      "example-host",
//				OsName:        "Linux",
//				OsVersion:     "5.4.0-42-generic",
//				Family:        "Unix",
//				Architecture:  "x86_64",
//				KernelVersion: "5.4.0",
//				BootTime:      time.Now().Add(-time.Hour * 24 * 5), // 5일 전으로 설정
//			}
//			err = s.insertValue(data)
//			if err != nil {
//				panic(err)
//			}
//			data, err = s.selectValue()
//			if err != nil {
//				panic(err)
//			}
//			fmt.Println("sql result : \n", data.createAt)
//		})
//	}
//}
//
//func TestNewSystemInfoDB_selectRecord(t *testing.T) {
//	type fields struct {
//		dbName string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		{name: "insert Record into DB"},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			var err error
//			s := NewSystemInfoDB()
//
//			data := &DsystemInfoDB{}
//
//			data, err = s.selectValue()
//			if err != nil {
//				panic(err)
//			}
//			fmt.Println("sql result : \n", data.createAt)
//		})
//	}
//}
