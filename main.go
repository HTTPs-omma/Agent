package main

import (
	"agent/Core"
	"agent/Execute"
	"agent/Extension"
	"agent/Model"
	"agent/Network"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"github.com/joho/godotenv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5)
		return
	}

	sysutil, err := Extension.NewSysutils()
	if err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5)
		return
	}
	uuid, err := HSProtocol.HexStringToByteArray(sysutil.GetUniqueID())
	if err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5)
		return
	}

	// 새로 생성된 에이전트를 등록
	hsItem := HSProtocol.HS{
		ProtocolID:     HSProtocol.HTTP, //
		HealthStatus:   HSProtocol.NEW,  //
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: 12345, // 아직 구현 안함
		Checksum:       0,     // 자동으로 채워줌
		TotalLength:    0,     // 자동으로 채워줌
		UUID:           uuid,
		Data:           []byte{},
	}
	ack, err := Network.SendPacket(hsItem)
	if err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5)
		return
	}
	if ack.Command == HSProtocol.ERROR_ACK {
		return
	}

	// stage 1 : 초기 정보 수집 단계 (정찰)
	err = Network.SendApplicationInfo()
	if err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5)
		return
	}
	err = Network.SendSystemInfo()
	if err != nil {
		fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
		time.Sleep(5)
		return
	}

	// stage 2-3 : 반복 실행
	for {
		fmt.Print("fetch instruction : ")
		time.Sleep(3 * time.Second)
		agsdb, err := Model.NewAgentStatusDB()
		if err != nil {
			fmt.Println("(5 초뒤 종료)에러 발생 : " + err.Error())
			time.Sleep(5)
			return
		}
		agsRcrd, err := agsdb.SelectAllRecords()
		protocol := agsRcrd[0].Protocol
		// ======= stage 2. fetch Data =========
		hsItem := HSProtocol.HS{
			ProtocolID:     uint8(protocol),
			HealthStatus:   HSProtocol.RUN,
			Command:        HSProtocol.FETCH_INSTRUCTION,
			Identification: 12345,
			Checksum:       6789,
			TotalLength:    50,
			UUID:           uuid,
			Data:           []byte{},
		}
		inst := &Core.InstructionData{}
		ack, err := Network.SendPacket(hsItem)
		if err != nil {
			fmt.Println("Error : ", err)
			continue
		}

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

		// ======= stage 3. run Command =========
		err = Core.ChangeStatusToRun(&hsItem)
		if err != nil {
			fmt.Println("Error : ", err)
			continue
		}
		shell := Execute.Cmd{}
		cmdLog, err := shell.Execute(instD.Command)
		fmt.Println("cmdLog : " + cmdLog)
		if err != nil {
			err = Network.SendLogData(&hsItem, cmdLog, instD, Network.EXIT_FAIL)
			if err != nil {
				fmt.Println("Error : ", err)
				continue
			}
			fmt.Printf("commmand: %b \n", ack.Command)
			continue
		}

		cmdLog, err = shell.Execute(instD.Cleanup)
		if err != nil {
			fmt.Println("Error : ", err)
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
