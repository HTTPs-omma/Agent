package Network

import (
	"agent/Core"
	"agent/Execute"
	"encoding/json"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
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

func Test_getPayload(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Test case 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := godotenv.Load()
			if err != nil {
				log.Fatalf("Error loading .env file")
			}

			uuid, err := HSProtocol.HexStringToByteArray("937640a858ad48e9bc2787e8c4456ced")
			if err != nil {
				panic(err)
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

			//fmt.Println(uuid)

			inst := &Core.InstructionData{}
			ack, err := SendHTTPRequest(hsItem)
			//fmt.Println(ack, err)
			instD, err := inst.GetInstData(ack.Data)
			if err != nil {
				fmt.Println("Error : ", err)
				return
			}
			//fmt.Println(instD.Command)

			shell := Execute.Cmd{}
			cmdLog, err := shell.Execute(instD.Command)
			fmt.Println(cmdLog)
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
					return
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
				ack, err := SendHTTPRequest(hsItem)
				fmt.Printf("commmand: %b \n", ack.Command)
				return
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
			ack, err = SendHTTPRequest(hsItem)
			//fmt.Println(ack, err)

			cmdLog, err = shell.Execute(instD.Cleanup)
			//fmt.Println(cmdLog)
			if err != nil {
				fmt.Println("Error : ", err)
			}

			fmt.Println()
		})
	}
}
