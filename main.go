package main

import (
	"agent/Extension"
	"agent/Model"
	"agent/Network"
	"encoding/json"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

type OperationLogDocument struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	AgentUUID       string             `bson:"agentUUID"`
	ProcedureID     string             `bson:"procedureID"`
	InstructionUUID string             `bson:"instructionUUID"`
	ConductAt       time.Time          `bson:"conductAt"`
	ExitCode        int                `bson:"exitCode"`
	Log             string             `bson:"log"`
	Command         string             `bson:"command"` // Command 필드로 변경
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// stage 1 : 초기 정보 수집 단계 (정찰)

	sysdb := Model.NewSystemInfoDB()

	appdb := Model.NewApplicationDB()
	applist, err := appdb.SelectAllRecords()
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

	for i := 0; i < 15; i++ {
		offset := len(applist) / 15
		var chunk []Model.DapplicationDB
		chunk = applist[offset*i : offset*(i+1)]

		bapplist, err := appdb.ToJSON(chunk)
		if err != nil {
			return
		}

		// stage 2 : 네트워크 전송 단계
		hsItem := HSProtocol.HS{
			ProtocolID:     HSProtocol.HTTP, //
			HealthStatus:   HSProtocol.NEW,  //
			Command:        HSProtocol.SEND_AGENT_APP_INFO,
			Identification: 12345, // 아직 구현 안함
			Checksum:       0,     // 자동으로 채워줌
			TotalLength:    0,     // 자동으로 채워줌
			UUID:           uuid,
			Data:           bapplist,
		}
		ack, err := Network.SendHTTPRequest(hsItem)
		if err != nil {
			fmt.Println(err)
			return
		}
		if ack.Command == HSProtocol.ERROR_ACK {
			fmt.Println("Application 정보 송신 실패")
			break
		}
		//fmt.Printf("command : ")
		//fmt.Println(ack.Command)
	}

	bsysinfo, err := json.Marshal(sysInfo)
	hsItem := HSProtocol.HS{
		ProtocolID:     HSProtocol.HTTP, //
		HealthStatus:   HSProtocol.NEW,  //
		Command:        HSProtocol.SEND_AGENT_SYS_INFO,
		Identification: 12345, // 아직 구현 안함
		Checksum:       0,     // 자동으로 채워줌
		TotalLength:    0,     // 자동으로 채워줌
		UUID:           uuid,
		Data:           bsysinfo,
	}
	ack, err := Network.SendHTTPRequest(hsItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(len(bsysinfo))
	if ack.Command == HSProtocol.ERROR_ACK {
		fmt.Println("sysinfo 정보 송신 실패")
	}

	hsItem = HSProtocol.HS{
		ProtocolID:     HSProtocol.HTTP, //
		HealthStatus:   HSProtocol.WAIT, //
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: 12345, // 아직 구현 안함
		Checksum:       0,     // 자동으로 채워줌
		TotalLength:    0,     // 자동으로 채워줌
		UUID:           uuid,
		Data:           []byte{},
	}
	ack, err = Network.SendHTTPRequest(hsItem)
	if err != nil {
		fmt.Println(err)
		return
	}
	if ack.Command == HSProtocol.ERROR_ACK {
		fmt.Println("sysinfo 정보 송신 실패")
	}

	for {

		break
	}
}

// ToJSON: 구조체를 JSON 바이트로 마샬링
//func (s *ApplicationDB) ToJSON(data []DapplicationDB) ([]byte, error) {
//	return json.Marshal(data)
//}
//
//// FromJSON: JSON 바이트를 구조체로 언마샬링
//func (s *ApplicationDB) FromJSON(data []byte) ([]DapplicationDB, error) {
//	var result []DapplicationDB
//	err := json.Unmarshal(data, &result)
//	return result, err
//}
