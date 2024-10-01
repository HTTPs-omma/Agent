package Network

import (
	"agent/Core"
	"agent/Extension"
	"agent/Model"
	"encoding/json"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func SendAck(hs *HSProtocol.HS) error {
	hsItem := HSProtocol.HS{
		ProtocolID:     hs.ProtocolID,   //
		HealthStatus:   HSProtocol.WAIT, //
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: 12345, // 아직 구현 안함
		Checksum:       0,     // 자동으로 채워줌
		TotalLength:    0,     // 자동으로 채워줌
		UUID:           hs.UUID,
		Data:           []byte{},
	}

	_, err := SendPacket(hsItem)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func SendErrorAck(hs *HSProtocol.HS) error {
	hsItem := HSProtocol.HS{
		ProtocolID:     hs.ProtocolID,   //
		HealthStatus:   HSProtocol.WAIT, //
		Command:        HSProtocol.UPDATE_AGENT_STATUS,
		Identification: 12345, // 아직 구현 안함
		Checksum:       0,     // 자동으로 채워줌
		TotalLength:    0,     // 자동으로 채워줌
		UUID:           hs.UUID,
		Data:           []byte{},
	}

	_, err := SendPacket(hsItem)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

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

func SendLogData(hs *HSProtocol.HS, cmdLog string, instD *Core.InstructionData, resultCode int) error {

	logD, err := json.Marshal(&OperationLogDocument{
		ID:              primitive.ObjectID{},
		Command:         instD.Command,
		AgentUUID:       HSProtocol.ByteArrayToHexString(hs.UUID),
		ProcedureID:     instD.ID,
		InstructionUUID: "",
		ConductAt:       time.Now(),
		Log:             cmdLog,
		ExitCode:        resultCode,
	})
	if err != nil {
		fmt.Println("Error : ", err)
	}

	hsItem := HSProtocol.HS{
		ProtocolID:     hs.ProtocolID,
		HealthStatus:   0,
		Command:        HSProtocol.SEND_PROCEDURE_LOG,
		Identification: 12345,
		Checksum:       6789,
		TotalLength:    50,
		UUID:           hs.UUID,
		Data:           logD,
	}
	ack, err := SendPacket(hsItem)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if ack.Command == HSProtocol.ERROR_ACK {
		return fmt.Errorf("Error Ack!")
	}

	return nil
}

// 데이터를 수집하고 보내는 과정
func SendApplicationInfo() error {
	sysuil, err := Extension.NewSysutils()
	if err != nil {
		return err
	}
	strUuid := sysuil.GetUniqueID()

	fmt.Println("application 정보 가져오는 중...")
	//applist := Model.GetApplicationList_WMI() // 현재 사용 안함

	fileNames, err := Model.GetApplicationList_fileBase()
	if err != nil {
		return fmt.Errorf("GetApplication File Names error : %v", err)
	}

	prgdb, err := Model.NewProgramsDB()
	if err != nil {
		return fmt.Errorf("NewProgramsDB error: %v", err)
	}

	err = prgdb.DeleteAllRecords()
	if err != nil {
		return fmt.Errorf("prgdb.DeleteAllRecords error: %v", err)
	}

	for _, fileName := range fileNames {
		prgdb.InsertRecord(strUuid, fileName)
	}

	fmt.Printf("완료 len(data) : ")
	fmt.Println(len(fileNames))

	fileList, err := prgdb.SelectAllRecords()
	if err != nil {
		return err
	}

	uuid, err := HSProtocol.HexStringToByteArray(sysuil.GetUniqueID())
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		offset := len(fileNames) / 10
		var chunk []Model.ProgramsRecord
		chunk = fileList[offset*i : offset*(i+1)]

		bapplist, err := prgdb.ToJSON(chunk)
		if err != nil {
			return err
		}

		hsItem := HSProtocol.HS{
			ProtocolID:     HSProtocol.TCP,
			HealthStatus:   HSProtocol.NEW,
			Command:        HSProtocol.SEND_AGENT_APP_INFO,
			Identification: 12345, // 아직 구현 안함
			Checksum:       0,     // 자동으로 채워줌
			TotalLength:    0,     // 자동으로 채워줌
			UUID:           uuid,
			Data:           bapplist,
		}
		ack, err := SendPacket(hsItem)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if ack.Command == HSProtocol.ERROR_ACK {
			fmt.Println("Application 정보 송신 실패")
			break
		}
	}

	return nil
}

func SendSystemInfo() error {
	sysdb := Model.NewSystemInfoDB()
	sysuil, err := Extension.NewSysutils()
	if err != nil {
		return err
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
	fmt.Printf("해당 PC UUID : ")
	fmt.Println(sysuil.GetUniqueID())
	fmt.Println("sysInfo : " + sysInfo.Uuid)

	err = sysdb.InsertRecord(sysInfo)
	if err != nil {
		return err
	}

	uuid, err := HSProtocol.HexStringToByteArray(sysuil.GetUniqueID())
	if err != nil {
		return fmt.Errorf("HexStringToByteArray error : %v", err)
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
	ack, err := SendPacket(hsItem)
	if err != nil {
		return err
	}
	if ack.Command == HSProtocol.ERROR_ACK {
		return fmt.Errorf("sysinfo 정보 송신 실패")
	}

	return nil
}
