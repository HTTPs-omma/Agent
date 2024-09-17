package main

import (
	"agent/Extension"
	"agent/Model"
	"agent/Network"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// stage 1 : 초기 정보 수집 단계 (정찰)

	sysdb := Model.NewSystemInfoDB()

	//appdb := Model.NewApplicationDB()
	//applist, err := appdb.SelectAllRecords()
	//applist := Model.GetApplicationList()
	//err = appdb.DeleteAllRecord()
	//if err != nil {
	//	return
	//}
	//
	//for _, appD := range applist {
	//	err = appdb.InsertRecord(appD)
	//	if err != nil {
	//		return
	//	}
	//}

	sysuil, err := Extension.NewSysutils()
	if err != nil {
		return
	}

	sysInfo := &Model.DsystemInfoDB{
		0,
		sysuil.GetUniqueID(),
		sysuil.GetHostName(),
		sysuil.GetOsName(),
		sysuil.GetOsVersion(),
		sysuil.GetFamily(),
		sysuil.GetArchitecture(),
		sysuil.GetKernelVersion(),
		sysuil.GetBootTime(),
		time.Now(),
		time.Now(),
	}
	fmt.Println(sysuil.GetUniqueID())

	err = sysdb.InsertRecord(sysInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	uuid, err := HSProtocol.HexStringToByteArray(sysuil.GetUniqueID())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// stage 2 : 네트워크 전송 단계
	hsItem := HSProtocol.HS{
		ProtocolID:     HSProtocol.HTTP, //
		HealthStatus:   HSProtocol.NEW,  //
		Command:        HSProtocol.FETCH_INSTRUCTION,
		Identification: 12345, // 아직 구현 안함
		Checksum:       0,     // 자동으로 채워줌
		TotalLength:    0,     // 자동으로 채워줌
		UUID:           uuid,
		Data:           []byte{},
	}
	_, err = Network.SendHTTPRequest(hsItem)
	if err != nil {
		fmt.Println(err)
		return
	}

}
