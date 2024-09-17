package Network

import (
	"agent/Execute"
	"bytes"
	"fmt"
	"github.com/HTTPs-omma/HTTPsBAS-HSProtocol/HSProtocol"
	"net"
)

func send() {
	// TCP 서버에 연결
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to TCP server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to TCP server")

	// HS 객체 생성
	hsItem := HSProtocol.HS{
		ProtocolID:     1,
		HealthStatus:   0,
		Command:        0b0000000110,
		Identification: 12345,
		Checksum:       6789,
		TotalLength:    50,
		UUID:           [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		Data:           []byte{0xAA, 0xBB, 0xCC},
	}

	// HS 객체를 직렬화 (예: ToBytes 함수 사용)
	HSMgr := HSProtocol.NewHSProtocolManager()
	data, err := HSMgr.ToBytes(&hsItem)
	if err != nil {
		fmt.Println("Error serializing HS object:", err)
		return
	}

	// 서버로 데이터 전송 (Payload 요청)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data to server:", err)
		return
	}
	fmt.Println("Data sent to server")

	// 데이터 응답 (PayLoad 받아옴)
	msg := make([]byte, 1024)
	conn.Read(msg)
	msg = bytes.ReplaceAll(msg, []byte{0x00}, []byte{})

	fmt.Println("Data received from server : ", string(msg))

	shell := Execute.PowerShell{
		Name:        "powershell",
		IsAvailable: true,
	}
	rst, err := shell.Execute(string(msg))
	if err != nil {
		fmt.Println("Error executing powershell command: ", err)
	}
	fmt.Println(rst)

	// 서버로 데이터 전송 (Payload 요청)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data to server:", err)
		return
	}
	fmt.Println("Data sent to server")

	// log 데이터 전송 객체 생성
	hsItem = HSProtocol.HS{
		ProtocolID:     1,
		HealthStatus:   0,
		Command:        0b0000000111,
		Identification: 12345,
		Checksum:       6789,
		TotalLength:    0,
		UUID:           [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
		Data:           []byte(rst),
	}

	// HS 객체를 직렬화 (예: ToBytes 함수 사용)
	HSMgr = HSProtocol.NewHSProtocolManager()
	data, err = HSMgr.ToBytes(&hsItem)

	hsItem2, err := HSMgr.Parsing(data)
	le := hsItem2.TotalLength - 24
	fmt.Println("legnth : ", le)
	fmt.Println("legnth : ", len(rst))
	fmt.Println("debug : ", string(hsItem2.Data[:le]))

	if err != nil {
		fmt.Println("Error serializing HS object:", err)
		return
	}

	// 서버로 데이터 전송 (Payload 요청)
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data to server:", err)
		return
	}
	fmt.Println("Data sent to server")

	// 데이터 응답 (PayLoad 받아옴)
	msg = make([]byte, 1024)
	conn.Read(msg)
	msg = bytes.ReplaceAll(msg, []byte{0x00}, []byte{})

	fmt.Println("Data received from server : ", string(msg))

}
