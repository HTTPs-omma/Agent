package main

import (
	"agent/Core"
	"agent/Execute"
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

const (
	EXIT_SUCCESS = 1
	EXIT_Unknown = 0
	EXIT_FAIL    = -1
)

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
		HealthStatus:   HSProtocol.NEW,  //
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

	// stage 3 : 반복 실행
	for {
		fmt.Print("fetch instruction : ")
		time.Sleep(3 * time.Second)

		uuid, err := HSProtocol.HexStringToByteArray(sysuil.GetUniqueID())
		if err != nil {
			break
		}

		hsItem := HSProtocol.HS{
			ProtocolID:     1,
			HealthStatus:   0,
			Command:        HSProtocol.FETCH_INSTRUCTION,
			Identification: 12345,
			Checksum:       6789,
			TotalLength:    50,
			//UUID:           [16]byte{0xc3, 0xcb, 0x84, 0x23, 0x34, 0x16, 0x49, 0x76, 0x94, 0x56, 0x9d, 0x75, 0x9a, 0x8a, 0x13, 0xe7},
			UUID: uuid,
			Data: []byte{},
		}

		inst := &Core.InstructionData{}
		ack, err := Network.SendHTTPRequest(hsItem)
		//fmt.Println(ack, err)
		instD, err := inst.GetInstData(ack.Data)
		if err != nil {
			fmt.Println("Error : ", err)
			continue
		}

		if len(ack.Data) < 1 {
			fmt.Println("... NoData Wait")
			continue
		}
		fmt.Println("... success")

		hsItem = HSProtocol.HS{
			ProtocolID:     HSProtocol.HTTP, //
			HealthStatus:   HSProtocol.RUN,  //
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
			fmt.Println("Status 정보 변경 실패")
		}

		shell := Execute.Cmd{}
		cmdLog, err := shell.Execute(instD.Command)
		fmt.Println("cmdLog : " + cmdLog)
		if err != nil {
			fmt.Println("Error : ", err)

			logD, err := json.Marshal(&OperationLogDocument{
				ID:              primitive.ObjectID{},
				Command:         instD.Command,
				AgentUUID:       "937640a858ad48e9bc2787e8c4456ced",
				ProcedureID:     instD.ID,
				InstructionUUID: "",
				ConductAt:       time.Now(),
				Log:             "",
				ExitCode:        EXIT_FAIL,
			})
			if err != nil {
				fmt.Println("Error : ", err)
				continue
			}

			hsItem = HSProtocol.HS{
				ProtocolID:     1,
				HealthStatus:   0,
				Command:        HSProtocol.SEND_PROCEDURE_LOG,
				Identification: 12345,
				Checksum:       6789,
				TotalLength:    50,
				//UUID:           [16]byte{0xc3, 0xcb, 0x84, 0x23, 0x34, 0x16, 0x49, 0x76, 0x94, 0x56, 0x9d, 0x75, 0x9a, 0x8a, 0x13, 0xe7},
				UUID: uuid,
				Data: logD,
			}
			ack, err := Network.SendHTTPRequest(hsItem)
			fmt.Printf("commmand: %b \n", ack.Command)
			continue
		}

		logD, err := json.Marshal(&OperationLogDocument{
			ID:              primitive.ObjectID{},
			Command:         instD.Command,
			AgentUUID:       "937640a858ad48e9bc2787e8c4456ced",
			ProcedureID:     instD.ID,
			InstructionUUID: "",
			ConductAt:       time.Now(),
			Log:             cmdLog,
			ExitCode:        EXIT_SUCCESS,
		})
		if err != nil {
			fmt.Println("Error : ", err)
		}

		hsItem = HSProtocol.HS{
			ProtocolID:     1,
			HealthStatus:   0,
			Command:        HSProtocol.SEND_PROCEDURE_LOG,
			Identification: 12345,
			Checksum:       6789,
			TotalLength:    50,
			//UUID:           [16]byte{0xc3, 0xcb, 0x84, 0x23, 0x34, 0x16, 0x49, 0x76, 0x94, 0x56, 0x9d, 0x75, 0x9a, 0x8a, 0x13, 0xe7},
			UUID: uuid,
			Data: logD,
		}
		ack, err = Network.SendHTTPRequest(hsItem)
		//fmt.Println(ack, err)

		cmdLog, err = shell.Execute(instD.Cleanup)
		//fmt.Println(cmdLog)
		if err != nil {
			fmt.Println("Error : ", err)
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
			fmt.Println("Status 정보 변경 실패")
		}
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
