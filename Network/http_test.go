package Network

import (
	"agent/Core"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"github.com/joho/godotenv"
	"log"
	"testing"
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
			fmt.Println(ack, err)
			instD, err := inst.GetInstData(ack.Data)
			if err != nil {
				fmt.Println("Error : ", err)
			}
			fmt.Println(instD.Command)

			//shell := Execute.PowerShell{}
			//cmdLog, err := shell.Execute(instD.Command)
			//if err != nil {
			//
			//}
			//shell.Execute(instD.Cleanup)
			fmt.Println()
		})
	}
}
